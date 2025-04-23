package customer

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

const (
	createQuery     = `INSERT INTO customers (name, email) VALUES ($1, $2) RETURNING id`
	getByEmailQuery = `SELECT id, name, email FROM customers WHERE email = $1`
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	return &Repository{conn: conn}, nil
}

func (r *Repository) CreateCustomer(ctx context.Context, customer Customer) (Customer, error) {
	err := r.conn.QueryRow(ctx, createQuery, customer.Name, customer.Email).Scan(&customer.Id)
	return customer, err
}

func (r *Repository) GetCustomerByEmail(ctx context.Context, email string) (Customer, error) {
	var customer Customer
	err := r.conn.QueryRow(ctx, getByEmailQuery, email).Scan(&customer.Id, &customer.Name, &customer.Email)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}
