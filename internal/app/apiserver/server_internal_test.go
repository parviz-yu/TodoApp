package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
