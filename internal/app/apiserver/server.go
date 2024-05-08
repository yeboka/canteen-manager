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
	"strconv"
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
	s.router.HandleFunc("/category", s.handleCategoriesGet()).Methods("GET")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/orders", s.handleCreateOrder()).Methods("POST")
	private.HandleFunc("/orders/{id}", s.handleDeleteOrder()).Methods("DELETE")
	private.HandleFunc("/allMyOrders", s.handleGetAllOrders()).Methods("GET")
	private.HandleFunc("/whoami", s.handleWhoAmI()).Methods("GET")
	private.HandleFunc("/users/{id}", s.handleUserUpdate()).Methods("PATCH")

	admin := s.router.PathPrefix("/admin").Subrouter()
	admin.Use(s.authenticateUser)
	admin.Use(s.checkAdmin)
	admin.HandleFunc("/users/{id}/role", s.handleRoleChange()).Methods("PATCH")
	admin.HandleFunc("/menu-item/{id}", s.handleMenuItemUpdate()).Methods("PATCH")
	admin.HandleFunc("/menu-item/{id}", s.handleMenuItemDelete()).Methods("DELETE")
	admin.HandleFunc("/users/{id}", s.handleDeleteUser()).Methods("DELETE")
	admin.HandleFunc("/menu-item", s.handleMenuItemCreate()).Methods("POST")
	admin.HandleFunc("/category", s.handleCategoryCreate()).Methods("POST")
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

func (s *server) checkAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ctxKeyUser).(*model.User)
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errors.New("unauthorized access: missing user information"))
			return
		}

		if user.Role != "admin" {
			s.error(w, r, http.StatusForbidden, errors.New("insufficient privileges: requires admin role"))
			return
		}

		s.logger.Printf("user %s accessed admin resource (id: %d)", user.Username, user.ID)

		next.ServeHTTP(w, r)
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			userId, _ := s.getUserId(writer, request)

			u, err := s.store.User().Find(userId)
			if err != nil {
				s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
				return
			}

			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, u)))
		},
	)
}

func (s *server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, errors.New("missing user ID in URL"))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid user ID"))
			return
		}

		_, err = s.store.User().Find(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if err := s.store.User().Delete(id); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]string{"message": "User deleted successfully"})
	}
}

