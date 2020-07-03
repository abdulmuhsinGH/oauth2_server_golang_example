package server

import (
	"log"
	"net/http"
	"oauth2-server/pkg/auth"
	"oauth2-server/pkg/clientstore"
	"oauth2-server/pkg/cors"
	"os"

	"github.com/go-redis/redis"
	oredis "gopkg.in/go-oauth2/redis.v3"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-session/session"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
)

var (
	manager *manage.Manager
	logger  *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "oauth_server \n", log.LstdFlags|log.Lshortfile)
}

/*
New instantiate a new server
*/
func New(router *mux.Router) {

	logger.Println("AuthServer", "Server is running at 9096 port.")

	logger.Fatal(http.ListenAndServe(":9096", cors.CORS(router)))
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/auth/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(uuid.UUID).String()
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

/*
Oauth sets up and return the oauth2 server
*/
func Oauth(clientStore *clientstore.ClientStore) *server.Server {
	manager = manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr:     os.Getenv("REDIS_SERVER_HOST") + ":" + os.Getenv("REDIS_SERVER_PORT"),
		Password: os.Getenv("REDIS_SERVER_PASS"),
		DB:       15,
	}))

	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(os.Getenv("JWT_SECRET")), jwt.SigningMethodHS512))

	//clientStore.Create
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)

	/* srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		user, err := authService.Login(username, password)
		if err != nil {
			return "", err
		}
		return user.ID.String(), nil
	}) */
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetClientInfoHandler(func(r *http.Request) (clientID, clientSecret string, err error) {
		clientID = r.FormValue("client_id")
		clientSecret = r.FormValue("client_secret")

		if clientID == "" || clientSecret == "" {
			err = errors.ErrAccessDenied
			return
		}

		return
	})

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		logger.Println("Internal Error:", err.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		logger.Println("Response Error:", re.Error.Error())
	})

	return srv
}
/* 
SetPasswordAuthorizationHandler sets the function that handles login attepmts to the oauth server
*/
func SetPasswordAuthorizationHandler(srv *server.Server, service auth.Service) *server.Server {
	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		user, err := service.Login(username, password)
		if err != nil {
			return "", err
		}
		return user.ID.String(), nil
	})

	return srv
}
