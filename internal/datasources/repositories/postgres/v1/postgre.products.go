package v1

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/datasources/records"
)

type postgreProductRepository struct {
	conn *sqlx.DB
}

func NewProductRepository(conn *sqlx.DB) V1Domains.ProductRepository {
	return &postgreProductRepository{
		conn: conn,
	}
}

func (r *postgreProductRepository) Store(ctx context.Context, p *V1Domains.ProductDomain) (V1Domains.ProductDomain, error) {
	query := `
		INSERT INTO products (name, description, price, stock, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING product_id, name, description, price, stock, created_at, updated_at
	`
	var result records.Product
	now := time.Now()
	err := r.conn.GetContext(ctx, &result, query, p.Name, p.Description, p.Price, p.Stock, now)
	if err != nil {
		return V1Domains.ProductDomain{}, err
	}

	return result.ToV1Domain(), nil
}

func (r *postgreProductRepository) GetAll(ctx context.Context) ([]V1Domains.ProductDomain, error) {
	query := `SELECT product_id, name, description, price, stock, created_at, updated_at FROM products`
	var productsFromDB []records.Product
	err := r.conn.SelectContext(ctx, &productsFromDB, query)
	if err != nil {
		return nil, err
	}

	var convertedProducts []V1Domains.ProductDomain
	for _, val := range productsFromDB {
		convertedProducts = append(convertedProducts, val.ToV1Domain())
	}

	return convertedProducts, nil
}

func (r *postgreProductRepository) GetById(ctx context.Context, id int) (V1Domains.ProductDomain, error) {
	query := `SELECT product_id, name, description, price, stock, created_at, updated_at FROM products WHERE product_id = $1`
	var product records.Product
	err := r.conn.GetContext(ctx, &product, query, id)
	if err != nil {
		return V1Domains.ProductDomain{}, err
	}

	return product.ToV1Domain(), nil
}

func (r *postgreProductRepository) Update(ctx context.Context, p *V1Domains.ProductDomain) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, updated_at = $5
		WHERE product_id = $6
	`
	now := time.Now()
	_, err := r.conn.ExecContext(ctx, query, p.Name, p.Description, p.Price, p.Stock, now, p.Id)
	return err
}

func (r *postgreProductRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE product_id = $1`
	_, err := r.conn.ExecContext(ctx, query, id)
	return err
}
