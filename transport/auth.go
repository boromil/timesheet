package transport

import (
	"crypto/sha512"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"gitlab.com/boromil/goslashdb/slashdb"

	"golang.org/x/crypto/pbkdf2"
)

var (
	defaultTransport = &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ResponseHeaderTimeout: time.Second * 10,
		IdleConnTimeout:       time.Second * 10,
		MaxIdleConns:          30,
		MaxIdleConnsPerHost:   3,
	}
	defaultClient = &http.Client{Transport: defaultTransport}
)

func logAndWrite(err error, logMsg string, w http.ResponseWriter) {
	log.Println(logMsg+",", err.Error())
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
	sdbService slashdb.Service,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				logAndWrite(err, "", w)
				return
			}
			defer r.Body.Close()

			un := r.FormValue("username")
			unErrors := basicValidation("username", un, 3, 35)

			userReq := slashdb.NewDataRequest("")
			userReq.AddParts(
				slashdb.Part{Name: "timesheet"},
				slashdb.Part{
					Name: "user",
					Filter: slashdb.Filter{
						Values: map[string][]string{"username": []string{un}},
						Order:  []string{"username"},
					},
					Fields: []string{"username"},
				},
			)

			userNames := []string{}
			if err = sdbService.Get(r.Context(), userReq, &userNames); err != nil {
				logAndWrite(err, fmt.Sprintf("couldn't find user %q or SlashDB instance unavailable", un), w)
				return
			}

			if len(userNames) > 0 {
				unErrors = append(unErrors, fmt.Sprintf("user %q exists, please select a diffrent user name", un))
			}

			validationData := map[string][]string{}
			if len(unErrors) > 0 {
				validationData["username"] = unErrors
			}

			email := r.FormValue("email")
			emailErrors := basicValidation("email", email, 5, 45)
			if len(emailErrors) > 0 {
				validationData["email"] = emailErrors
			}

			passwd := r.FormValue("password")
			passwdErrors := basicValidation("password", passwd, 6, 45)
			if len(passwdErrors) > 0 {
				validationData["password"] = passwdErrors
			}

			if len(validationData) > 0 {
				writeValidationErrors(w, validationData)
				return
			}

			userData := User{
				Username: un,
				Passwd:   genPassword(un+passwd, nil),
				Email:    email,
			}
			userReq = slashdb.NewDataRequest("")
			userReq.AddParts(
				slashdb.Part{
					Name:   "timesheet",
					Fields: []string{"user"},
				},
			)
			if _, err = sdbService.Create(r.Context(), userReq, userData); err != nil {
				logAndWrite(err, fmt.Sprintf("couldn't create user %q SlashDB instance unavailable", un), w)
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
	sdbService slashdb.Service,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				logAndWrite(err, "", w)
				return
			}
			defer r.Body.Close()

			un := r.FormValue("username")
			unErrors := basicValidation("username", un, 3, 35)

			userReq := slashdb.NewDataRequest("")
			userReq.AddParts(
				slashdb.Part{Name: "timesheet"},
				slashdb.Part{
					Name: "user",
					Filter: slashdb.Filter{
						Values: map[string][]string{
							"username": []string{un},
						},
					},
				},
			)

			userData := []User{}
			if len(unErrors) == 0 {
				if err = sdbService.Get(r.Context(), userReq, &userData); err != nil {
					logAndWrite(err, fmt.Sprintf("couldn't find user %q or SlashDB instance unavailable", un), w)
					return
				}

				if len(userData) != 1 {
					unErrors = append(unErrors, "no such user")
				}
			}

			validationData := map[string][]string{}
			if len(unErrors) > 0 {
				validationData["username"] = unErrors
			}

			passwd := r.FormValue("password")
			passErrors := basicValidation("password", passwd, 6, 45)
			if len(passErrors) > 0 {
				validationData["password"] = passErrors
			}

			if len(validationData) > 0 {
				writeValidationErrors(w, validationData)
				return
			}

			dataUn := userData[0].Username
			dataPasswd := userData[0].Passwd
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

			st, err := genJWTToken(un, userData[0].ID, nil)
			if err != nil {
				logAndWrite(err, "error generating JWT token", w)
				return
			}
			w.Write([]byte(`{"accessToken":"` + st + `"}`))
		}
	}
}

// SetupAuthHandlers adds user auth endpoints
func SetupAuthHandlers(
	sdbService slashdb.Service,
) {
	http.HandleFunc("/app/reg/", regHandler(sdbService))
	http.HandleFunc("/app/login/", loginHandler(sdbService))
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
