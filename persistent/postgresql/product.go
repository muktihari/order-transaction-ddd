package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/muktihari/order-transaction-ddd/transaction"
)

type productRepository struct {
	db    *sql.DB
	table string
}

// NewProductRepository creates new product repository in memory
func NewProductRepository(db *sql.DB) transaction.ProductRepository {
	return &productRepository{db, "products"}
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*transaction.Product, error) {
	stmt, err := r.db.PrepareContext(ctx, "select * from "+r.table+" where id = $1")
	if err != nil {
		return nil, err
	}
	sqlRows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	var p transaction.Product
	if sqlRows.Next() == false {
		return nil, transaction.ErrProductNotFound
	}
	if err := Map(sqlRows, &p); err != nil {
		return nil, err
	}
	sqlRows.Close()

	return &p, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]transaction.Product, error) {
	stmt, err := r.db.PrepareContext(ctx, "select * from "+r.table)
	if err != nil {
		return nil, err
	}
	sqlRows, err := stmt.Query(nil)
	if err != nil {
		return nil, err
	}

	var ps []transaction.Product
	for sqlRows.Next() {
		var p transaction.Product
		if err := Map(sqlRows, &p); err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	if len(ps) == 0 {
		return nil, transaction.ErrProductNotFound
	}

	return ps, nil
}

func (r *productRepository) Update(ctx context.Context, product *transaction.Product) error {
	p, err := r.FindByID(ctx, product.ID)
	if err != nil {
		return err
	}

	keyVals, err := KeyValsDiff(p, product)
	if err != nil {
		return err
	}

	var (
		columnValues []string
		args         []interface{}
	)
	var i int
	for key, val := range keyVals {
		columnValues = append(columnValues, fmt.Sprintf("%v = $%d", key, i+1))
		args = append(args, val)
		i++
	}

	query := fmt.Sprintf("update %s set %s where id = %s", r.table, strings.Join(columnValues, ", "), p.ID)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		return err
	}
	return nil
}
