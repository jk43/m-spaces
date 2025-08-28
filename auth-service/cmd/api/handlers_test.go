package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/moly-space/molylibs/network"
// 	"github.com/moly-space/molylibs/pb"
// )

// // func Test_SaveCredentials(t *testing.T) {
// // 	validateTests := []struct {
// // 		name         string
// // 		body         string
// // 		expectedCode int
// // 	}{
// // 		{"Incorrect JSON", `{"test`, http.StatusNotAcceptable},
// // 		{"Missing password", `{"password":"", "objectId":"123", "email":"test@test.com", "firstName":"alex", "lastName":"kim"}`, http.StatusBadRequest},
// // 		{"Missing objectId", `{"password":"password", "objectId":"", "email":"test@test.com", "firstName":"alex", "lastName":"kim"}`, http.StatusBadRequest},
// // 		{"Missing email", `{"password":"password", "objectId":"123", "email":"", "firstName":"alex", "lastName":"kim"}`, http.StatusBadRequest},
// // 		{"Incorrect email", `{"password":"password", "objectId":"123", "email":"bademail", "firstName":"alex", "lastName":"kim"}`, http.StatusBadRequest},
// // 		{"Missing first name", `{"password":"password", "objectId":"123", "email":"bademail", "firstName":"", "lastName":"kim"}`, http.StatusBadRequest},
// // 		{"Missing last name", `{"password":"password", "objectId":"123", "email":"bademail", "firstName":"alex", "lastName":""}`, http.StatusBadRequest},
// // 		{"InsertCredentials Error", `{"password":"password", "objectId":"123", "email":"error@test.com", "firstName":"alex", "lastName":"kim"}`, http.StatusExpectationFailed},
// // 		{"OK", `{"password":"password", "objectId":"123", "email":"test@test.com", "firstName":"alex", "lastName":"kim"}`, http.StatusAccepted},
// // 	}
// // 	for _, test := range validateTests {
// // 		var reader io.Reader
// // 		reader = strings.NewReader(test.body)
// // 		req := httptest.NewRequest("POST", "/", reader)
// // 		rr := httptest.NewRecorder()
// // 		handler := http.HandlerFunc(app.SaveCredentials)

// // 		handler.ServeHTTP(rr, req)

// // 		if rr.Code != test.expectedCode {
// // 			t.Errorf("%s expected code %d but got the %d", test.name, test.expectedCode, rr.Code)
// // 		}
// // 	}
// // }

// func TestLogin(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		body         string
// 		expectedCode int
// 	}{
// 		{"json error", "{", http.StatusNotAcceptable},
// 		{"mysql error", `{"password":"password1", "email":"error@test.com"}`, http.StatusExpectationFailed},
// 		{"incorrect password", `{"password":"password1", "email":"test@test.com"}`, http.StatusUnauthorized},
// 		{"ok", `{"password":"password", "email":"test@test.com"}`, http.StatusOK},
// 	}

// 	for _, test := range tests {
// 		rr, _ := utils.SimpleHttpTest("POST", app.Login, test.body)
// 		if rr.Code != test.expectedCode {
// 			t.Errorf("%s expected code %d but got the %d", test.name, test.expectedCode, rr.Code)
// 		}
// 		if rr.Code == http.StatusOK {
// 			var res map[string]any
// 			json.Unmarshal(rr.Body.Bytes(), &res)
// 			accessKey, ok := res["accessToken"]
// 			if !ok || accessKey == "" {
// 				t.Errorf("Unexpected accessKey")
// 			}
// 			refreshToken, ok := res["refreshToken"]
// 			if !ok || refreshToken == "" {
// 				t.Errorf("Unexpected refreshToken")
// 			}
// 		}
// 	}
// }

// func TestRefreshToken(t *testing.T) {

