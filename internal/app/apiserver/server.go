package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/yeboka/final-project/internal/app/model"
	"github.com/yeboka/final-project/internal/app/store"
	"net/http"
	"time"
)

const (
	sessionName        = "canteen"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

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

	s.logger.Info("Server started successfully!")
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestId)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoAmI()).Methods("GET")
}

func (s *server) setRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			id := uuid.New().String()
			writer.Header().Set("X-Request-ID", id)
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyRequestID, id)))
		},
	)
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			logger := s.logger.WithFields(logrus.Fields{
				"remote_addr": request.RemoteAddr,
				"request_id":  request.Context().Value(ctxKeyRequestID),
			})

			logger.Infof("Started %s %s", request.Method, request.RequestURI)

			start := time.Now()
			rw := &ResponseWriter{writer, http.StatusOK}
			next.ServeHTTP(rw, request)

			logger.Infof("Completed wuth %d %s in %v", rw.code, http.StatusText(rw.code), time.Now().Sub(start))
		},
	)
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			session, err := s.sessionStore.Get(request, sessionName)
			if err != nil {
				s.error(writer, request, http.StatusInternalServerError, err)
				return
			}

			id, ok := session.Values["user_id"]

			if !ok {
				s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
				return
			}

			u, err := s.store.User().Find(id.(int))
			if err != nil {
				s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
				return
			}

			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, u)))
		},
	)
}

func (s *server) handleWhoAmI() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		s.respond(writer, request, http.StatusOK, request.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type requests struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
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
			Username: req.Username,
			Role:     "user",
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

		session.Values["user_id"] = u.ID
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
