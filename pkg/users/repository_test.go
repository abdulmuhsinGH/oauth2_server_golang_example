package users

import (
	postgres "oauth2-server/pkg/storage/postgres"
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
)

var dbTest *pg.DB
var userRepositoryTest Repository

func init() {
	godotenv.Load(os.ExpandEnv("$GOPATH/src/oauth2-server/.env"))
}

func setupTestCase(t *testing.T, db *pg.DB) func(t *testing.T) {
	t.Log("setup test case")

	err := db.Insert(&User{
		Firstname: "test first name",
		Lastname:  "test last name",
		Password:  "test password",
		Gender:    "male",
		Username:  "test.username",
		Email:     "name@mailcom",
	})

	if err != nil {
		t.Errorf("Test Failed; Could not insert  user seed data: \n" + err.Error())
	}

	return func(t *testing.T) {
		t.Log("teardown test case")
		_, err = db.Model((*User)(nil)).Exec(`TRUNCATE TABLE ?TableName`)
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
func TestAddUser(t *testing.T) {

	dbTest := pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
	}

	defer dbTest.Close()
	userRepositoryTest = NewRepository(dbTest)
	// var user User

	teardownTestCase := setupTestCase(t, dbTest)
	defer teardownTestCase(t)
	user := &User{
		Firstname: "a",
		Lastname:  "b",
		Password:  "c",
		Gender:    "d",
		Username:  "w",
		Email:     "a",
	}
	err := userRepositoryTest.AddUser(user)
	if err != nil {
		t.Errorf("Test Failed; Users was not added")
	}

}

func TestAddUserWithoutEmail(t *testing.T) {

	dbTest = pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
	}
	defer dbTest.Close()

	teardownTestCase := setupTestCase(t, dbTest)
	defer teardownTestCase(t)

	userRepositoryTest = NewRepository(dbTest)
	user := &User{
		Firstname: "a",
		Lastname:  "b",
		Password:  "c",
		Gender:    "d",
		Username:  "qwerty",
	}

	err := userRepositoryTest.AddUser(user)
	if err == nil {
		t.Errorf("Test Failed; Users added. User Added Without Email")
	}

}

func TestAddUserWithoutRole(t *testing.T) {

	dbTest = pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
	}
	defer dbTest.Close()

	teardownTestCase := setupTestCase(t, dbTest)
	defer teardownTestCase(t)

	userRepositoryTest = NewRepository(dbTest)
	user := &User{
		Firstname: "a",
		Lastname:  "b",
		Password:  "c",
		Gender:    "d",
		Username:  "asdf",
	}

	err := userRepositoryTest.AddUser(user)
	if err == nil {
		t.Errorf("Test Failed; Users added. User Added Without Role")
	}

}

func TestAddUserWithEmptyUserObj(t *testing.T) {

	dbTest = pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
	}
	defer dbTest.Close()

	teardownTestCase := setupTestCase(t, dbTest)
	defer teardownTestCase(t)

	userRepositoryTest = NewRepository(dbTest)

	err := userRepositoryTest.AddUser(&User{})
	if err == nil {
		t.Errorf("Test Failed; Users added. Empty user added")
	}

}

func TestGetAllusers(t *testing.T) {

	dbTest = pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
	}
	defer dbTest.Close()
	userRepositoryTest = NewRepository(dbTest)

	teardownTestCase := setupTestCase(t, dbTest)
	defer teardownTestCase(t)

	users, err := userRepositoryTest.GetAllUsers()
	if err != nil {
		t.Errorf("Test Failed; No users found")
	}
	if len(users) != 1 {
		t.Fatalf("Test Failed; Number of users expected %v; actual %v", 1, len(users))
	}
}

func TestFindUserByUsername(t *testing.T) {

	dbTest = pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()
	userRepositoryTest = NewRepository(dbTest)

	teardownTestCase := setupTestCase(t, dbTest)
	defer teardownTestCase(t)

	user, err := userRepositoryTest.FindUserByUsername("test.username")
	if err != nil {
		t.Errorf("Test Failed; No users found")
	}
	if user.Username != "test.username" {
		t.Fatalf("Test Failed; username expected %v; actual %v", "test.username", user.Username)
	}
}

func TestFindUserByUsernameNoUsername(t *testing.T) {
	dbTest = pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()
	userRepositoryTest = NewRepository(dbTest)

	teardownTestCase := setupTestCase(t, dbTest)
	defer teardownTestCase(t)

	_, err := userRepositoryTest.FindUserByUsername("test")
	if err == nil {
		t.Errorf("Test Failed; users found")
	}
}

func TestFindOrAddEmptyUser(t *testing.T) {
	dbTest = pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()
	userRepositoryTest = NewRepository(dbTest)

	teardownTestCase := setupTestCase(t, dbTest)
	defer teardownTestCase(t)

	_, err := userRepositoryTest.FindOrAddUser(&User{})
	if err == nil {
		t.Errorf("Test Failed; user added")
	}
}

func TestFindOrAddUser(t *testing.T) {
	dbTest = pgConnect()
	if dbTest == nil {
		t.Errorf("Test Failed; DB Connection failed")
		return
	}
	defer dbTest.Close()
	userRepositoryTest = NewRepository(dbTest)

	teardownTestCase := setupTestCase(t, dbTest)
	defer teardownTestCase(t)

	user, err := userRepositoryTest.FindOrAddUser(&User{
		Firstname: "test2 first name",
		Lastname:  "test2 last name",
		Password:  "test2 password",
		Gender:    "male",
		Username:  "test2.username",
		Email:     "name2@mailcom",
	})
	if err != nil {
		t.Errorf("Test Failed; user not added")
	}

	if user.Username != "test2.username" {
		t.Fatalf("Test Failed; Wrong user added %v; actual %v", "test2.username", user.Username)
	}

}
