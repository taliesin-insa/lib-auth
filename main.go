package lib_auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	RoleAdmin = iota
	RoleAnnotator = iota
)

type UserData struct {
	Username	string
	Role		int
}

type VerifyRequest struct {
	Token	string
}

var AuthAPI string

func getAPIUrl() {
	authEnvVal, authEnvExists := os.LookupEnv("AUTH_API_URL")

	if authEnvExists {
		AuthAPI = authEnvVal
	} else {
		AuthAPI = "http://auth-api.gitlab-managed-apps.svc.cluster.local:8080"
	}
}

func AuthenticateUser(r *http.Request) (*UserData, error, int) {
	getAPIUrl()

	authRequest, _ :=json.Marshal(VerifyRequest{Token: r.Header.Get("Authorization")})
	authRes, authReqErr := http.Post(AuthAPI+"/auth/verifyToken", "application/json", bytes.NewBuffer(authRequest))

	if authReqErr != nil {
		log.Printf("[LIB AUTH] Error in request to auth/verifyToken: %v", authReqErr.Error())
		return nil, errors.New("error in request to auth/verifyToken"), http.StatusInternalServerError
	}

	if authRes.StatusCode != http.StatusOK {
		log.Printf("[LIB AUTH] Error response from auth: %v", authRes.Body)
		return nil, errors.New("error response from auth"), http.StatusBadRequest
	}

	authResBody, _ := ioutil.ReadAll(authRes.Body)

	user := new(UserData)
	jsonErr := json.Unmarshal(authResBody, &user)

	if jsonErr != nil {
		log.Printf("[LIB AUTH] Error parsing auth/verifyToken: %v, error was %v", authResBody, jsonErr.Error())
		return nil, errors.New("error parsing auth/verifyToken"), http.StatusInternalServerError
	}

	return user, nil, http.StatusOK

}
