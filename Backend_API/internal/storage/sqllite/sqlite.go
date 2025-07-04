package sqllite

import (
	"database/sql"
	"fmt"
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
			shopID string, 
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
func (sqlDb *DbConnection) CreateProduct(name string, price float32, quantity int, shopID string) (int64, error) {
	stmt, err := sqlDb.db.Prepare(`
		INSERT INTO product (shopID, product_name, product_price, product_quantity)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(shopID, name, price, quantity)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetAllProduct method
func (sqlDb *DbConnection) GetAllProduct(shopID string, offset int64, limit int64, sortOrder string, sortField string) ([]types.Product, error) {
	query := fmt.Sprintf(`
		SELECT id, shopID, product_name, product_price, product_quantity
		FROM product
		WHERE shopID = ?
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, sortField, sortOrder)

	rows, err := sqlDb.db.Query(query, shopID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(&p.Id, &p.ShopID, &p.Name, &p.Price, &p.Quantity); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (sqlDb *DbConnection) GetUserById(id int64, shopID string) (types.Product, error) {
	stmt, err := sqlDb.db.Prepare("SELECT id, shopID, product_name, product_price, product_quantity FROM product WHERE id = ? AND shopID = ?")
	if err != nil {
		return types.Product{}, err
	}
	defer stmt.Close()

	var product types.Product
	err = stmt.QueryRow(id, shopID).Scan(&product.Id, &product.ShopID, &product.Name, &product.Price, &product.Quantity)
	if err != nil {
		return types.Product{}, err
	}

	return product, nil
}

func (sqlDb *DbConnection) DeleteUser(id int64, shopID string) error {

	// Then, delete the product
	stmt, err := sqlDb.db.Prepare("DELETE FROM product WHERE id = ? AND shopID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, shopID)
	if err != nil {
		return err
	}

	return nil
}

func (sqlDb *DbConnection) UpdateProduct(id int64, name string, price float32, quantity int, shopID string) (int64, error) {
	stmt, err := sqlDb.db.Prepare(`UPDATE product
		SET product_name = ?, product_price = ?, product_quantity = ?
		WHERE id = ? AND shopID = ?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, price, quantity, id, shopID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}


