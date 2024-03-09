package domain

import (
	"context"
	"log"
	"time"
)

const (
	timeout time.Duration = 100 // nanoseconds
)

// DBStatus returns DB connection status.
func (o *LibraryService) DBStatus() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	if err := o.db.PingContext(ctx); err != nil {
		log.Printf("[ERROR]: ping context error: %v", err)
		return false, err
	}
	return true, nil
}
