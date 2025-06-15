package sqllite

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/slangeres/Vypaar/backend_API/internal/config"
	"github.com/slangeres/Vypaar/backend_API/internal/types"
)

type DbConnection struct {
	db *sql.DB
}

func ConfigSQL(cnf *config.Config) (*DbConnection, error) {
	db, err := sql.Open("sqlite3", cnf.StoragePath)
	if err != nil {
		slog.Error("DB ERROR: Unable to connect to the storage path", "error", err)
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS product (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_name TEXT,
			product_price REAL,
			product_quantity INTEGER
		)
	`)
	if err != nil {
		slog.Error("DB ERROR: Unable to create the table in db", "error", err)
		return nil, err
	}

	return &DbConnection{db: db}, nil
}

// CreateProduct method
func (sqlDb *DbConnection) CreateProduct(name string, price float32, quantity int) (int64, error) {
	stmt, err := sqlDb.db.Prepare(`
		INSERT INTO product (product_name, product_price, product_quantity)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, price, quantity)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetAllProduct method
func (sqlDb *DbConnection) GetAllProduct() ([]types.Product, error) {
	rows, err := sqlDb.db.Query("SELECT id, product_name, product_price, product_quantity FROM product")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Quantity); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
func (sqlDb *DbConnection) GetUserById(id int64) (types.Product, error) {
	stmt, err := sqlDb.db.Prepare("SELECT id, product_name, product_price, product_quantity FROM product WHERE id = ?")
	if err != nil {
		return types.Product{}, err
	}
	defer stmt.Close()

	var product types.Product
	err = stmt.QueryRow(id).Scan(&product.Id, &product.Name, &product.Price, &product.Quantity)
	if err != nil {
		return types.Product{}, err
	}

	return product, nil
}

func (sqlDb *DbConnection) DeleteUser(id int64) error {

	// Then, delete the product
	stmt, err := sqlDb.db.Prepare("DELETE FROM product WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (sqlDb *DbConnection) UpdateProduct(id int64, name string, price float32, quantity int) (int64, error) {
	stmt, err := sqlDb.db.Prepare(`UPDATE product
		SET product_name = ?, product_price = ?, product_quantity = ?
		WHERE id = ?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, price, quantity, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
