package main

// https://sendgrid.com/blog/tokens-tokens-intro-json-web-tokens-jwt-go/

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// TODO, right key size?
	privateKey *rsa.PrivateKey
	publicKey  crypto.PublicKey
)

func Init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("error generating RSA key: %s", err)
	}
	publicKey = privateKey.Public()
}

func GenAuth(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.GetSigningMethod("RS256")) // Create a Token that will be signed with RSA 256.
	// { "typ":"JWT", "alg":"RS256" }

	token.Claims["ID"] = "This is my super fake ID"
	token.Claims["exp"] = time.Now().Unix() + 36000

	// The claims object allows you to store information in the actual token.
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Printf("error generating token: %s", err)
	}

	// tokenString Contains the actual token you should share with your client.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"token\": %s}", tokenString)
}

func CheckAuth(r *http.Request) (uid int, err error) {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if token.Valid {

		//YAY!
	} else {
		//Someone is being funny
	}
	return 0, nil
}
