package postgres

import (
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(os.ExpandEnv("$GOPATH/src/oauth2-server/.env"))
}

func pgOptions() pg.Options {

	return pg.Options{
		Addr:            os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		User:            os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASS"),
		Database:        os.Getenv("DB_TEST_NAME"),
		MaxRetries:      1,
		MinRetryBackoff: -1,

		DialTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		PoolSize:           10,
		MaxConnAge:         10 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        10 * time.Second,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

func TestPGConnection(t *testing.T) {
	db:= Connect(pgOptions())
	if db == nil{
		t.Errorf("DB connection failed")
	}
}
