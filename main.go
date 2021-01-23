package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/gofrs/uuid"
)

type UserClaims struct {
	jwt.StandardClaims
	SessionID int64
}

func(u *UserClaims) Valid() error {
	if !u.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("Token has expired")
	}

	if u.SessionID == 0 {
		return fmt.Errorf("Invalid session ID")
	}

	return nil
}

type person struct {
	First string
}

func main() {
	pass := "123456789"

	hashedPass, err := hashPassword(pass)
	if err != nil {
		panic(err)
	}

	err = comparePassword(pass, hashedPass)
	if err != nil {
		log.Fatalln("Not logged in")
	}

	log.Println("Logged in!")
}


func hashPassword(password string) ([]byte, error) {
	bs, err :=  bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error while generationg bcrypt from password: %w", err)
	}
	return bs, nil
}

func comparePassword(password string, hashedPass []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	if err != nil {
		return fmt.Errorf("Invalid password: %w", err)
	}
	return nil
}

func signMessage(msg []byte) ([]byte, error) {
	//second argument is private key
	h := hmac.New(sha512.New, keys[currentKid].key)
	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("Error in signMessage hashing message: %w", err)
	}

	signature := h.Sum(nil)
	return signature, nil
}

func checkSig(msg, sig []byte) (bool, error){
	newSig, err := signMessage(msg)
	if err != nil {
		return false, fmt.Errorf("Error in check Sig while getting signature of message: %w", err)
	}

	same := hmac.Equal(newSig, sig)
	return same, nil
}
func createToken(c *UserClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	signedToken, err := t.SignedString(keys[currentKid])
	if err != nil {
		return "", fmt.Errorf("Error in create token when signing token: %w", err)
	}

	return signedToken, nil
}

func generateNewKey() error {
	newKey := make([]byte, 64)
	_, err :=io.ReadFull(rand.Reader, newKey)
	if err != nil {
		return fmt.Errorf("Error in generate a new key %w", err)
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("Error in generate a new kid %w", err)
	}

	keys[uid.String()] = key {
		key: newKey,
		created: time.Now(),
	}
	currentKid = uid.String()
	return nil
}

type key struct {
	key []byte
	created time.Time
}

var currentKid = ""
var keys = map[string]key{}

func parseToken(signedToken string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(signedToken, &UserClaims{}, func (t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() == jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("Invalid a signing algorithm")
		}

		//kid = key id, using for rotating keys for better security
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("Invalid key ID")
		}

		k, ok := keys[kid]
		if !ok {
			return nil, fmt.Errorf("Invalid key ID")
		}

		return k, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error in parse token: %w", err)
	}

	if !t.Valid {
		return nil, fmt.Errorf("Error in validationg token: %w", err)
	}

	return t.Claims.(*UserClaims), nil
}