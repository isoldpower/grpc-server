package storage

import (
	"context"
	"database/sql"
	"fmt"
	"golang-grpc/internal/color"
	"golang-grpc/internal/database"
	"golang-grpc/internal/log"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/orders/types"
	"strings"
	"time"

	databaseTypes "golang-grpc/internal/database/types"
)

type PostgresDatabase struct {
	process  databaseTypes.Service
	database *database.Database
	timeout  time.Duration
}

func getItemFromRows(rows *sql.Rows) (*orders.Order, error) {
	item := &orders.Order{}
	if err := rows.Scan(&item.ID, &item.CustomerID, &item.Quantity, &item.ProductID); err != nil {
		log.PrintError("Error scanning row", err)
		return nil, err
	}

	return item, nil
}

func (p *PostgresDatabase) AddItem(item *orders.Order) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	query := `INSERT INTO orders (customerId, productId, quantity) VALUES ($1, $2, $3) RETURNING id`

	if err := p.database.Database.QueryRowContext(
		ctx,
		query,
		item.CustomerID,
		item.ProductID,
		item.Quantity,
	).Scan(&item.ID); err == nil {
		log.Debugln("AddItem finished, generated ID: %s", item.ID)

		return nil, true
	} else {
		return err, false
	}
}

func (p *PostgresDatabase) RemoveItem(id string) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	query := `DELETE FROM orders WHERE id = $1`

	if result, err := p.database.Database.ExecContext(
		ctx,
		query,
		id,
	); err == nil {
		affected, _ := result.RowsAffected()
		log.Debugln("RemoveItem finished with %d affected rows", affected)

		return nil, true
	} else {
		return err, false
	}
}

func (p *PostgresDatabase) ListItems(limit *uint64, offset *uint64) ([]*orders.Order, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	options := make([]string, 0)
	if limit != nil {
		options = append(options, fmt.Sprintf("LIMIT %d", *limit))
	}
	if offset != nil {
		options = append(options, fmt.Sprintf("OFFSET %d", *offset))
	}

	query := fmt.Sprintf("SELECT * FROM orders %s", strings.Join(options, " "))
	rows, err := p.database.Database.QueryContext(ctx, query)
	log.Debugln("ListItems query: %s", color.Yellow(query))
	if err != nil {
		log.Errorln("Error querying database: %v", err)
		return nil, false
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Errorln("Error closing rows: %v", closeErr)
		}
	}()

	items := []*orders.Order{}
	for rows.Next() {
		if item, err := getItemFromRows(rows); err != nil {
			log.PrintError("Error getting item from rows, skipping row", err)
		} else {
			items = append(items, item)
		}
	}

	log.Debugln("Retrieved %d items from database", len(items))
	return items, true
}

func (p *PostgresDatabase) UpdateItem(item *orders.Order) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	var setClauses []string
	var args []interface{}

	if item.CustomerID != "" {
		setClauses = append(setClauses, fmt.Sprintf("customerId = $%d", len(setClauses)))
		args = append(args, item.CustomerID)
	}
	if item.Quantity != 0 {
		setClauses = append(setClauses, fmt.Sprintf("quantity = $%d", len(setClauses)))
		args = append(args, item.Quantity)
	}
	if item.ProductID != "" {
		setClauses = append(setClauses, fmt.Sprintf("productId = $%d", len(setClauses)))
		args = append(args, item.ProductID)
	}

	args = append(args, item.ID)
	query := fmt.Sprintf(
		"UPDATE orders SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		len(setClauses)+1,
	)

	result, err := p.database.Database.ExecContext(ctx, query, args...)
	if err != nil {
		log.PrintError(fmt.Sprintf("Error updating item with id %s", item.ID), err)
		return err, false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.PrintError("Error checking rows affected", err)
		return err, false
	}
	if rowsAffected == 0 {
		log.Errorln("No item found with id %s to update", item.ID)
		return fmt.Errorf("no item found with id %s", item.ID), false
	}

	log.Debugf("Successfully updated item %s, %d rows affected", item.ID, rowsAffected)
	return nil, true
}

func (p *PostgresDatabase) OverwriteItem(item *orders.Order) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := `
        INSERT INTO orders (id, customerId, quantity, productId)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id)
        DO UPDATE SET
            customerId = EXCLUDED.customerId,
            quantity = EXCLUDED.quantity,
            productId = EXCLUDED.productId
    `

	if _, err := p.database.Database.ExecContext(ctx, query, item.ID, item.CustomerID, item.Quantity, item.ProductID); err != nil {
		log.PrintError(fmt.Sprintf("Error overwriting item %s", item.ID), err)
		return err, false
	}

	log.Debugf("Successfully overwrote item %s", item.ID)
	return nil, true
}

// Products
type ProductPostgreSQLStorage struct {
	database *database.Database
	timeout  time.Duration
}

