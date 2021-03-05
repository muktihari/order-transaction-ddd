package ordering_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/muktihari/order-transaction-ddd/ordering"
	"github.com/muktihari/order-transaction-ddd/persistent/inmem"
	"github.com/muktihari/order-transaction-ddd/transaction"
	"github.com/shopspring/decimal"
)

func TestMakeOrder(t *testing.T) {
	var (
		customers = inmem.NewCustomerRepository()
		products  = inmem.NewProductRepository()
		coupons   = inmem.NewCouponRepository()
		logistics = inmem.NewLogisticsParner()
		orders    = inmem.NewOrderRepository(coupons, products)
		s         = ordering.NewService(orders, customers, products, coupons, logistics)
	)

	tt := []struct {
		Name       string
		CustomerID string
		Expected   *transaction.Order
	}{
		{Name: "Make Order", CustomerID: "CUSTOMER1", Expected: &transaction.Order{
			Customer: transaction.Customer{
				ID: "CUSTOMER1", Name: "Hari", PhoneNumber: "+62-12345", Email: "example@email.com", Address: "No, Street, City, Indonesia",
			},
			Cart:   []transaction.CartItem{},
			Status: transaction.OrderStatusOpen,
		}},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			o, err := s.MakeOrder(context.Background(), tc.CustomerID)
			if err != nil {
				t.Fatalf("got err, expected nil")
			}
			if o.ID == "" {
				t.Errorf("got empty string, expected not empty")
			}
			o.ID = ""
			if diff := cmp.Diff(o, tc.Expected); diff != "" {
				fmt.Println(diff)
				t.Fatal("different")
			}
		})
	}
}

func TestAddProduct(t *testing.T) {
	var (
		customers = inmem.NewCustomerRepository()
		products  = inmem.NewProductRepository()
		coupons   = inmem.NewCouponRepository()
		logistics = inmem.NewLogisticsParner()
		orders    = inmem.NewOrderRepository(coupons, products)
		s         = ordering.NewService(orders, customers, products, coupons, logistics)
	)

	tt := []struct {
		Name      string
		OrderID   string
		ProductID string
		Expected  *transaction.Order
	}{
		{Name: "Add Product", OrderID: "ORDER_OPEN", ProductID: "PRODUCT1", Expected: &transaction.Order{
			ID: "ORDER_OPEN",
			Customer: transaction.Customer{
				ID: "CUSTOMER1", Name: "Hari", PhoneNumber: "+62-12345", Email: "example@email.com", Address: "No, Street, City, Indonesia",
			},
			Cart: []transaction.CartItem{
				{
					Product:  &transaction.Product{ID: "PRODUCT1", Name: "Sony Xperia 10", Price: decimal.NewFromInt(500), Quantity: 200},
					Quantity: 5,
				},
			},
			Price:  decimal.NewFromInt(500 * 5),
			Status: transaction.OrderStatusOpen,
		}},
	}

	ctx := context.Background()
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := s.AddProduct(ctx, tc.OrderID, tc.ProductID, 5); err != nil {
				t.Fatalf("got %v, expected nil", err)
			}

			o, err := orders.FindByID(ctx, tc.OrderID)
			if err != nil {
				t.Fatalf("got %v, expected nil", err)
			}

			if diff := cmp.Diff(o, tc.Expected); diff != "" {
				fmt.Println(diff)
				t.Fatal("different")
			}
		})
	}
}