func (s *server) handleRoleChange() http.HandlerFunc {
	type roleChangeRequest struct {
		Role string `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req roleChangeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		vars := mux.Vars(r)
		userIDStr, ok := vars["id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, errors.New("missing user ID in URL"))
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid user ID"))
			return
		}

		user, err := s.store.User().Find(userID)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		user.Role = req.Role
		if err := s.store.User().UpdateRole(userID, req.Role); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]string{"message": "User role updated successfully"})
	}
}

func (s *server) handleUserUpdate() http.HandlerFunc {
	type requests struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		idStr, ok := vars["id"]
		if !ok {
			s.error(writer, request, http.StatusBadRequest, errors.New("missing user ID in URL"))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(writer, request, http.StatusBadRequest, errors.New("invalid user ID"))
			return
		}

		req := &requests{}
		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().Find(id)
		if err != nil {
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}

		u.Email = req.Email
		u.Username = req.Username

		if err := s.store.User().Update(u); err != nil {
			s.error(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(writer, request, http.StatusOK, "updated")
	}
}

type CategoryTree struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	MenuItems []*model.MenuItem `json:"menu_items"`
	Children  []*CategoryTree   `json:"children,omitempty"`
}

func (s *server) handleCategoriesGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		categories, err := s.store.Category().GetAllCategories()
		if err != nil {
			s.error(writer, request, http.StatusInternalServerError, err)
			return
		}

		categoryMap := make(map[int]*CategoryTree)

		for _, category := range categories {
			items, err := s.store.MenuItem().FindByCategoryId(category.ID)
			if err != nil {
				s.error(writer, request, http.StatusInternalServerError, err)
				return
			}
			categoryMap[category.ID] = &CategoryTree{
				ID:        category.ID,
				Name:      category.Name,
				MenuItems: items,
			}
		}

		var roots []*CategoryTree
		s.logger.Info(categories)
		s.logger.Info(categoryMap)
		for _, category := range categories {
			if category.ParentID == -1 {
				roots = append(roots, categoryMap[category.ID])
			} else {
				parent := categoryMap[category.ParentID]
				if parent != nil {
					parent.Children = append(parent.Children, categoryMap[category.ID])
				}
			}
		}

		s.respond(writer, request, http.StatusOK, roots)
	}
}

func (s *server) handleMenuItemCreate() http.HandlerFunc {
	type requests struct {
		Name        string `json:"name"`
		CategoryId  int    `json:"categoryId"`
		Price       int    `json:"price"`
		Description string `json:"description"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		req := &requests{}
		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}

		mi := &model.MenuItem{
			Name:        req.Name,
			CategoryID:  req.CategoryId,
			Price:       req.Price,
			Description: req.Description,
		}

		if err := s.store.MenuItem().Create(mi); err != nil {
			s.error(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(writer, request, http.StatusCreated, mi)
	}
}

func (s *server) handleMenuItemDelete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		idStr, ok := vars["id"]
		if !ok {
			s.error(writer, request, http.StatusBadRequest, errors.New("missing user ID in URL"))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(writer, request, http.StatusBadRequest, errors.New("invalid user ID"))
			return
		}

		if err := s.store.MenuItem().Delete(id); err != nil {
			s.error(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(writer, request, http.StatusOK, id)
	}
}

func (s *server) handleMenuItemUpdate() http.HandlerFunc {
	type requests struct {
		Name        string `json:"name"`
		Price       int    `json:"price"`
		Description string `json:"description"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		req := &requests{}
		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}
		vars := mux.Vars(request)
		idStr, ok := vars["id"]
		if !ok {
			s.error(writer, request, http.StatusBadRequest, errors.New("missing user ID in URL"))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(writer, request, http.StatusBadRequest, errors.New("invalid user ID"))
			return
		}

		mi := &model.MenuItem{
			ID:          id,
			Name:        req.Name,
			Price:       req.Price,
			Description: req.Description,
		}

		if err := s.store.MenuItem().Update(mi); err != nil {
			s.error(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(writer, request, http.StatusOK, mi)
	}
}

func (s *server) handleCategoryCreate() http.HandlerFunc {
	type requests struct {
		Name     string `json:"name"`
		ParentId int    `json:"parentId,omitempty"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		req := &requests{}
		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}

		ctg := &model.Category{
			Name:     req.Name,
			ParentID: req.ParentId,
		}
		s.logger.Info(req)
		s.logger.Info(ctg)
		if req.ParentId > 0 {
			parentCategory, err := s.store.Category().Find(req.ParentId)
			if err != nil {
				s.error(writer, request, http.StatusBadRequest, err)
				return
			} else {
				ctg.ParentID = parentCategory.ID
			}
		}

		if err := s.store.Category().Create(ctg); err != nil {
			s.error(writer, request, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(writer, request, http.StatusCreated, ctg)
	}
}

func (s *server) handleWhoAmI() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		s.respond(writer, request, http.StatusOK, request.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handleCreateOrder() http.HandlerFunc {
	type respondOrder struct {
		Id         int                `json:"id"`
		OrderItems []*model.OrderItem `json:"order_item"`
		CreatedAt  time.Time          `json:"created_At"`
		TotalPrice int                `json:"total_price"`
	}

	type requests struct {
		MenuItemId []int `json:"menu_item_id"`
		Quantity   []int `json:"quantity"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		userId, _ := s.getUserId(writer, request)

		req := &requests{}
		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}

		var totalAmount int
		for i := 0; i < len(req.MenuItemId); i++ {
			price := s.store.MenuItem().GetPrice(req.MenuItemId[i])
			totalAmount += price * req.Quantity[i]
		}

		o := &model.Order{
			UserId:      userId,
			TotalAmount: totalAmount,
		}

		exception := s.store.Order().Create(o)
		if exception != nil {
			s.error(writer, request, http.StatusUnprocessableEntity, exception)
			return
		}

		var orderItems []*model.OrderItem

		for i := 0; i < len(req.MenuItemId); i++ {
			mi := &model.OrderItem{
				OrderId:    o.ID,
				MenuItemId: req.MenuItemId[i],
				Quantity:   req.Quantity[i],
			}

			if err := s.store.OrderItem().Create(mi); err != nil {
				if e := s.store.Order().Delete(o.ID); e != nil {
					s.error(writer, request, http.StatusInternalServerError, e)
					return
				}
				s.error(writer, request, http.StatusUnprocessableEntity, err)
				return
			}

			orderItems = append(orderItems, mi)
		}

		respondOrder := respondOrder{
			Id:         o.ID,
			CreatedAt:  o.CreatedAt,
			TotalPrice: totalAmount,
			OrderItems: orderItems,
		}

		s.respond(writer, request, http.StatusCreated, respondOrder)
	}
}

func (s *server) handleGetAllOrders() http.HandlerFunc {
	type respondOrder struct {
		Id         int                `json:"id"`
		OrderItems []*model.OrderItem `json:"order_item"`
		CreatedAt  time.Time          `json:"created_At"`
		TotalPrice int                `json:"total_price"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		userId, err := s.getUserId(writer, request)
		if err != nil {
			s.error(writer, request, http.StatusUnauthorized, err)
			return
		}

		orders, err := s.store.Order().GetOrders(userId)
		if err != nil {
			s.error(writer, request, http.StatusInternalServerError, err)
			return
		}

		var respondOrders []respondOrder
		for _, order := range orders {
			orderItems, err := s.store.OrderItem().GetOrderItems(order.ID)
			if err != nil {
				s.error(writer, request, http.StatusInternalServerError, err)
				return
			}

			respondOrder := respondOrder{
				Id:         order.ID,
				CreatedAt:  order.CreatedAt,
				TotalPrice: order.TotalAmount,
				OrderItems: orderItems,
			}
			respondOrders = append(respondOrders, respondOrder)
		}

		s.respond(writer, request, http.StatusOK, respondOrders)
	}
}

func (s *server) handleDeleteOrder() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		orderId := vars["id"]

		id, err := strconv.Atoi(orderId)
		if err != nil {
			s.error(writer, request, http.StatusBadRequest, err)
			return
		}

		ex := s.store.OrderItem().DeleteAllOrder(id)
		if ex != nil {
			s.error(writer, request, http.StatusBadRequest, ex)
			return
		}

		e := s.store.Order().Delete(id)
		if e != nil {
			s.error(writer, request, http.StatusBadRequest, e)
			return
		}

		s.respond(writer, request, http.StatusOK, nil)
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

func (s *server) getUserId(writer http.ResponseWriter, request *http.Request) (int, error) {
	session, err := s.sessionStore.Get(request, sessionName)

	if err != nil {
		s.error(writer, request, http.StatusInternalServerError, err)
		return 0, err
	}

	userId, authorized := session.Values["user_id"]

	if !authorized {
		s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
		return 0, nil
	}

	return userId.(int), nil
}
