package apiserver

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/yeboka/final-project/internal/app/model"
	"github.com/yeboka/final-project/internal/app/store"
	"net/http"
)

const (
	sessionName = "canteen"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionsStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionsStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type requests struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		req := &requests{}
		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(writer, request, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type requests struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		req := &requests{}
		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(writer, request, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(request, sessionName)
		if err != nil {
			s.error(writer, request, http.StatusInternalServerError, err)
			return
		}

		session.Values["values_id"] = u.ID
		if err := s.sessionStore.Save(request, writer, session); err != nil {
			s.error(writer, request, http.StatusInternalServerError, err)
			return
		}

		s.respond(writer, request, http.StatusOK, nil)
	}
}

func (s *server) error(writer http.ResponseWriter, request *http.Request, code int, err error) {
	s.respond(writer, request, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(writer http.ResponseWriter, request *http.Request, code int, data interface{}) {
	writer.WriteHeader(code)
	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}
