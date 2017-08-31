package filemanager

import (
	"crypto/rand"
    "crypto/ecdsa"
	"encoding/json"
	"net/http"
	"strings"
	"time"
    "fmt"
    "os"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

var jwtPubKey *ecdsa.PublicKey


const (
	iyoPubKey = `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAES5X8XrfKdx9gYayFITc89wad4usrk0n2
7MjiGYvqalizeSWTHEpnd7oea9IQ8T5oJjMVH5cc0H5tFSKilFFeh//wngxIyny6
6+Vq5t5B0V0Ehy01+2ceEon2Y0XDkIKv
-----END PUBLIC KEY-----`
)

/*
// staging
const (
	iyoPubKey = `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEkmd07vxBqoCiHsaplIpjlonDeOnpvPam
ORMdBcAlHNXbzwplcdK4qlZGPBz9mxDSrBOv9SZH+Et6r8gn9Fx/+ZjlvRwowqOU
FpCIijAEx6A3BhfRUbmwl1evBKzWB/qw
-----END PUBLIC KEY-----`
)
*/

func init() {
	var err error

	jwtPubKey, err = jwt.ParseECPublicKeyFromPEM([]byte(iyoPubKey))
	if err != nil {
		fmt.Printf("failed to parse public key: %v\n", err)
		os.Exit(1)
	}
}

// authHandler proccesses the authentication for the user.
func authHandler(c *RequestContext, w http.ResponseWriter, r *http.Request) (int, error) {
	// NoAuth instances shouldn't call this method.
	if c.NoAuth {
		return 0, nil
	}

	// Receive the credentials from the request and unmarshal them.
	var cred User
	if r.Body == nil {
		return http.StatusForbidden, nil
	}

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		return http.StatusForbidden, nil
	}

	// Checks if the user exists.
	u, ok := c.Users[cred.Username]
	if !ok {
		return http.StatusForbidden, nil
	}

	// Checks if the password is correct.
	if !checkPasswordHash(cred.Password, u.Password) {
		return http.StatusForbidden, nil
	}

	c.User = u
	return printToken(c, w)
}

// renewAuthHandler is used when the front-end already has a JWT token
// and is checking if it is up to date. If so, updates its info.
func renewAuthHandler(c *RequestContext, w http.ResponseWriter, r *http.Request) (int, error) {
	ok, u := validateAuth(c, r)
	if !ok {
		return http.StatusForbidden, nil
	}

	c.User = u
	return printToken(c, w)
}

// claims is the JWT claims.
type claims struct {
	User
	NoAuth bool `json:"noAuth"`
	jwt.StandardClaims
}

// printToken prints the final JWT token to the user.
func printToken(c *RequestContext, w http.ResponseWriter) (int, error) {
	// Creates a copy of the user and removes it password
	// hash so it never arrives to the user.
	u := User{}
	u = *c.User
	u.Password = ""

	// Builds the claims.
	claims := claims{
		u,
		c.NoAuth,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "File Manager",
		},
	}

	// Creates the token and signs it.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(c.key)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Writes the token.
	w.Header().Set("Content-Type", "cty")
	w.Write([]byte(signed))
	return 0, nil
}

type extractor []string

func (e extractor) ExtractToken(r *http.Request) (string, error) {
	token, _ := request.AuthorizationHeaderExtractor.ExtractToken(r)

	// Checks if the token isn't empty and if it contains two dots.
	// The former prevents incompatibility with URLs that previously
	// used basic auth.
	if token != "" && strings.Count(token, ".") == 2 {
		return token, nil
	}

	cookie, err := r.Cookie("caddyauth")
	if err != nil {
		return "", request.ErrNoTokenInRequest
	}

    fmt.Println(cookie)

	return cookie.Value, nil
}

// validateAuth is used to validate the authentication and returns the
// User if it is valid.
func validateAuth(c *RequestContext, r *http.Request) (bool, *User) {
	if c.NoAuth {
		c.User = c.DefaultUser
		return true, c.User
	}

    var tokenStr = ""

    for _, v := range r.Cookies() {
        if v.Name == "caddyoauth" {
            tokenStr = v.Value
        }
    }

    // verify token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if token.Method != jwt.SigningMethodES384 {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtPubKey, nil
	})

	if err != nil {
        fmt.Printf("parse from request failed: %v\n", err)
		return false, nil
	}

	// get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
        fmt.Errorf("invalid token")
		return false, nil
	}

    u := c.Users["admin"]
    username, ok := claims["username"].(string)
    if !ok {
        fmt.Errorf("username not set, this should not happen")
		return false, nil
    }

    scopes, ok := claims["scope"].([]interface{})
    if !ok {
        fmt.Errorf("scopes not set, this should not happen")
		return false, nil
    }

    for _, value := range scopes {
        scope := value.(string)

        // this scope is not for this, out of bounds
        if len(scope) < 12 {
            continue
        }

        if scope[0:13] == "[user:email]:" {
            u.Email = scope[13:]
        }

        if scope[0:12] == "[user:name]:" {
            u.RealName = scope[12:]
        }
    }


    u.Username = username
    u.Admin = false

    fmt.Printf("User logged in as %s (%s, %s)\n", u.Username, u.RealName, u.Email)

	c.User = u
	return true, u
}

// hashPassword generates an hash from a password using bcrypt.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash compares a password with an hash to check if they match.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
