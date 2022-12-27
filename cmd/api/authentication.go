package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"homeHandler/models"
	"log"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"github.com/spf13/viper"
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

func createJwt(email string, expirationTime int, secretKey ...string) (string, error) {
	var claims jwt.Claims
	claims.Subject = email
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(time.Duration(expirationTime)))
	// Those fields had to be filled with domain name
	claims.Issuer = "localhost:4000"
	claims.Audiences = []string{"localhost:4000"}

	// secret used to sign the token
	var secret string

	if secretKey != nil {
		secret = secretKey[0]
	} else {
		secret = hashingSecret(viper.GetString("JWT_SECRET"))
	}

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

	duration := 24 * time.Hour
	token, err := createJwt(newUser.Email, int(duration))

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

	userID, hashedPwd, err := app.models.DB.GetUserStoredCredentials(cred.Username)

	if err != nil || len(hashedPwd) == 0 {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(cred.Password))

	if err != nil {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	duration := 24 * time.Hour
	token, err := createJwt(cred.Username, int(duration))

	if err != nil || len(hashedPwd) == 0 {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var response struct {
		ID    int    `json:"id"`
		Token string `json:"token"`
	}

	response.Token = token
	response.ID = userID

	res, err := json.MarshalIndent(&response, "", " ")

	if err != nil || len(hashedPwd) == 0 {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (app *application) resetPassword(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Search in database if user exists and eventually retrive it
	user, err := app.models.DB.GetUserByUsername(payload.Username)

	if err != nil {
		app.errorJSON(w, errors.New("this user doesn't exists"), http.StatusBadRequest)
		return

	}

	date := time.Now().Format("2006-01-02 15:04:05")

	// Saving the date when the token was require
	err = app.models.DB.SavingResetPasswordRequestDate(user.ID, date)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Generate a token that will be embed in the link in the email for password reset
	jwt, err := createJwt(user.Email, int(5*time.Minute))

	if err != nil {
		app.errorJSON(w, errors.New("error while saving the temporary auth token"))
		return
	}

	resetLink := fmt.Sprintf("http://localhost:3000/passwordReset/%s", jwt)
	msg := fmt.Sprintf(`Hello %s %s ,<br>
	this email was sent to you with the purpouse of resetting your password. <br> 
	The procedure is very simple, click on the link below and follow the instructions.<br>
	<a href="%s">Go to password reset page</a>`, user.Name, user.LastName, resetLink)

	ch := app.mailChan

	var emailInfo struct {
		To  string
		Msg string
	}

	emailInfo.Msg = msg
	emailInfo.To = user.Email

	// Send email text content through the email channel
	ch <- emailInfo

	app.writeJSON(w, http.StatusOK, "Email sended", "message")
}

func (app *application) SaveNewPasword(w http.ResponseWriter, r *http.Request) {
	var payload models.PasswordResetSchema

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if len([]rune(payload.NewPassword)) == 0 {
		err = errors.New("password cannot be empty")
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	claims, err := jwt.ParseWithoutCheck([]byte(payload.Token))

	if err != nil {
		app.errorJSON(w, errors.New("error while parsing token"), http.StatusUnauthorized)
		return
	}

	userEmail := claims.Subject
	user, err := app.models.DB.GetUserByUsername(userEmail)

	if err != nil || user.ResetRequestDate == "not-set" {
		app.errorJSON(w, errors.New("any reset password email was emitted before"), http.StatusUnauthorized)
		return
	}

	secret := viper.GetString("JWT_SECRET")
	parsedToken, err := jwt.HMACCheck([]byte(payload.Token), []byte(hashingSecret(secret)))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if !parsedToken.Valid(time.Now()) {
		app.errorJSON(w, errors.New("token has expired"))
		return
	}

	// Get the new password from the request payload
	newPwd, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), 12)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.SetNewPassword(user.ID, string(newPwd))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, "Password changed correctly", "message")
}
