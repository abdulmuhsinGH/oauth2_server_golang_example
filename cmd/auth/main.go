package main

import (
	"oauth2-server/pkg/auth"
	"oauth2-server/pkg/clientstore"
	"oauth2-server/pkg/server"
	"oauth2-server/pkg/storage/postgres"
	"oauth2-server/pkg/users"
	"os"

	"github.com/go-pg/pg/v9"
	"github.com/gorilla/mux"
)

func main() {

	var (
		// local db credential
		DbHost     = os.Getenv("DB_HOST")
		DbUser     = os.Getenv("DB_USER")
		DbPassword = os.Getenv("DB_PASS")
		DbPort     = os.Getenv("DB_PORT")
		DbName     = os.Getenv("DB_NAME")
	)

	dbInfo := pg.Options{
		Addr:     DbHost + ":" + DbPort,
		User:     DbUser,
		Password: DbPassword,
		Database: DbName,
	}

	db := postgres.Connect(dbInfo)
	defer db.Close()

	clientStore := clientstore.New(db)
	_ = clientStore.Create(clientstore.OauthClient{
		ID:     os.Getenv("ADMIN_CLIENT_ID"),
		Secret: os.Getenv("ADMIN_CLIENT_SECRET"),
		Domain: os.Getenv("ADMIN_CLIENT_DOMAIN"),
		Data:   nil,
	})

	router := mux.NewRouter()
	userRepository := users.NewRepository(db)

	oauthSrv := server.Oauth(clientStore)

	authHandler := auth.NewHandlers(oauthSrv, clientStore, userRepository)
	authHandler.SetupRoutes(router)

	server.SetPasswordAuthorizationHandler(oauthSrv, authHandler.AuthService)
	go server.New(router)
	auth.Client()
}
