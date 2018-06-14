package shopping

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") // pool of numbers and letters to generate random code

)

// HandleHealthcheck responds with an empty JSON object
func (s *Service) HandleHealthcheck(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, struct{}{})
}

// HandleNotFound writes error and logs the requested method
func (s *Service) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, &Error{Code: http.StatusNotFound, Message: r.RequestURI + " not found"})
}

// HandleLogin validates credentials for a user
func (s *Service) HandleLogin(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		s.writeError(w, errInvalidUsernameOrPassowrd.msg("Authorization header missing"))
		return
	}

	user, err := s.dao.GetUser(username)
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetUser: "+err.Error()))
		return
	}

	if user == nil {
		s.writeError(w, errInvalidUsernameOrPassowrd.msg("User not found: "+username))
		return
	}
	if password != user.Password {
		s.writeError(w, errInvalidUsernameOrPassowrd.msg("Incorrect password"))
		return
	}

	token := rand10letterString()
	user.Token = token

	if err := s.dao.SetUser(user); err != nil {
		s.writeError(w, errInternalServerError.msg("dao.SetUser: "+err.Error()))
		return
	}
	// init a cart if it's empty
	if _, err := s.dao.GetCartByUserID(user.ID); err != nil {
		s.dao.SetCart(&Cart{UserID: user.ID, Products: map[int][]*CartProduct{}})
	}

	s.writeJSON(w, struct {
		Token string `json:"token"`
	}{Token: token})
}

// HandleGetProducts gets a list of products
func (s *Service) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.dao.GetProducts()
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetCartByUserID: "+err.Error()))
		return
	}
	res := struct {
		Products Products `json:"products"`
	}{
		Products: products,
	}
	s.writeJSON(w, res)
}

// HandleGetPromotions gets a list of promotions
func (s *Service) HandleGetPromotions(w http.ResponseWriter, r *http.Request) {
	promos, err := s.dao.GetPromotions()
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetCartByUserID: "+err.Error()))
		return
	}
	res := struct {
		Promotions Promotions `json:"promotions"`
	}{
		Promotions: promos,
	}
	s.writeJSON(w, res)
}

// HandleGetCart gets the latest shopping cart information for a user
func (s *Service) HandleGetCart(w http.ResponseWriter, r *http.Request) {
	user, ctxErr := getUserFromContext(r.Context())
	if ctxErr != nil {
		s.writeError(w, ctxErr)
		return
	}

	cart, err := s.dao.GetCartByUserID(user.ID)
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetCartByUserID: "+err.Error()))
		return
	}
	s.writeJSON(w, struct {
		Cart *Cart `json:"cart"`
	}{Cart: cart})
}

// HandleCartAddItem adds an item to the cart
func (s *Service) HandleCartAddItem(w http.ResponseWriter, r *http.Request) {
	req := struct {
		ProductID int `json:"product_id,omitempty"`
		Quantity  int `json:"quantity,omitempty"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, errInternalServerError.msg("failed to decode req: "+err.Error()))
		return
	}

	user, ctxErr := getUserFromContext(r.Context())
	if ctxErr != nil {
		s.writeError(w, ctxErr)
		return
	}

	cart, err := s.dao.GetCartByUserID(user.ID)
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetCartByUserID: "+err.Error()))
		return
	}
	products, err := s.dao.GetProducts()
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetProducts: "+err.Error()))
		return
	}
	p := products[req.ProductID]
	if p == nil {
		s.writeError(w, errInternalServerError.msg("product type doesn't exist"))
		return
	}
	if req.Quantity > p.Stock {
		s.writeError(w, errBadRequestNotEnoughStock.msg("not enough stock"))
		return
	}
	for i := 0; i < req.Quantity; i++ {
		cart.Products[req.ProductID] = append(cart.Products[req.ProductID], &CartProduct{
			ProductType: p.Type,
			ProductName: p.Name,
			UnitPrice:   p.Price,
		})
	}
	if len(cart.Products[req.ProductID]) > p.Stock {
		cart.Products[req.ProductID] = append(cart.Products[req.ProductID][:0], cart.Products[req.ProductID][1:]...)
		s.writeError(w, errBadRequestNotEnoughStock.msg("not enough stock"))
		return
	}

	if err := s.dao.SetCart(cart); err != nil {
		s.writeError(w, errInternalServerError.msg("dao.SetCart: "+err.Error()))
		return
	}
	s.writeJSON(w, struct {
		Cart *Cart `json:"cart"`
	}{Cart: cart})
}

// HandleCartRemoveItem removes an item fromthe cart
func (s *Service) HandleCartRemoveItem(w http.ResponseWriter, r *http.Request) {
	req := struct {
		ProductID int `json:"product_id,omitempty"`
		Quantity  int `json:"quantity,omitempty"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, errInternalServerError.msg("failed to decode req: "+err.Error()))
		return
	}

	user, ctxErr := getUserFromContext(r.Context())
	if ctxErr != nil {
		s.writeError(w, ctxErr)
		return
	}

	cart, err := s.dao.GetCartByUserID(user.ID)
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetCartByUserID: "+err.Error()))
		return
	}
	for i := 0; i < req.Quantity; i++ {
		if len(cart.Products[req.ProductID]) < 2 {
			cart.Products[req.ProductID] = []*CartProduct{}
		} else {
			cart.Products[req.ProductID] = append(cart.Products[req.ProductID][:0], cart.Products[req.ProductID][1:]...)
		}
	}
	// reset all promotions
	for _, products := range cart.Products {
		for _, product := range products {
			product.Discount = false
			product.SpecialPrice = false
			product.DiscountAmount = 0
			product.DiscountPercentage = 0
			product.DiscountedPrice = 0

		}
	}
	if err := s.dao.SetCart(cart); err != nil {
		s.writeError(w, errInternalServerError.msg("dao.SetCart: "+err.Error()))
		return
	}
	s.writeJSON(w, struct {
		Cart *Cart `json:"cart"`
	}{Cart: cart})
}