func getProductFromRows(rows *sql.Rows) (*orders.Product, error) {
	item := &orders.Product{}
	if err := rows.Scan(&item.ID, &item.Title, &item.Description); err != nil {
		log.PrintError("Error scanning product row", err)
		return nil, err
	}

	return item, nil
}

func (p *ProductPostgreSQLStorage) AddItem(item *orders.Product) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	query := `INSERT INTO products (title, description) VALUES ($1, $2) RETURNING id`

	if err := p.database.Database.QueryRowContext(
		ctx,
		query,
		item.Title,
		item.Description,
	).Scan(&item.ID); err == nil {
		log.Debugln("AddProduct finished, generated ID: %s", item.ID)

		return nil, true
	} else {
		return err, false
	}
}

func (p *ProductPostgreSQLStorage) RemoveItem(id string) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	query := `DELETE FROM products WHERE id = $1`

	if result, err := p.database.Database.ExecContext(ctx, query, id); err == nil {
		affected, _ := result.RowsAffected()
		log.Debugln("RemoveProduct finished with %d affected rows", affected)
		return nil, true
	} else {
		return err, false
	}
}

func (p *ProductPostgreSQLStorage) GetItem(id string) (*orders.Product, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := "SELECT * FROM products WHERE id = $1 LIMIT 1"
	rows, err := p.database.Database.QueryContext(ctx, query, id)
	if err != nil {
		log.Errorln("Error querying database: %v", err)
		return nil, false
	}
	defer rows.Close()

	if rows.Next() {
		if item, err := getProductFromRows(rows); err != nil {
			log.PrintError("Error getting product from rows", err)
			return nil, false
		} else {
			return item, true
		}
	}
	return nil, false
}

func (p *ProductPostgreSQLStorage) ListItems(limit *uint64, offset *uint64) ([]*orders.Product, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	options := make([]string, 0)
	if limit != nil {
		options = append(options, fmt.Sprintf("LIMIT %d", *limit))
	}
	if offset != nil {
		options = append(options, fmt.Sprintf("OFFSET %d", *offset))
	}

	query := fmt.Sprintf("SELECT * FROM products %s", strings.Join(options, " "))
	rows, err := p.database.Database.QueryContext(ctx, query)
	if err != nil {
		log.Errorln("Error querying database: %v", err)
		return nil, false
	}
	defer rows.Close()

	items := []*orders.Product{}
	for rows.Next() {
		if item, err := getProductFromRows(rows); err != nil {
			log.PrintError("Error getting product from rows, skipping row", err)
		} else {
			items = append(items, item)
		}
	}

	return items, true
}

func (p *ProductPostgreSQLStorage) CountItems() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := "SELECT COUNT(*) FROM products"
	var count uint64
	err := p.database.Database.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Errorln("Error counting products: %v", err)
		return 0, err
	}

	return count, nil
}

func (p *ProductPostgreSQLStorage) UpdateItem(item *orders.Product) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := `UPDATE products SET title = $1, description = $2 WHERE id = $3`
	if _, err := p.database.Database.ExecContext(ctx, query, item.Title, item.Description, item.ID); err != nil {
		log.PrintError(fmt.Sprintf("Error updating product %s", item.ID), err)
		return err, false
	}

	return nil, true
}

func (p *ProductPostgreSQLStorage) OverwriteItem(item *orders.Product) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := `
        INSERT INTO products (id, title, description)
        VALUES ($1, $2, $3)
        ON CONFLICT (id)
        DO UPDATE SET
            title = EXCLUDED.title,
            description = EXCLUDED.description
    `

	if _, err := p.database.Database.ExecContext(ctx, query, item.ID, item.Title, item.Description); err != nil {
		log.PrintError(fmt.Sprintf("Error overwriting product %s", item.ID), err)
		return err, false
	}

	return nil, true
}

// Customers
type CustomerPostgreSQLStorage struct {
	database *database.Database
	timeout  time.Duration
}

func getCustomerFromRows(rows *sql.Rows) (*orders.Customer, error) {
	item := &orders.Customer{}
	if err := rows.Scan(&item.ID, &item.Name); err != nil {
		log.PrintError("Error scanning customer row", err)
		return nil, err
	}

	return item, nil
}

func (p *CustomerPostgreSQLStorage) AddItem(item *orders.Customer) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	query := `INSERT INTO customers (name) VALUES ($1) RETURNING id`

	if err := p.database.Database.QueryRowContext(
		ctx,
		query,
		item.Name,
	).Scan(&item.ID); err == nil {
		log.Debugln("AddCustomer finished, generated ID: %s", item.ID)

		return nil, true
	} else {
		return err, false
	}
}

func (p *CustomerPostgreSQLStorage) RemoveItem(id string) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	query := `DELETE FROM customers WHERE id = $1`

	if result, err := p.database.Database.ExecContext(ctx, query, id); err == nil {
		affected, _ := result.RowsAffected()
		log.Debugln("RemoveCustomer finished with %d affected rows", affected)
		return nil, true
	} else {
		return err, false
	}
}

