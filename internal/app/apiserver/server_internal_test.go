package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_handleUserCreate(t *testing.T) {
	s := newServer(teststore.New(), nil)

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
			req, _ := http.NewRequest(http.MethodPost, "/sign-up", buf)
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
			req, _ := http.NewRequest(http.MethodPost, "/sign-in", buf)
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

	cookieStore, secureCookie := TestSession(t)
	s := newServer(store, cookieStore)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/users", nil)
			cookieStr, _ := secureCookie.Encode(sessionName, tc.cookieValue)
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

	cookieStore, secureCookie := TestSession(t)
	s := newServer(store, cookieStore)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := secureCookie.Encode(sessionName, tc.cookieValue)
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
			req, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, tc.user_id))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Result().StatusCode)
		})
	}
}

func TestServer_handleTaskCreate(t *testing.T) {
	store := teststore.New()
	task := model.TestTask(t)
	srv := newServer(store, nil)

	testCases := []struct {
		name         string
		user_id      interface{}
		payload      map[string]string
		expectedCode int
	}{
		{
			name:    "valid",
			user_id: task.UserID,
			payload: map[string]string{
				"title":       task.Title,
				"description": task.Description,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:    "invalid title",
			user_id: task.UserID,
			payload: map[string]string{
				"title":       "",
				"description": task.Description,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:    "invalid description",
			user_id: task.UserID,
			payload: map[string]string{
				"title":       task.Title,
				"description": "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			buf := &bytes.Buffer{}
			json.NewEncoder(buf).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/users/tasks", buf)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, tc.user_id))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Result().StatusCode)
		})
	}
}

func TestServer_handleTaskDelete(t *testing.T) {
	store := teststore.New()
	task := model.TestTask(t)
	store.Task().Create(task)
	srv := newServer(store, nil)

	testCases := []struct {
		name         string
		user_id      interface{}
		queryString  string
		expectedCode int
	}{
		{
			name:         "valid id",
			user_id:      task.UserID,
			queryString:  "/users/tasks/1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid id",
			user_id:      task.UserID,
			queryString:  "/users/tasks/asdas",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "not existing id",
			user_id:      task.UserID,
			queryString:  "/users/tasks/564",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, tc.queryString, nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, tc.user_id))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Result().StatusCode)
		})
	}
}

func TestServer_handleTaskDone(t *testing.T) {
	store := teststore.New()
	task := model.TestTask(t)
	store.Task().Create(task)
	srv := newServer(store, nil)

	testCases := []struct {
		name         string
		user_id      interface{}
		queryString  string
		expectedCode int
	}{
		{
			name:         "valid id",
			user_id:      task.UserID,
			queryString:  "/users/tasks/1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid id",
			user_id:      task.UserID,
			queryString:  "/users/tasks/asdas",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "not existing id",
			user_id:      task.UserID,
			queryString:  "/users/tasks/564",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, tc.queryString, nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, tc.user_id))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Result().StatusCode)
		})
	}
}

func TestServer_handleTaskGet(t *testing.T) {
	store := teststore.New()
	task := model.TestTask(t)
	store.Task().Create(task)
	srv := newServer(store, nil)

	testCases := []struct {
		name         string
		userId       interface{}
		queryString  string
		expectedCode int
	}{
		{
			name:         "valid id",
			userId:       task.UserID,
			queryString:  "/users/tasks/1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid id",
			userId:       task.UserID,
			queryString:  "/users/tasks/id",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "not existing id",
			userId:       task.UserID,
			queryString:  "/users/tasks/150",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.queryString, nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, tc.userId))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Result().StatusCode)
		})
	}
}

func TestServer_handleTaskGetDone(t *testing.T) {
	store := teststore.New()
	task := model.TestTask(t)
	store.Task().Create(task)
	srv := newServer(store, nil)

	testCases := []struct {
		name         string
		userId       interface{}
		queryString  string
		expectedCode int
	}{
		{
			name:         "valid argument",
			userId:       task.UserID,
			queryString:  "/users/tasks?done=false",
			expectedCode: http.StatusOK,
		},
		{
			name:         "not found task",
			userId:       task.UserID,
			queryString:  "/users/tasks?done=true",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid argument",
			userId:       task.UserID,
			queryString:  "/users/tasks?done=2",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.queryString, nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, tc.userId))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Result().StatusCode)
		})
	}
}

func TestServer_handleTaskGetAll(t *testing.T) {
	store := teststore.New()
	task := model.TestTask(t)
	srv := newServer(store, nil)

	testCases := []struct {
		name         string
		userId       interface{}
		create       bool
		expectedCode int
	}{
		{
			name:         "no tasks in storage",
			userId:       task.UserID,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "valid",
			userId:       task.UserID,
			create:       true,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			if tc.create {
				store.Task().Create(task)
			}
			req, _ := http.NewRequest(http.MethodGet, "/users/tasks", nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, tc.userId))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Result().StatusCode)
		})
	}
}
