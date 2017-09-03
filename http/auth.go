package timesheetHttp

import (
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"golang.org/x/crypto/pbkdf2"
)

var (
	defaultTransport = &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ResponseHeaderTimeout: time.Millisecond * 200,
		IdleConnTimeout:       time.Millisecond * 800,
	}
	defaultClient = &http.Client{Transport: defaultTransport}
)

func logAndWrite(err error, logMsg string, w http.ResponseWriter) {
	log.Println(logMsg+",", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

var defaultSalt = []byte("timesheet app salt")

func genPassword(p string, salt []byte) string {
	if salt == nil {
		salt = defaultSalt
	}
	dk := pbkdf2.Key([]byte(p), salt, 4096, 32, sha512.New)
	return hex.EncodeToString(dk)
}

func basicValidation(name, val string, minLen, maxLen int) []string {
	errs := []string{}
	if len(val) < minLen {
		errs = append(errs, fmt.Sprintf("%q needs to be at least %d characters long", name, minLen))
	}
	if len(val) > maxLen {
		errs = append(errs, fmt.Sprintf("%q needs to be less than %d characters long", name, maxLen))
	}
	return errs
}

func writeValidationErrors(w http.ResponseWriter, vData map[string][]string) error {
	data, err := json.Marshal(vData)
	if err != nil {
		log.Printf("data: %v, error: %v\n", vData, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(data)
	return nil
}

func regHandler(
	sdbDBName, sdbInstanceAddr, sdbAPIKey, sdbAPIValue, address string,
) func(w http.ResponseWriter, r *http.Request) {
	userDataFormatString := fmt.Sprintf(
		"%s/db/%s/user/username/%s.json", sdbInstanceAddr, sdbDBName, "%s",
	)
	createUserFormatString := fmt.Sprintf(
		"%s/db/%s/user.json", sdbInstanceAddr, sdbDBName,
	)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				logAndWrite(err, "", w)
				return
			}
			defer r.Body.Close()

			un := r.FormValue("username")
			unErrors := basicValidation("username", un, 3, 35)

			req, _ := http.NewRequest("GET", fmt.Sprintf(userDataFormatString, un), nil)
			req.Header.Set(sdbAPIKey, sdbAPIValue)
			resp, err := defaultClient.Do(req)
			if err != nil {
				logAndWrite(err, fmt.Sprintf("couldn't find user %q, or service %q unavailable", un, sdbInstanceAddr), w)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logAndWrite(err, "error reading user response body", w)
				return
			}
			if len(body) > 2 {
				unErrors = append(unErrors, fmt.Sprintf("user %q exists, please select a diffrent user name", un))
			}

			email := r.FormValue("email")
			emailErrors := basicValidation("email", email, 5, 45)

			passwd := r.FormValue("password")
			passwdErrors := basicValidation("password", passwd, 6, 45)

			validationData := map[string][]string{}
			if len(unErrors) > 0 {
				validationData["username"] = unErrors
			}

			if len(emailErrors) > 0 {
				validationData["email"] = emailErrors
			}

			if len(passwdErrors) > 0 {
				validationData["password"] = passwdErrors
			}

			if len(validationData) > 0 {
				writeValidationErrors(w, validationData)
				return
			}

			encodedPass := genPassword(un+passwd, nil)
			payload := map[string]string{
				"username": un,
				"passwd":   encodedPass,
				"email":    email,
			}
			data, err := json.Marshal(payload)
			if err != nil {
				log.Printf("data: %v, error: %v", payload, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			req, _ = http.NewRequest("POST", createUserFormatString, bytes.NewReader(data))
			req.Header.Set(sdbAPIKey, sdbAPIValue)
			resp, err = defaultClient.Do(req)
			if err != nil {
				logAndWrite(err, fmt.Sprintf("couldn't create user %q, or service %q unavailable", un, address), w)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					logAndWrite(err, "error reading user creation response body", w)
					return
				}

				w.WriteHeader(resp.StatusCode)
				w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
				w.Write(body)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(fmt.Sprintf("User %q was created successfully!", un)))
		}
	}
}

var defaultSecret = []byte("timesheet app secret")

func genJWTToken(username string, id int, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": username,
		"id":       id,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	if len(secret) == 0 {
		secret = defaultSecret
	}
	return token.SignedString(defaultSecret)
}

func loginHandler(
	sdbDBName, sdbInstanceAddr, sdbAPIKey, sdbAPIValue string,
) func(w http.ResponseWriter, r *http.Request) {
	userDataFormatString := fmt.Sprintf(
		"%s/db/%s/user/username/%s.json?%s=%s", sdbInstanceAddr, sdbDBName, "%s",
	)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				logAndWrite(err, "", w)
				return
			}
			defer r.Body.Close()

			un := r.FormValue("username")
			unErrors := basicValidation("username", un, 3, 35)

			var body []byte
			if len(unErrors) == 0 {
				req, _ := http.NewRequest("GET", fmt.Sprintf(userDataFormatString, un), nil)
				req.Header.Set(sdbAPIKey, sdbAPIValue)
				resp, err := defaultClient.Do(req)
				if err != nil || resp.StatusCode != http.StatusOK {
					logAndWrite(
						err, fmt.Sprintf("couldn't find user %q, or service %q unavailable", un, sdbInstanceAddr), w,
					)
					return
				}
				defer resp.Body.Close()

				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					logAndWrite(err, "error reading user response body", w)
					return
				}

				if len(body) == 2 {
					unErrors = append(unErrors, "no such user")
				}
			}

			passwd := r.FormValue("password")
			passErrors := basicValidation("password", passwd, 6, 45)

			validationData := map[string][]string{}
			if len(unErrors) > 0 {
				validationData["username"] = unErrors
			}

			if len(passErrors) > 0 {
				validationData["password"] = passErrors
			}

			if len(validationData) > 0 {
				writeValidationErrors(w, validationData)
				return
			}

			var data []map[string]interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				logAndWrite(err, fmt.Sprintf("error marshalling %v", body), w)
				return
			}

			dataUn := data[0]["username"].(string)
			dataPasswd := data[0]["passwd"].(string)
			encodedPass := genPassword(un+passwd, nil)
			if dataUn != un || dataPasswd != encodedPass {
				errMsg := "wrong username or password"
				log.Printf(
					"%s, expected: u: %q, p: %q, got: u: %q, p: %q\n", errMsg, dataUn, dataPasswd, un, encodedPass,
				)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("{\"form\": \"" + errMsg + "\"}"))
				return
			}

			st, err := genJWTToken(un, int(data[0]["id"].(float64)), nil)
			if err != nil {
				logAndWrite(err, "error generating JWT token", w)
				return
			}
			tc := struct {
				Token string `json:"accessToken"`
			}{st}
			td, err := json.Marshal(tc)
			if err != nil {
				log.Printf("data: %v, error: %v", tc, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Write(td)
		}
	}
}

