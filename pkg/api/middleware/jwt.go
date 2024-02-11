package middleware

import (
	"carbide-images-api/pkg/api/utils"
	"carbide-images-api/pkg/objects"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	// "github.com/rs/zerolog/log"
)

func checkAuth(w http.ResponseWriter, r *http.Request) {
	if !authenticated(w, r) {
		log.Info("user is unauthorized\n")
		utils.RespondWithJSON(w, "user is unauthorized")
		return
	}
}

func authenticated(w http.ResponseWriter, r *http.Request) bool {
	_, err := verifyJWT(r)
	if err == nil {
		return true
	}
	return false
}

func setAuthCookie(w http.ResponseWriter, user objects.User) error {
	token, err := generateJWT(user)
	if err != nil {
		return err
	}
	ck := http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &ck)
	return nil
}

// generate JWT from given user - returns err and token
func generateJWT(user objects.User) (string, error) {
	secret := os.Getenv("JWTSECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	claims["userid"] = fmt.Sprint(user.Id)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// checks if http request is authorized/logged in - returns error and username string; empty if err
func verifyJWT(r *http.Request) (int64, error) {
	secret := os.Getenv("JWTSECRET")
	// get token from cookie
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			return 0, err
		}
		// For any other type of error, return a bad request status
		return 0, err
	}
	// Get the JWT string from the cookie
	tokenString := c.Value
	// parse and check token validity
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid JWT")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("failed to parse JWT claims")
	}
	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return 0, errors.New("token expired")
	}
	userIdString := claims["userid"].(string)
	userid, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return 0, errors.New("failed to parse userid")
	}
	return userid, nil
}

func terminateJWT() {
	// replace jwt with another that expires immediately
}