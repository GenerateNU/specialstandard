package config

import (
	"database/sql"
	"fmt"
	"time"
)

type DB struct {
	Host     string `env:"DB_HOST, required"`
	Port     string `env:"DB_PORT, required"`
	User     string `env:"DB_USER, required"`
	Password string `env:"DB_PASSWORD, required"`
	Name     string `env:"DB_NAME, required"`
	// Add these:
	MaxOpenConns    int `env:"DB_MAX_OPEN_CONNS" envDefault:"50"`         // max connections to keep open
	MaxIdleConns    int `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`          // max idle connections
	MaxConnIdleTime int `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"100000"` // max idle time in seconds
	ConnMaxLifetime int `env:"DB_CONN_MAX_LIFETIME"`                      // connection lifetime in seconds
}

func (db *DB) Connection() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		db.Host, db.User, db.Password, db.Name, db.Port)
}

// New method to configure the pool
func (db *DB) ConfigurePool(sqlDB *sql.DB) {
	// Set sensible defaults if not configured
	maxOpen := db.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 25 // start conservative
	}

	maxIdle := db.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 5
	}

	lifetime := db.ConnMaxLifetime
	if lifetime == 0 {
		lifetime = 5 * 60 // 5 minutes
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(lifetime) * time.Second)
}
