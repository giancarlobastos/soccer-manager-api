package security

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/giancarlobastos/soccer-manager-api/domain"
	"github.com/giancarlobastos/soccer-manager-api/service"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

var jwtKey = []byte("you should get it from OS environment!")

type AuthenticationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticationResponse struct {
	Token string `json:"token"`
}

type AuthenticationMiddleware struct {
	accountService  *service.AccountService
	routeProfileMap map[string]string
}

func NewAuthenticationMiddleware(as *service.AccountService, m map[string]string) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		accountService:  as,
		routeProfileMap: m,
	}
}

func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		if route != nil {
			if role, ok := amw.routeProfileMap[route.GetName()]; ok && role == "USER" {
				if claims, err := amw.getClaims(w, r); err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
				} else {
					accountId, _ := strconv.Atoi(claims.Id)
					user := domain.User{
						AccountId: accountId,
						Username:  claims.Subject,
					}

					next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
				}
			} else {
				next.ServeHTTP(w, r)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (amw *AuthenticationMiddleware) Authenticate(w http.ResponseWriter, r *http.Request) {
	var authenticationRequest AuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&authenticationRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := amw.generateToken(authenticationRequest.Username, authenticationRequest.Password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	response, _ := json.Marshal(AuthenticationResponse{Token: *token})

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (amw *AuthenticationMiddleware) getClaims(w http.ResponseWriter, r *http.Request) (*jwt.StandardClaims, error) {
	bearerToken := r.Header.Get("Authorization")

	if len(bearerToken) < 8 {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, errors.New("bearer token not provided")
	}

	bearerToken = bearerToken[7:]
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, errors.New("unauthorized")
	}

	return claims, nil
}

func (amw *AuthenticationMiddleware) generateToken(username, password string) (*string, error) {
	account, err := amw.accountService.GetAccountByUsername(username)

	if err != nil {
		return nil, err
	} else if !account.Confirmed {
		return nil, errors.New("account not verified")
	} else if account.Locked {
		return nil, errors.New("account locked")
	}

	expiresAt := time.Now().Add(time.Hour).Unix()

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		_ = amw.accountService.RegisterFailedLoginAttempt(username)
		return nil, errors.New("password mismatch")
	}

	_ = amw.accountService.ResetLoginAttempts(username)
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"),
		&jwt.StandardClaims{
			Id:        strconv.Itoa(account.Id),
			Subject:   account.Username,
			ExpiresAt: expiresAt,
		})

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