func (p *CustomerPostgreSQLStorage) GetItem(id string) (*orders.Customer, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := "SELECT * FROM customers WHERE id = $1 LIMIT 1"
	rows, err := p.database.Database.QueryContext(ctx, query, id)
	if err != nil {
		log.Errorln("Error querying database: %v", err)
		return nil, false
	}
	defer rows.Close()

	if rows.Next() {
		if item, err := getCustomerFromRows(rows); err != nil {
			log.PrintError("Error getting customer from rows", err)
			return nil, false
		} else {
			return item, true
		}
	}
	return nil, false
}

func (p *CustomerPostgreSQLStorage) ListItems(limit *uint64, offset *uint64) ([]*orders.Customer, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	options := make([]string, 0)
	if limit != nil {
		options = append(options, fmt.Sprintf("LIMIT %d", *limit))
	}
	if offset != nil {
		options = append(options, fmt.Sprintf("OFFSET %d", *offset))
	}

	query := fmt.Sprintf("SELECT * FROM customers %s", strings.Join(options, " "))
	rows, err := p.database.Database.QueryContext(ctx, query)
	if err != nil {
		log.Errorln("Error querying database: %v", err)
		return nil, false
	}
	defer rows.Close()

	items := []*orders.Customer{}
	for rows.Next() {
		if item, err := getCustomerFromRows(rows); err != nil {
			log.PrintError("Error getting customer from rows, skipping row", err)
		} else {
			items = append(items, item)
		}
	}

	return items, true
}

func (p *CustomerPostgreSQLStorage) CountItems() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := "SELECT COUNT(*) FROM customers"
	var count uint64
	err := p.database.Database.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Errorln("Error counting customers: %v", err)
		return 0, err
	}

	return count, nil
}

func (p *CustomerPostgreSQLStorage) UpdateItem(item *orders.Customer) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := `UPDATE customers SET name = $1 WHERE id = $2`
	if _, err := p.database.Database.ExecContext(ctx, query, item.Name, item.ID); err != nil {
		log.PrintError(fmt.Sprintf("Error updating customer %s", item.ID), err)
		return err, false
	}

	return nil, true
}

func (p *CustomerPostgreSQLStorage) OverwriteItem(item *orders.Customer) (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := `
        INSERT INTO customers (id, name)
        VALUES ($1, $2)
        ON CONFLICT (id)
        DO UPDATE SET
            name = EXCLUDED.name
    `

	if _, err := p.database.Database.ExecContext(ctx, query, item.ID, item.Name); err != nil {
		log.PrintError(fmt.Sprintf("Error overwriting customer %s", item.ID), err)
		return err, false
	}

	return nil, true
}

func (p *PostgresDatabase) GetItem(id string) (*orders.Order, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := "SELECT * FROM orders WHERE id = $1 LIMIT 1"
	rows, err := p.database.Database.QueryContext(ctx, query, id)
	if err != nil {
		log.Errorln("Error querying database: %v", err)
		return nil, false
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Errorln("Error closing rows: %v", closeErr)
		}
	}()

	if item, err := getItemFromRows(rows); err != nil {
		log.PrintError("Error getting item from rows", err)
		return nil, false
	} else {
		return item, true
	}
}

func (p *PostgresDatabase) CountItems() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := "SELECT COUNT(*) FROM orders"
	var count uint64
	err := p.database.Database.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Errorln("Error counting items: %v", err)
		return 0, err
	}

	return count, nil
}

func NewPostgresStorageTimed(
	config *database.Config,
	timeout time.Duration,
) types.ObjectStore[orders.Order] {
	instance := database.NewDatabase(config)
	process, err := instance.Instantiate()
	if err != nil {
		panic(err)
	}

	return &PostgresDatabase{
		database: instance,
		process:  process,
		timeout:  timeout,
	}
}

func NewPostgresStorage(config *database.Config) types.ObjectStore[orders.Order] {
	instance := database.NewDatabase(config)
	_, err := instance.Instantiate()
	if err != nil {
		panic(err)
	}

	return &PostgresDatabase{
		database: instance,
		timeout:  3 * time.Second,
	}
}

func NewCustomerPostgresStorage(config *database.Config) types.ObjectStore[orders.Customer] {
	instance := database.NewDatabase(config)
	_, err := instance.Instantiate()
	if err != nil {
		panic(err)
	}

	return &CustomerPostgreSQLStorage{
		database: instance,
		timeout:  3 * time.Second,
	}
}

func NewProductPostgresStorage(config *database.Config) types.ObjectStore[orders.Product] {
	instance := database.NewDatabase(config)
	_, err := instance.Instantiate()
	if err != nil {
		panic(err)
	}

	return &ProductPostgreSQLStorage{
		database: instance,
		timeout:  3 * time.Second,
	}
}
