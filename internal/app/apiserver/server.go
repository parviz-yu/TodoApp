package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "todoapp"
	ctxKeyUser  ctxKey = iota
)

var (
	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

// newStore returns a new instance of server.
func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()
	s.logger.Infof("Listening...")
	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/user/create", s.handleUserCreate()).Methods("POST")
	s.router.HandleFunc("/user/login", s.handleUserLogin()).Methods("POST")

	auth := s.router.PathPrefix("").Subrouter()
	auth.Use(s.authUserMW)
	auth.HandleFunc("/user/logout", s.handleUserLogout()).Methods("POST")
	auth.HandleFunc("/user/whoami", s.handleWhoAmI()).Methods("GET")

	auth.HandleFunc("/task/add", s.handleTaskAdd()).Methods("POST")
	auth.HandleFunc("/task/delete", s.handleTaskDelete()).Methods("DELETE").Queries("id", "{id}")
	auth.HandleFunc("/task/done", s.handleTaskDone()).Methods("PATCH").Queries("id", "{id}")
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) authUserMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		userId, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, ErrNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, userId)))
	})
}

func (s *server) handleUserCreate() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Name:     req.Name,
			Email:    strings.ToLower(req.Email),
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleUserLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !user.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, ErrIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 3, // 86400 seconds = 24 hours
			HttpOnly: true,
		}

		session.Values["user_id"] = user.ID
		if err = s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]string{"info": "you've successfully logged in"})
	}
}

func (s *server) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Options.MaxAge = -1
		delete(session.Values, "user_id")
		if err := session.Save(r, w); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]string{"info": "you've successfully logged out"})
	}
}

func (s *server) handleWhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(ctxKeyUser).(int)
		user, err := s.store.User().FindById(userId)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, ErrNotAuthenticated)
			return
		}

		s.respond(w, r, http.StatusOK, user)
	}
}

func (s *server) handleTaskAdd() http.HandlerFunc {
	type Request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		task := &model.Task{
			UserID:       r.Context().Value(ctxKeyUser).(int),
			Title:        req.Title,
			Description:  req.Description,
			Done:         false,
			CreationDate: time.Now().Format("02/01/06"),
		}

		if err := s.store.Task().Create(task); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, task)
	}
}

func (s *server) handleTaskDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(ctxKeyUser).(int)
		taskId, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, store.ErrInvalidTaskId)
			return
		}

		if err := s.store.Task().Delete(userId, taskId); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]string{
			"info": "you've successfully deleted a task",
		})
	}
}

func (s *server) handleTaskDone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(ctxKeyUser).(int)
		taskId, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Task().Done(userId, taskId); err != nil {
			if err == store.ErrInvalidTaskId {
				s.error(w, r, http.StatusNotFound, err)
				return
			}

			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]string{
			"info": "congrats! you've done a task",
		})
	}
}

// error wrapper for respond function
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// helper function for writing a json respond
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
