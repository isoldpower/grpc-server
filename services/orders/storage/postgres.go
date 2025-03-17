package storage

import (
	"context"
	"fmt"
	"golang-grpc/internal/database"
	"golang-grpc/internal/log"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/orders/types"
	"time"

	databaseTypes "golang-grpc/internal/database/types"
)

type PostgresDatabase struct {
	process  databaseTypes.Service
	database *database.Database
}

func (p *PostgresDatabase) AddItem(item *orders.Order) (error, bool) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	query := "INSERT INTO orders (customerId, productId, quantity) VALUES ($1, $2, $3)"

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
	log.Infoln("Remove from database")
	// Remove from database
	return nil, true
}

func (p *PostgresDatabase) ListItems(from int, to int) ([]*orders.Order, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(
		"SELECT * FROM orders LIMIT %d OFFSET %d",
		to-from+1,
		from,
	)

	rows, err := p.database.Database.QueryContext(ctx, query)
	if err != nil {
		log.Errorln("Error querying database: %v", err)
		return nil, false
	}
	defer rows.Close()

	var items []*orders.Order
	for rows.Next() {
		item := orders.Order{}
		if err := rows.Scan(&item.ID, &item.CustomerID, &item.Quantity, &item.ProductID); err != nil {
			log.PrintError("Error scanning row", err)
			return nil, false
		}

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		log.Errorln("Error during row iteration: %v", err)
		return nil, false
	}

	log.Debugln("Retrieved %d items from database", len(items))
	return items, true
}

func (p *PostgresDatabase) UpdateItem(item *orders.Order) (error, bool) {
	log.Infoln("Update item")
	// Update item
	return nil, true
}

func (p *PostgresDatabase) OverwriteItem(item *orders.Order) (error, bool) {
	log.Infoln("Overwrite item")
	// Overwrite item
	return nil, true
}

func (p *PostgresDatabase) GetItem(id string) (*orders.Order, bool) {
	log.Infoln("Get Item")
	// Get Item
	return nil, true
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
	}
}
