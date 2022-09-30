package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"homeHandler/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func hashingSecret(secret string) string {
	data := "data"

	if len(secret) == 0 {
		log.Fatal("no secret was provided to sign the jwt")
	}

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	return sha

}

func createJwt(email string) (string, error) {
	var claims jwt.Claims
	claims.Subject = email
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	// Those fields had to be filled with domain name
	claims.Issuer = "localhost:4000"
	claims.Audiences = []string{"localhost:4000"}

	// secret used to sign the token
	secret := hashingSecret(os.Getenv("jwt_secret"))

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(secret))

	if err != nil {
		return "", err
	}

	return string(jwtBytes), err
}

func (app *application) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.User

	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	userExists := app.models.DB.CheckExistingUser(newUser.Email)

	if userExists {
		err = errors.New("the user you are trying to register already exists")
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	userPwd := newUser.Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userPwd), 12)
	newUser.Password = string(hashedPassword)

	err = app.models.DB.CreateUser(newUser)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var cred Credentials

	cred.Username = newUser.Email
	cred.Password = newUser.Password

	token, err := createJwt(newUser.Email)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	type jwtResponse struct {
		Token string `json:"token"`
	}

	var response jwtResponse

	response.Token = token

	resJson, err := json.MarshalIndent(&response, "", " ")

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}

func (app *application) SignIn(w http.ResponseWriter, r *http.Request) {
	var cred Credentials

	err := json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	hashedPwd, err := app.models.DB.GetUserPassword(cred.Username)

	if err != nil || len(hashedPwd) == 0 {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(cred.Password))

	if err != nil || len(hashedPwd) == 0 {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	token, err := createJwt(cred.Username)

	if err != nil || len(hashedPwd) == 0 {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var response struct {
		Token string `json:"token"`
	}

	response.Token = token

	res, err := json.MarshalIndent(&response, "", " ")

	if err != nil || len(hashedPwd) == 0 {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
