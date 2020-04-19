package lib_auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var validUserData = UserData{Username:"morpheus", Role:RoleAdmin}

func MockAuthMicroservice() *httptest.Server {

	mockedServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/auth/verifyToken" {

				body, _ := ioutil.ReadAll(r.Body)
				var request VerifyRequest
				json.Unmarshal(body, &request)

				if request.Token == "valid_token" {
					w.WriteHeader(http.StatusOK)
					r, _ := json.Marshal(validUserData)
					w.Write(r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
				}
			}
		}))

	return mockedServer
}


func TestMain(m *testing.M) {
	previousAuthUrl := os.Getenv("AUTH_API_URL")

	m.Run()

	os.Setenv("AUTH_API_URL", previousAuthUrl)
}


func TestAuthOk(t *testing.T) {

	/* Mocking Database API Response */
	mockedAuthServer := MockAuthMicroservice()

	os.Setenv("AUTH_API_URL", mockedAuthServer.URL)

	header := http.Header{}
	header.Add("Authorization", "valid_token")

	/* Mock http request, here in createDatabase we don't use the request struct so we can pass a blank one */
	request := &http.Request{
		Method: http.MethodPost,
		Header: header,
	}

	result, err, _ := AuthenticateUser(request)

	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
		return
	}

	if result.Username != validUserData.Username && result.Role != validUserData.Role {
		t.Errorf("unexpected result data, got %v want %v",
			result, validUserData)
		return
	}

	mockedAuthServer.Close()
}

func TestAuthUnauthorized(t *testing.T) {

	/* Mocking Database API Response */
	mockedAuthServer := MockAuthMicroservice()

	os.Setenv("AUTH_API_URL", mockedAuthServer.URL)

	header := http.Header{}
	header.Add("Authorization", "invalid_token")

	/* Mock http request, here in createDatabase we don't use the request struct so we can pass a blank one */
	request := &http.Request{
		Method: http.MethodPost,
		Header: header,
	}

	_, err, statusCode := AuthenticateUser(request)

	if err.Error() != "error response from auth" || statusCode != http.StatusBadRequest {
		t.Errorf("unexpected error: got '%v' want '%v'", err.Error(), "error response from auth")
		return
	}

	mockedAuthServer.Close()
}