package database

import (
	"context"
	"sync"
)

var (
	instance *Database
	once     sync.Once
)

// GetDatabase returns a singleton database instance
func GetDatabase(ctx context.Context, databaseURL string) (*Database, error) {
	var err error
	once.Do(func() {
		instance, err = New(ctx, databaseURL)
	})
	return instance, err
}