func SetupAuthHandlers(sdbDBName, sdbInstanceAddr, sdbAPIKey, sdbAPIValue, address string) {
	http.HandleFunc("/app/reg/", regHandler(sdbDBName, sdbInstanceAddr, sdbAPIKey, sdbAPIValue, address))
	http.HandleFunc("/app/login/", loginHandler(sdbDBName, sdbInstanceAddr, sdbAPIKey, sdbAPIValue))
}

func authorizationMiddleware(
	sdbDBName string,
	fn func(http.ResponseWriter, *http.Request),
	secret []byte,
) func(w http.ResponseWriter, r *http.Request) {
	baseURL := "/db/" + sdbDBName + "/timesheet/user_id/"
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			// we simply check the token claims, but this is a good place
			// to parse the r.URL.Path or other request parameters
			// and determine if a given user can access requested data
			// i.e. check if user of ID = 8 can access /db/timesheet/timesheet/user_id/8/project.json etc.
			mc := token.Claims.(jwt.MapClaims)
			userID, ok := mc["id"]
			if !ok {
				return nil, fmt.Errorf("token lacks 'id' claim")
			}

			userPath := baseURL + strconv.FormatFloat(userID.(float64), 'f', 0, 64)
			userPathLen := len(userPath) + 1
			// if: userID = 10
			// and: r.URL.Path = "/db/timesheet/timesheet/user_id/10/project.json
			// userPath = "/db/timesheet/timesheet/user_id/10/"
			// then: check if r.URL.Path starts with userPath
			if len(r.URL.Path) < userPathLen || (r.URL.Path[:userPathLen] != userPath+"/" && r.URL.Path[:userPathLen] != userPath+".") {
				return nil, fmt.Errorf("restricted access to this resource")
			}

			if _, ok = mc["username"]; !ok {
				return nil, fmt.Errorf("token lacks 'username' claim")
			}

			if len(secret) == 0 {
				secret = defaultSecret
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, http.StatusText(http.StatusUnauthorized)+": "+err.Error(), http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}