// HandleCartClear removes all items from a cart
func (s *Service) HandleCartClear(w http.ResponseWriter, r *http.Request) {
	user, ctxErr := getUserFromContext(r.Context())
	if ctxErr != nil {
		s.writeError(w, ctxErr)
		return
	}

	cart, err := s.dao.GetCartByUserID(user.ID)
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetCartByUserID: "+err.Error()))
		return
	}
	cart.Products = map[int][]*CartProduct{}
	cart.Checkout = []*CartProduct{}

	if err := s.dao.SetCart(cart); err != nil {
		s.writeError(w, errInternalServerError.msg("dao.SetCart: "+err.Error()))
		return
	}
	s.writeJSON(w, struct {
		Cart *Cart `json:"cart"`
	}{Cart: cart})
}

// HandleCartCheckout calculates final price before buying
func (s *Service) HandleCartCheckout(w http.ResponseWriter, r *http.Request) {
	user, ctxErr := getUserFromContext(r.Context())
	if ctxErr != nil {
		s.writeError(w, ctxErr)
		return
	}

	cart, err := s.dao.GetCartByUserID(user.ID)
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetCartByUserID: "+err.Error()))
		return
	}

	updatedProducts, err := s.calculatePromotions(cart)
	if err != nil {
		s.writeError(w, errInternalServerError.msg("s.calculatePromotions: "+err.Error()))
		return
	}
	total, discountedTotal := calculateTotalPrice(updatedProducts)
	cart.Checkout = updatedProducts
	cart.TotalPrice = total
	cart.TotalDiscount = discountedTotal

	if err := s.dao.SetCart(cart); err != nil {
		s.writeError(w, errInternalServerError.msg("dao.SetCart: "+err.Error()))
		return
	}

	res := struct {
		Cart *Cart `json:"cart"`
	}{cart}
	s.writeJSON(w, res)
}

// HandleCartBuy process items as bought and removes them from inventory
func (s *Service) HandleCartBuy(w http.ResponseWriter, r *http.Request) {
	user, ctxErr := getUserFromContext(r.Context())
	if ctxErr != nil {
		s.writeError(w, ctxErr)
		return
	}

	cart, err := s.dao.GetCartByUserID(user.ID)
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetCartByUserID: "+err.Error()))
		return
	}
	inventory, err := s.dao.GetProducts()
	if err != nil {
		s.writeError(w, errInternalServerError.msg("dao.GetCartByUserID: "+err.Error()))
		return
	}

	for productType, products := range cart.Products {
		inventory[productType].Stock -= len(products)
		if inventory[productType].Stock < 0 {
			inventory[productType].Stock = 0
		}
	}

	if err := s.dao.UpdateProducts(inventory); err != nil {
		s.writeError(w, errInternalServerError.msg("dao.UpdateProducts: "+err.Error()))
		return
	}
	s.writeJSON(w, struct{}{})
}

func (s *Service) writeError(w http.ResponseWriter, e *Error) {

	log.WithError(e).WithField("error-message", e.message).Error()

	js, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	w.Write(js)
}

func (s *Service) writeJSON(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	js, err := json.Marshal(i)
	if err != nil {
		s.writeError(w, errInternalServerError.msg("json.Marshal: "+err.Error()))
		return
	}
	w.Write(js)
}

func rand10letterString() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getUserFromContext(ctx context.Context) (*User, *Error) {
	userData := ctx.Value(ctxUser)
	if userData == nil {
		return nil, errInternalServerError.msg("missing user in context")

	}
	user, ok := userData.(*User)
	if !ok {
		return nil, errInternalServerError.msg("failed to get user data")

	}
	return user, nil
}
