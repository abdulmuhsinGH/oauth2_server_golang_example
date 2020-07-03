package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"oauth2-server/pkg/clientstore"
	"oauth2-server/pkg/format"
	"oauth2-server/pkg/users"
	"os"
	"path/filepath"
	"time"

	"github.com/go-session/session"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/oauth2.v3/server"
)

var ()

/*
Handlers define auth
*/
type Handlers struct {
	AuthService Service
	clientStore *clientstore.ClientStore
	oauthServer *server.Server
	authLogging *log.Logger
}

var googleOauthConfig *oauth2.Config

func (h *Handlers) handlePostLoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	oauthState := h.AuthService.GenerateState(w)
	newURL := googleOauthConfig.AuthCodeURL(oauthState)

	http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
}

func (h *Handlers) handleGoogleAuthCallback(w http.ResponseWriter, r *http.Request) {

	oauthState, err := r.Cookie("oauth-state")
	if err != nil {
		h.authLogging.Println("getting_cookie_err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var decodedState string

	err = SecuredCookie.Decode("oauth-state", oauthState.Value, &decodedState)
	if err != nil {
		h.authLogging.Println("cookie err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}

	store.Set("ReturnUri", form)

	if r.FormValue("state") != decodedState {
		h.authLogging.Println("google_oauth_error", "invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := h.AuthService.GetUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		h.authLogging.Println("google_oauth_get_user_data_error", err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	user, err := h.AuthService.SignUpViaGoogle(data)
	if err != nil {
		h.authLogging.Println("google_oauth_sign_up_users_error", err.Error())
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	store.Set("LoggedInUserID", user.ID)
	store.Save()

	w.Header().Set("Location", "/auth")
	w.WriteHeader(http.StatusFound)
}

func (h *Handlers) handleToken(response http.ResponseWriter, request *http.Request) {
	//request.GetBody()
	err := h.oauthServer.HandleTokenRequest(response, request)
	if err != nil {
		h.authLogging.Printf("handle_token_Error: %v\n", err.Error())
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) handleUserAuthTest(response http.ResponseWriter, request *http.Request) {
	token, err := h.oauthServer.ValidationBearerToken(request)
	if err != nil {
		h.authLogging.Printf("Error: %v\n", err.Error())
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}
	e := json.NewEncoder(response)
	e.SetIndent("", "  ")
	e.Encode(data)
}

/*
HandleAddUser gets data from http request and sends to
*/
func (h *Handlers) handleAuthorize(response http.ResponseWriter, request *http.Request) {
	store, err := session.Start(nil, response, request)
	if err != nil {
		h.authLogging.Printf("Error: %v\n", err.Error())
		format.Send(response, 500, format.Message(false, "Error while starting session", nil))
		return
	}
	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {

		form = v.(url.Values)
	}
	request.Form = form

	store.Delete("ReturnUri")
	store.Save()

	err = h.oauthServer.HandleAuthorizeRequest(response, request)
	if err != nil {
		h.authLogging.Println("HandleAuthorizeRequestError:", err.Error())
		format.Send(response, 500, format.Message(false, "Error handling authorization", nil))
	}
}

func (h *Handlers) handleLogin(w http.ResponseWriter, r *http.Request) {
	outputHTML(w, r, "/login.html")
}

func (h *Handlers) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := h.AuthService.Login(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		h.authLogging.Printf("Error: %v", err.Error())
		format.Send(w, http.StatusUnauthorized, format.Message(false, err.Error(), nil))
		return
	}
	store.Set("LoggedInUserID", user.ID)
	store.Save()

	redirectURI := r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	h.authLogging.Println(redirectURI, clientID)

	w.Header().Set("Location", "/auth")
	w.WriteHeader(http.StatusFound)
}

func (h *Handlers) handleSignUp(w http.ResponseWriter, r *http.Request) {
	outputHTML(w, r, "/signup.html")
}
func (h *Handlers) handlePostSignUp(response http.ResponseWriter, request *http.Request) {
	//body, err := ioutil.ReadAll(request.Body)
	//h.authLogging.Printlog("request_body: ", string(body))
	newUser := users.User{
		Firstname: request.FormValue("firstname"),
		Username:  request.FormValue("username"),
		Email:     request.FormValue("email"),
		Lastname:  request.FormValue("lastname"),
		Gender:    request.FormValue("gender"),
		Password:  request.FormValue("password"),
	}

	err := h.AuthService.SignUp(newUser)
	if err != nil {
		h.authLogging.Printf("Error: %v", err.Error())
		format.Send(response, http.StatusUnauthorized, format.Message(false, err.Error(), nil))
		return
	}

	response.Header().Set("Location", "/auth/login")
	format.Send(response, http.StatusCreated, format.Message(true, "User Created", nil))
}

func (h *Handlers) handleAddClient(response http.ResponseWriter, request *http.Request) {
	oauthClient := clientstore.OauthClient{}
	body, err := ioutil.ReadAll(request.Body)

	err = json.Unmarshal([]byte(body), &oauthClient) //NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		h.authLogging.Printf("Error while decoding request body: %v", err.Error())
		format.Send(response, 500, format.Message(false, "Error while decoding request body", nil))
		return
	}
	err = h.clientStore.Create(oauthClient)
	if err != nil {
		h.authLogging.Printf("Error: %v", err.Error())
		format.Send(response, http.StatusUnauthorized, format.Message(false, err.Error(), nil))
		return
	}
	format.Send(response, http.StatusCreated, format.Message(true, "Client Created", nil))
}

func (h *Handlers) handleAuth(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		h.authLogging.Printf("Error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/auth/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "/auth.html")
}

// outputHTML renders static html files
func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	filePrefix, _ := filepath.Abs("./view/")
	file, err := os.Open(filePrefix + filename)
	if err != nil {

		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

/*
SetupRoutes sets up routes to respective handlers
*/
func (h *Handlers) SetupRoutes(mux *mux.Router) {
	mux.HandleFunc("/auth", h.Httplog(h.handleAuth)).Methods("GET")
	mux.HandleFunc("/auth/login", h.Httplog(h.handleLogin)).Methods("GET")
	mux.HandleFunc("/auth/login", h.Httplog(h.handlePostLogin)).Methods("POST")

	mux.HandleFunc("/auth/signup", h.Httplog(h.handleSignUp)).Methods("GET")
	mux.HandleFunc("/auth/signup", h.Httplog(h.handlePostSignUp)).Methods("POST")

	mux.HandleFunc("/auth/authorize", h.Httplog(h.handleAuthorize)).Methods("GET", "POST")
	mux.HandleFunc("/auth/token", h.Httplog(h.handleToken)).Methods("POST")
	mux.HandleFunc("/auth/test", h.Httplog(ValidateToken(h.handleUserAuthTest, h.oauthServer))).Methods("GET")
	mux.HandleFunc("/auth/google/login", h.Httplog(h.handlePostLoginWithGoogle))
	mux.HandleFunc("/auth/google/callback", h.Httplog(h.handleGoogleAuthCallback))
	mux.HandleFunc("/auth/client", h.Httplog(ValidateToken(h.handleAddClient, h.oauthServer))).Methods("POST")
}

/*
Httplog handles how long it takes for a request to process
*/
func (h *Handlers) Httplog(next http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		defer h.authLogging.Printf("%s request processed in %s\n", request.URL.Path, time.Now().Sub(startTime))
		next(response, request)
	}
}

/*
NewHandlers initiates auth handler
*/
func NewHandlers(oauthServerArg *server.Server, clientStoreArg *clientstore.ClientStore, userRepository users.Repository) *Handlers {

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "profile"},
		Endpoint:     google.Endpoint,
	}
	return &Handlers{
		AuthService: NewAuthService(userRepository),
		oauthServer: oauthServerArg,
		clientStore: clientStoreArg,
		authLogging: log.New(os.Stdout, "auth_handler: ", log.LstdFlags|log.Lshortfile),
	}
}
