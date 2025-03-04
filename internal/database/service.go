package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type service struct {
	db       *sql.DB
	settings Config
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("DB down: %v", err)

		fmt.Printf("DB down: %v\n", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	DBStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(DBStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(DBStats.InUse)
	stats["idle"] = strconv.Itoa(DBStats.Idle)
	stats["wait_count"] = strconv.FormatInt(DBStats.WaitCount, 10)
	stats["wait_duration"] = DBStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(DBStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(DBStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if DBStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if DBStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if DBStats.MaxIdleClosed > int64(DBStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if DBStats.MaxLifetimeClosed > int64(DBStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	fmt.Printf("Disconnected from database: %s\n", s.settings.Database)
	return s.db.Close()
}
