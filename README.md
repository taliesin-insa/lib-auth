# auth_lib

This is the authentication library used by the microservices of the Taliesin project.
It is a simple wrapper around the auth microservice API, it delegates all the logic to it.

## How to use

1. Add the library to the imports of your go files `import "github.com/taliesin-insa/lib-auth"`
2. Make sure that the environment variable `AUTH_API_URL` contains the URL to the auth microservice
3. Call the method  `AuthenticateUser(r *http.Request)` with a request as a parameter to verify that the request was issued by a legitimate User and get its identity

## API

`AuthenticateUser(r *http.Request) (*UserData, error, int)`

Parameters:

1. r : the http.Request you want to authenticate, it shall contain the `Authorization` header with the user's token (enforced by the frontend)

Returned data:

1. a struct containing information about the authenticated user (see `UserData` definition), `nil` if an error occurred or the user wasn't authenticated
2. an error struct, `nil` if there was no error
3. an int corresponding to the HTTP error code associated to the error, value is http.StatusOK if there was no error 

## Example

```go
package test
import (
	"github.com/taliesin-insa/lib-auth"
	"net/http"
)

func createDatabase(w http.ResponseWriter, r *http.Request) {

	user, authErr, authStatusCode := lib_auth.AuthenticateUser(r)

    // first check if there was an error during the authentication or if the user wasn't authenticated
	if authErr != nil {
		w.WriteHeader(authStatusCode)
		w.Write([]byte("[AUTH] "+authErr.Error()))
		return
	}

    // check if the authenticated user has sufficient permissions to 
	if user.Role != lib_auth.RoleAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("[AUTH] Insufficient permissions to create database"))
	}

    // do some stuff ...
}
```
