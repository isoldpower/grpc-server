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
	query := `INSERT INTO orders (customerId, productId, quantity) VALUES ($1, $2, $3)`

	if result, err := p.database.Database.ExecContext(
		ctx,
		query,
		item.CustomerID,
		item.ProductID,
		item.Quantity,
	); err == nil {
		affected, _ := result.RowsAffected()
		log.Debugln("AddItem finished with %d affected rows", affected)

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

	var items []*orders.Order
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

func (p *PostgresDatabase) GetItem(id string) (*orders.Order, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	query := fmt.Sprintf("SELECT * FROM orders WHERE id = $1 LIMIT 1")
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
	process, err := instance.Instantiate()
	if err != nil {
		panic(err)
	}

	return &PostgresDatabase{
		database: instance,
		process:  process,
		timeout:  3 * time.Second,
	}
}
