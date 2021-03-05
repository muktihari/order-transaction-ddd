package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/muktihari/decimalcodec"
	"github.com/muktihari/order-transaction-ddd/handling"
	"github.com/muktihari/order-transaction-ddd/ordering"
	"github.com/muktihari/order-transaction-ddd/persistent/inmem"
	"github.com/muktihari/order-transaction-ddd/persistent/mongodb"
	"github.com/muktihari/order-transaction-ddd/persistent/mongodb/migration"
	"github.com/muktihari/order-transaction-ddd/transaction"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	httpAddr    = flag.String("httpAddr", ":8080", "server http address")
	mongoURI    = flag.String("mongoURI", "mongodb://localhost:27017", "mongodb connection URI")
	repo        = flag.String("repo", "inmem", "use repository: inmem, mongo")
	migrate     = flag.Bool("migrate", false, "migrate predefined data to mongo")
	httpAddrEnv = os.Getenv("HTTP_ADDRESS")
	mongoURIEnv = os.Getenv("MONGO_URI")
	repoEnv     = os.Getenv("REPO")
	migrateEnv  = os.Getenv("MIGRATE")
)

func main() {
	flag.Parse()

	if httpAddrEnv != "" {
		*httpAddr = httpAddrEnv
	}
	if mongoURIEnv != "" {
		*mongoURI = mongoURIEnv
	}
	if repoEnv != "" {
		*repo = repoEnv
	}
	if migrateEnv != "" {
		m, _ := strconv.ParseBool(migrateEnv)
		*migrate = m
	}

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})

	var logistics transaction.LogisticsPartner
	var customers transaction.CustomerRepository
	var products transaction.ProductRepository
	var coupons transaction.CouponRepository
	var orders transaction.OrderRepository

	// since logistics partner is another service's domain, use inmem mock
	logistics = inmem.NewLogisticsParner()

	// inmem
	switch *repo {
	case "inmem":
		customers = inmem.NewCustomerRepository()
		products = inmem.NewProductRepository()
		coupons = inmem.NewCouponRepository()
		orders = inmem.NewOrderRepository(coupons, products)
	case "mongo":
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx,
			options.Client().
				ApplyURI(*mongoURI).
				SetRegistry(func() *bsoncodec.Registry {
					rb := bsoncodec.NewRegistryBuilder()
					bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
					bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
					decimalcodec.RegisterEncodeDecoder(rb)
					return rb.Build()
				}()),
		)
		if err != nil {
			logger.Fatalf("could not connect to mongodb: %v", err)
		}
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			logger.Fatalf("could not ping mongodb: %v", err)
		}

		db := client.Database("transaction-order")
		customers = mongodb.NewCustomerRepository(db)
		products = mongodb.NewProductRepository(db)
		coupons = mongodb.NewCouponRepository(db)
		orders = mongodb.NewOrderRepository(client, db)

		if *migrate {
			if err := migration.MigratePredefinedData(context.Background(), client); err != nil {
				logger.Fatalf("could not migrate pre-defined data: %v", err)
			}
		}
	}

	var orderingService ordering.Service
	orderingService = ordering.NewService(orders, customers, products, coupons, logistics)
	orderingService = ordering.NewLoggingService(logger, orderingService)
	orderingService = ordering.NewInstrumentinService(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "ordering",
			Name:      "request_counter",
			Help:      "Total number of processed request",
		}, []string{"method", "error"}),
		prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "ordering",
			Name:      "request_latency",
			Help:      "Summary of request latency",
		}, []string{"method", "err"}),
		orderingService,
	)
	orderingHandler := ordering.MakeHandler(orderingService)

	var handlingService handling.Service
	handlingService = handling.NewService(orders, products, logistics)
	handlingService = handling.NewLoggingService(logger, handlingService)
	handlingService = handling.NewInstrumentingService(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "handling",
			Name:      "request_counter",
			Help:      "Total ordering request",
		}, []string{"method", "error"}),
		prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "handling",
			Name:      "request_latency",
			Help:      "Summary of request latency",
		}, []string{"method", "err"}),
		handlingService,
	)
	handlingHandler := handling.MakeHandler(handlingService)

	r := chi.NewMux()
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) { promhttp.Handler().ServeHTTP(w, r) })

	r.Mount("/ordering/v1", orderingHandler)
	r.Mount("/handling/v1", handlingHandler)

	_ = chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		logger.Infof("[%s] %s", method, route)
		return nil
	})

	errs := make(chan error, 2)
	go func() {
		logger.Infof("listening to %s", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Infof("terminated: %v", <-errs)

}
