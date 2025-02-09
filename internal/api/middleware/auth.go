package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rafael-cagliari/doggo-api/config"
)

// Cognito info
func init(){
	config.LoadEnv()
}

var userPooolID = config.GetEnv("COGNITO_USER_POOL_ID")
var region = config.GetEnv("AWS_REGION")

type JWK struct {
	Kid string   `json:"kid"`
	Kty string   `json:"kty"`
	Alg string   `json:"alg"`
	N   string   `json:"n"`
	E   string   `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWT auth middleware 
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		authHeader := r.Header.Get("Authorization")	
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Pegar a chave pública do Cognito (JWKs)
			return GetCognitoPublicKey(token)
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetCognitoPublicKey(token *jwt.Token) (interface{}, error) {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPooolID)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, errors.New("failed to fetch JWKs")
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, errors.New("erro ao decodificar JWKS")
	}

	// Finds matching public key on JWKS
	for _, key := range jwks.Keys {
		if key.Kid == token.Header["kid"]{
			return jwt.ParseRSAPublicKeyFromPEM([]byte(key.N))
		}
	}

	return nil, errors.New("public key not found")
}