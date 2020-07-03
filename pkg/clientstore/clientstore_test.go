package clientstore

import (
	"oauth2-server/pkg/storage/postgres"
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(os.ExpandEnv("$GOPATH/src/oauth2-server/.env"))
}

func setupTestCase(t *testing.T, db *pg.DB) func(t *testing.T) {
	t.Log("setup test case")

	err := db.Insert(&OauthClient{
		ID:     "test_id",
		Secret: "test last name",
		Domain: "http://127.0.0.1:1234",
	})

	if err != nil {
		t.Errorf("Test Failed; Could not insert  user seed data: \n" + err.Error())
	}

	return func(t *testing.T) {
		t.Log("teardown test case")
		_, err = db.Model((*OauthClient)(nil)).Exec(`TRUNCATE TABLE ?TableName`)
		if err != nil {
			t.Errorf("Test Failed; Users Table truncate failed")
		}

	}
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

func pgConnect() *pg.DB {
	dbInfo := pgOptions()
	return postgres.Connect(dbInfo)
}

func TestCreateOauthClient(t *testing.T) {
	dbTest := pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()

	cs := New(dbTest)
	err := cs.Create(OauthClient{
		ID:     "1234",
		Secret: "qwedfvbhytrew",
		Domain: "http://127.0.0.1:1234",
	})

	if err != nil {
		t.Errorf("Test Failed; Unable to create OauthClient: %v", err.Error())
	}
}

func TestOauthClientGetByID(t *testing.T) {
	dbTest := pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()

	tearDownTestCase := setupTestCase(t, dbTest)
	defer tearDownTestCase(t)

	cs := New(dbTest)
	ci, err := cs.GetByID("test_id")
	t.Log(ci)

	if err != nil {
		t.Errorf("Test Failed; OauthClient Not found: %v", err.Error())
	}

	if ci.GetDomain() != "http://127.0.0.1:1234" {
		t.Errorf("Test Failed; OauthClient not found; expected %v; actual: %v", "http://127.0.0.1:1234", ci.GetDomain())
	}
}

func TestOauthClientGetByIDNoClient(t *testing.T) {
	dbTest := pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()

	tearDownTestCase := setupTestCase(t, dbTest)
	defer tearDownTestCase(t)

	cs := New(dbTest)
	ci, err := cs.GetByID("test")
	t.Log(ci)

	if err == nil {
		t.Errorf("Test Failed; OauthClient Not found: %v", err.Error())
	}

	if ci != nil {
		t.Errorf("Test Failed; OauthClient not found; expected %v; actual: %v", nil, ci)
	}
}
