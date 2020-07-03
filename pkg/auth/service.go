package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"oauth2-server/pkg/users"
	"os"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/oauth2.v3/server"
)

// Service provides user adding operations.
type Service interface {
	Login(string, string) (*users.User, error)
	HashPassword(string) (string, error)
	CheckPasswordHash(string, string) bool
	SignUp(users.User) error
	SignUpViaGoogle(users.User) (*users.User, error)
	GenerateState(http.ResponseWriter) string
	GetUserDataFromGoogle(string) (users.User, error)
	//ValidateToken(http.HandlerFunc, *server.Server) http.HandlerFunc
}

type service struct {
	userRepository users.Repository
	logger         *log.Logger
}

/*
SecuredCookie instance for encoding and decoding cookies
*/
var SecuredCookie *securecookie.SecureCookie

/*
NewAuthService creates a auth service with the necessary dependencies
*/
func NewAuthService(r users.Repository) Service {
	return &service{r, log.New(os.Stdout, "auth_service", log.LstdFlags|log.Lshortfile)}
}

/*
GenerateState genearate a state for verifying request to prevent CSRF
*/
func (s *service) GenerateState(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)
	SecuredCookie = securecookie.New([]byte(os.Getenv("SESSION_KEY")), []byte(os.Getenv("STATE_HASH_KEY")))

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	encodedState, err := SecuredCookie.Encode("oauth-state", state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	cookie := http.Cookie{
		Name:    "oauth-state",
		Expires: expiration,
		Value:   encodedState,
		/* Path:     "/",
		Secure:   true,
		HttpOnly: true, */
	}

	http.SetCookie(w, &cookie)
	return state
}

/*
GetUserDataFromGoogle feches user information from Google api
*/
func (s *service) GetUserDataFromGoogle(code string) (users.User, error) {
	const oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	// Use code to get token and get user info from Google.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return users.User{}, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		return users.User{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	var userInfo googleUserInfo
	err = json.NewDecoder(response.Body).Decode(&userInfo) //ioutil.ReadAll(response.Body)
	if err != nil {
		return users.User{}, fmt.Errorf("failed read response: %s", err.Error())
	}
	user := googleUserInfoToUserStruct(userInfo)

	return user, nil
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

func googleUserInfoToUserStruct(userInfo googleUserInfo) users.User {
	var user users.User
	user.ID = uuid.NewV4()
	user.Email = userInfo.Email
	user.Firstname = userInfo.GivenName
	user.Lastname = userInfo.FamilyName
	user.Username = strings.Split(userInfo.Email, "@")[0]
	user.Password = userInfo.ID
	user.Gender = "_" //Not provided by google api

	return user
}

/*
ValidateToken checks if user token is valid and authorises user to access route
*/
func ValidateToken(next http.HandlerFunc, srv *server.Server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next(w, r)
	})
}

/*
Login authenticates users
*/
func (s *service) Login(username string, password string) (*users.User, error) {
	user, err := s.userRepository.FindUserByUsername(username)
	if err != nil {
		return &users.User{}, err
	}
	passwordMatched := s.CheckPasswordHash(password, user.Password)
	if !passwordMatched {
		return &users.User{}, errors.New("password does not match")
	}

	return user, nil
}

/*
SignUp creates a new user
*/
func (s *service) SignUp(user users.User) error {
	var err error
	user.Password, err = s.HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = s.userRepository.AddUser(&user)
	if err != nil {
		return err
	}
	return nil

}

/*
SignUpViaGoogle creates a new user
*/
func (s *service) SignUpViaGoogle(user users.User) (*users.User, error) {
	var err error
	user.Password, err = s.HashPassword(user.Password)
	if err != nil {
		return &users.User{}, err
	}

	registeredUser, err := s.userRepository.FindOrAddUser(&user)
	if err != nil {

		return &users.User{}, err
	}
	return registeredUser, nil

}

/*
HashPassword encrypts user password
*/
func (s *service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/*
CheckPasswordHash checks if password entered matches the hashed password
*/
func (s *service) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
