package auth

import (
	"oauth2-server/pkg/storage/postgres"
	"oauth2-server/pkg/users"
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

	err := db.Insert(&users.User{
		Firstname: "test first name",
		Lastname:  "test last name",
		Password:  "$2a$14$TH23lPu7kA9QiRqW8SCNJOg182LKQ7okjhCThCN.ICSw9dgmBk2a2",
		Gender:    "male",
		Username:  "test.username",
		Email:     "name@mailcom",
	})

	if err != nil {
		t.Errorf("Test Failed; Could not insert  user seed data: \n" + err.Error())
	}

	return func(t *testing.T) {
		t.Log("teardown test case")
		_, err = db.Model((*users.User)(nil)).Exec(`TRUNCATE TABLE ?TableName`)
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

func TestLogin(t *testing.T) {
	dbTest := pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()

	tearDownTestCase := setupTestCase(t, dbTest)
	defer tearDownTestCase(t)

	authServiceTest := NewAuthService(users.NewRepository(dbTest))

	user, err := authServiceTest.Login("test.username", "secret")
	if err != nil {
		t.Errorf("Test Failed, user failed to login it; error: %v", err.Error())
	}

	if user.Firstname != "test first name" {
		t.Errorf("Test Failed, returned wrong user; expected: %v; actual: %v", "test first name", user.Firstname)
	}
}

func TestLoginWrongPassowrd(t *testing.T) {
	dbTest := pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()

	tearDownTestCase := setupTestCase(t, dbTest)
	defer tearDownTestCase(t)

	authServiceTest := NewAuthService(users.NewRepository(dbTest))

	_, err := authServiceTest.Login("test.username", "qwerty")
	if err == nil {
		t.Errorf("Test Failed, user logged in")
	}
}


func TestLoginWrongUsername(t *testing.T) {
	dbTest := pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()

	tearDownTestCase := setupTestCase(t, dbTest)
	defer tearDownTestCase(t)

	authServiceTest := NewAuthService(users.NewRepository(dbTest))

	_, err := authServiceTest.Login("wrong.username", "secret")
	if err == nil {
		t.Errorf("Test Failed, user logged in")
	}
}