// 	//get the tokens from login
// 	body := `{"password":"password", "email":"test@test.com"}`
// 	rr, _ := utils.SimpleHttpTest("POST", app.Login, body)
// 	if rr.Code != http.StatusOK {
// 		t.Errorf("Expected code %d but got the %d", http.StatusOK, rr.Code)
// 	}
// 	var res map[string]any
// 	json.Unmarshal(rr.Body.Bytes(), &res)
// 	refreshToken, _ := res["refreshToken"].(string)
// 	// fmt.Println("rt", refreshToken)
// 	// return
// 	tests := []struct {
// 		name         string
// 		refreshToken string
// 		cookieName   string
// 		expectedCode int
// 	}{
// 		{"no refresh token", "", "", http.StatusUnauthorized},
// 		{"expired token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI4NTIwMDAsInN1YiI6IjYzYjRiMjI1MTMyZjlmYWQxYTMyNWYwZiJ9.mTQ-NtzdMp6eNmFEjxVe5oHc2XpNbfwSbeSG9AlnaCU", "_host_refresh_token", http.StatusUnauthorized},
// 		{"invalid token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI4NTIwMDAsInN1YiI6IjYzYjRiMjI1MTMyZjlmYWQxYTMyNWYwZiJ9.qZBawB9q0-8_wMpTOauCkICnEz3xLQLz2QMUHc5tK5c", "_host_refresh_token", http.StatusUnauthorized},
// 		{"OK", refreshToken, "_host_refresh_token", http.StatusOK},
// 	}

// 	for _, test := range tests {
// 		var reader io.Reader
// 		reader = strings.NewReader("")
// 		req, _ := http.NewRequest("POST", "/", reader)

// 		cookie := http.Cookie{
// 			Name:     test.cookieName,
// 			Value:    test.refreshToken,
// 			Path:     "/",
// 			HttpOnly: true,
// 			Secure:   false,
// 			SameSite: http.SameSiteStrictMode,
// 		}
// 		req.AddCookie(&cookie)

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(app.RefreshToken)

// 		handler.ServeHTTP(rr, req)
// 		if rr.Code != test.expectedCode {
// 			t.Errorf("%s expected code %d but got the %d", test.name, test.expectedCode, rr.Code)
// 		}
// 	}
// }

// func Test_GRPCSaveCredentials(t *testing.T) {
// 	validateTests := []struct {
// 		name           string
// 		in             pb.CredentialsRequest
// 		expectingError bool
// 	}{
// 		{"Missing all", pb.CredentialsRequest{Email: "", HashedPassword: "", UserID: "", FirstName: "", LastName: ""}, true},
// 		{"Missing Email", pb.CredentialsRequest{Email: "", HashedPassword: "password", UserID: "12345", FirstName: "alex", LastName: "kim"}, true},
// 		{"Missing ObjectId", pb.CredentialsRequest{Email: "test@test.com", HashedPassword: "password", UserID: "", FirstName: "alex", LastName: "kim"}, true},
// 		{"Missing HashedPassword", pb.CredentialsRequest{Email: "test@test.com", HashedPassword: "", UserID: "12345", FirstName: "alex", LastName: "kim"}, true},
// 		{"Missing First Name", pb.CredentialsRequest{Email: "test@test.com", HashedPassword: "", UserID: "12345", FirstName: "alex", LastName: "kim"}, true},
// 		{"Missing Last Name", pb.CredentialsRequest{Email: "test@test.com", HashedPassword: "", UserID: "12345", FirstName: "", LastName: "kim"}, true},
// 		{"Error InsertCredentials", pb.CredentialsRequest{Email: "error@test.com", HashedPassword: "password", UserID: "12345", FirstName: "alex", LastName: ""}, true},
// 		{"No Error", pb.CredentialsRequest{Email: "test@test.com", HashedPassword: "password", UserID: "12345", FirstName: "alex", LastName: "kim"}, false},
// 	}
// 	fmt.Println(validateTests[0].in.Email)
// 	server.SaveCredentials(context.Background(), &validateTests[0].in)

// 	for _, test := range validateTests {
// 		out, err := server.SaveCredentials(context.Background(), &test.in)
// 		if test.expectingError && err == nil {
// 			t.Errorf("%s is expecting error but no error found", test.name)
// 		}
// 		if !test.expectingError && err != nil {
// 			t.Errorf("%s is expecting no error but got an error %v", test.name, err)
// 		}
// 		if !test.expectingError && out.Id != 1 {
// 			t.Errorf("%s out.Id should be 1 but got an %v", test.name, out.Id)
// 		}
// 	}
// }
