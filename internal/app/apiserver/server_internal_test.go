package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_handleUserCreate(t *testing.T) {
	s := newServer(teststore.New(), sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"name":     "User_1",
				"email":    "user1@spe.com",
				"password": "user_passport",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email":    "user1@sp",
				"password": "urt",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid payload",
			payload:      "some input",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			buf := &bytes.Buffer{}
			json.NewEncoder(buf).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/user/create", buf)
			s.ServeHTTP(rec, req)
			result := rec.Result()
			assert.Equal(t, tc.expectedCode, result.StatusCode)
		})
	}
}

func TestServer_handleUserLogin(t *testing.T) {
	store := teststore.New()
	user := model.TestUser(t)
	store.User().Create(user)

	srv := newServer(store, sessions.NewCookieStore([]byte("secret")))
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    user.Email,
				"password": user.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "some text",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "email",
				"password": user.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    user.Email,
				"password": "somepassword",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			buf := &bytes.Buffer{}
			json.NewEncoder(buf).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/user/login", buf)
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Result().StatusCode)
		})
	}
}

func TestServer_authUserMW(t *testing.T) {
	store := teststore.New()
	u := model.TestUser(t)
	store.User().Create(u)

	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("secret")
	s := newServer(store, sessions.NewCookieStore(secretKey))
	sc := securecookie.New(secretKey, nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.authUserMW(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_handleUserLogout(t *testing.T) {
	store := teststore.New()
	user := model.TestUser(t)
	store.User().Create(user)

	testCases := []struct {
		name        string
		cookieValue map[interface{}]interface{}
	}{
		{
			name: "logged out",
			cookieValue: map[interface{}]interface{}{
				"user_id": user.ID,
			},
		},
	}

	secretKey := []byte("secret")
	s := newServer(store, sessions.NewCookieStore(secretKey))
	sc := securecookie.New(secretKey, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.handleUserLogout().ServeHTTP(rec, req)
			assert.Empty(t, rec.Result().Header.Get("Cookie"))
		})
	}
}

func TestServer_handleWhoAmI(t *testing.T) {
	store := teststore.New()
	user := model.TestUser(t)
	store.User().Create(user)

	testCases := []struct {
		name         string
		user_id      interface{}
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name:         "authorized",
			user_id:      user.ID,
			expectedCode: http.StatusOK,
		},
		{
			name:         "unauthorized",
			user_id:      45,
			expectedCode: http.StatusUnauthorized,
		},
	}

	srv := newServer(store, nil)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/user/whoami", nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, tc.user_id))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Result().StatusCode)
		})
	}
}
