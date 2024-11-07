package records

import (
	"time"
)

type Product struct {
	Id          int        `db:"product_id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	Price       int        `db:"price"`
	Stock       int        `db:"stock"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