func TestApplyCoupon(t *testing.T) {
	var (
		customers = inmem.NewCustomerRepository()
		products  = inmem.NewProductRepository()
		coupons   = inmem.NewCouponRepository()
		logistics = inmem.NewLogisticsParner()
		orders    = inmem.NewOrderRepository(coupons, products)
		s         = ordering.NewService(orders, customers, products, coupons, logistics)
	)

	tt := []struct {
		Name       string
		OrderID    string
		CouponCode string
		Expected   *transaction.Order
	}{
		{Name: "Apply Coupon", OrderID: "ORDER_WITH_PRODUCT", CouponCode: "DISCOUNT_20%", Expected: &transaction.Order{
			ID: "ORDER_WITH_PRODUCT",
			Coupon: transaction.Coupon{
				Code:     "DISCOUNT_20%",
				Quantity: 100,
				Amount:   decimal.NewFromFloat(0.2),
				Type:     transaction.CouponTypePercentage,
			},
			Customer: transaction.Customer{
				ID: "CUSTOMER1", Name: "Hari", PhoneNumber: "+62-12345", Email: "example@email.com", Address: "No, Street, City, Indonesia",
			},
			Cart: []transaction.CartItem{
				{
					Product:  &transaction.Product{ID: "PRODUCT1", Name: "Sony Xperia 10", Price: decimal.NewFromInt(500), Quantity: 200},
					Quantity: 5,
				},
			},
			Price:               decimal.NewFromInt(500 * 5),
			PriceAfterReduction: decimal.NewFromInt(500 * 5).Sub(decimal.NewFromInt(500 * 5).Mul(decimal.NewFromFloat(0.2))),
			Status:              transaction.OrderStatusOpen,
		}},
	}

	ctx := context.Background()
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := s.ApplyCoupon(ctx, tc.OrderID, tc.CouponCode); err != nil {
				t.Fatalf("got %v, expected nil", err)
			}

			o, err := orders.FindByID(ctx, tc.OrderID)
			if err != nil {
				t.Fatalf("got %v, expected nil", err)
			}
			o.Coupon.Begin = time.Time{}
			o.Coupon.End = time.Time{}

			if diff := cmp.Diff(o, tc.Expected); diff != "" {
				fmt.Println(diff)
				t.Fatal("different")
			}
		})
	}
}

func TestSubmitOrder(t *testing.T) {
	var (
		customers = inmem.NewCustomerRepository()
		products  = inmem.NewProductRepository()
		coupons   = inmem.NewCouponRepository()
		logistics = inmem.NewLogisticsParner()
		orders    = inmem.NewOrderRepository(coupons, products)
		s         = ordering.NewService(orders, customers, products, coupons, logistics)
	)

	tt := []struct {
		Name       string
		OrderID    string
		CouponCode string
		Expected   *transaction.Order
	}{
		{Name: "Apply Coupon", OrderID: "ORDER_WITH_PRODUCT_AND_COUPON", CouponCode: "DISCOUNT_20%", Expected: &transaction.Order{
			ID: "ORDER_WITH_PRODUCT_AND_COUPON",
			Coupon: transaction.Coupon{
				Code:     "DISCOUNT_20%",
				Quantity: 100,
				Amount:   decimal.NewFromFloat(0.2),
				Type:     transaction.CouponTypePercentage,
			},
			Customer: transaction.Customer{
				ID: "CUSTOMER1", Name: "Hari", PhoneNumber: "+62-12345", Email: "example@email.com", Address: "No, Street, City, Indonesia",
			},
			Cart: []transaction.CartItem{
				{
					Product:  &transaction.Product{ID: "PRODUCT1", Name: "Sony Xperia 10", Price: decimal.NewFromInt(500), Quantity: 200},
					Quantity: 5,
				},
			},
			Price:               decimal.NewFromInt(500 * 5),
			PriceAfterReduction: decimal.NewFromInt(500 * 5).Sub(decimal.NewFromInt(500 * 5).Mul(decimal.NewFromFloat(0.2))),
			Status:              transaction.OrderStatusSubmitted,
		}},
	}

	ctx := context.Background()
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := s.SubmitOrder(ctx, tc.OrderID); err != nil {
				t.Fatalf("got %v, expected nil", err)
			}

			o, err := orders.FindByID(ctx, tc.OrderID)
			if err != nil {
				t.Fatalf("got %v, expected nil", err)
			}

			o.Coupon.Begin = time.Time{}
			o.Coupon.End = time.Time{}

			if diff := cmp.Diff(o, tc.Expected); diff != "" {
				fmt.Println(diff)
				t.Fatal("different")
			}
		})
	}
}
