package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yeboka/final-project/internal/app/model"
	"github.com/yeboka/final-project/internal/app/store"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// APIServer ...
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

// New ...
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("Starting Api server !")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handlerHello())
	s.router.HandleFunc("/create", s.handleCreate())
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}

func (s *APIServer) handlerHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "hello")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *APIServer) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		u, err := s.store.User().Create(&model.User{
			Email:             "mukanyerbolat@gmail.com",
			EncryptedPassword: "fgfdgfd",
		})

		res := strings.Join([]string{u.Email, strconv.Itoa(u.ID)}, " ")
		_, err = io.WriteString(w, res)
		if err != nil {
			log.Fatal(err)
		}
	}
}
