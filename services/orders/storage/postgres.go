package storage

import (
	"context"
	"fmt"
	"golang-grpc/internal/database"
	"golang-grpc/internal/log"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/orders/types"
	"strconv"
	"time"

	databaseTypes "golang-grpc/internal/database/types"
)

type IndexedOrder struct {
	*orders.Order
}

func (io IndexedOrder) GetID() string {
	return strconv.Itoa(int(io.ID))
}

type PostgresDatabase[T types.IndexedStoreItem] struct {
	process  databaseTypes.Service
	database *database.Database
}

func (p *PostgresDatabase[T]) AddItem(item *T) (error, bool) {
	var order IndexedOrder
	if parsed, ok := any(*item).(IndexedOrder); ok == false {
		return fmt.Errorf("item is not IndexedOrder"), false
	} else {
		order = parsed
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	query := fmt.Sprintf(
		"INSERT INTO %s (customerId, productId, quantity) VALUES (%d, %d, %d)",
		"orders",
		order.GetCustomerID(),
		order.GetProductID(),
		order.GetQuantity(),
	)

	if result, err := p.database.Database.ExecContext(ctx, query); err == nil {
		affected, _ := result.RowsAffected()
		log.Debugln("AddItem finished with %d affected rows", affected)
		return nil, true
	} else {
		return err, false
	}
}

func (p *PostgresDatabase[T]) RemoveItem(id string) (error, bool) {
	log.Infoln("Remove from database")
	// Remove from database
	return nil, true
}

func (p *PostgresDatabase[T]) ListItems(from int, to int) ([]*T, bool) {
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
	var items []IndexedOrder
	for rows.Next() {
		item := IndexedOrder{Order: &orders.Order{}}
		err := rows.Scan(&item.Order.ID, &item.CustomerID, &item.Quantity, &item.ProductID)
		if err != nil {
			log.Errorln("Error scanning row: %v", err)
			return nil, false
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Errorln("Error during row iteration: %v", err)
		return nil, false
	}

	log.Debugln("Retrieved %d items from database", len(items))
	parsedItems := make([]*T, len(items))
	for i, item := range items {
		parsed := any(item).(T)
		parsedItems[i] = &parsed
	}
	return parsedItems, true
}

func (p *PostgresDatabase[T]) UpdateItem(item *T) (error, bool) {
	log.Infoln("Update item")
	// Update item
	return nil, true
}

func (p *PostgresDatabase[T]) OverwriteItem(item *T) (error, bool) {
	log.Infoln("Overwrite item")
	// Overwrite item
	return nil, true
}

func (p *PostgresDatabase[T]) GetItem(id string) (*T, bool) {
	log.Infoln("Get Item")
	// Get Item
	return nil, true
}

func NewPostgresStorage(config *database.Config) types.IndexedObjectStore[IndexedOrder] {
	instance := database.NewDatabase(config)
	process, err := instance.Instantiate()
	if err != nil {
		panic(err)
	}

	return &PostgresDatabase[IndexedOrder]{
		database: instance,
		process:  process,
	}
}
