package shopping

import (
	"errors"
)

// Memory implements DAO interface for the Shopping Service
type Memory struct {
	Users      Users
	Carts      map[int]*Cart
	Products   Products
	Promotions Promotions
}

// NewMemory initialises in-memory DAO for the shopping API
func NewMemory() *Memory {
	return &Memory{
		Users:      Users{},
		Carts:      map[int]*Cart{},
		Products:   Products{},
		Promotions: Promotions{},
	}
}

// SetProductInventory saves inventory in a map
func (d *Memory) SetProductInventory(products []*Product) error {
	for _, product := range products {
		d.Products[product.ID] = product
	}
	if len(d.Products) == 0 {
		return errors.New("product inventory is empty")
	}
	return nil
}

// SetPromotions saves promotions in a map, promotions could be empty
func (d *Memory) SetPromotions(promotions []*Promotion) error {
	for _, promotion := range promotions {
		d.Promotions[promotion.ProductType] = promotion
	}
	return nil
}

// SetUsers saves users in a map
func (d *Memory) SetUsers(users []*User) error {
	for _, user := range users {
		d.SetUser(user)
	}
	if len(d.Users) == 0 {
		return errors.New("user array is empty")
	}
	return nil
}

// GetProducts returns product inventory from memory
func (d *Memory) GetProducts() (Products, error) {
	if len(d.Products) == 0 {
		return nil, errors.New("no products found")
	}
	return d.Products, nil

}

// UpdateProducts updates inventory in memory
func (d *Memory) UpdateProducts(products Products) error {
	if len(products) == 0 {
		return errors.New("products empty")
	}
	d.Products = products
	return nil
}

// GetPromotions from memory
func (d *Memory) GetPromotions() (Promotions, error) {
	if len(d.Promotions) == 0 {
		return nil, errors.New("no promotions found")
	}
	return d.Promotions, nil
}

// SetUser saves data in memory
func (d *Memory) SetUser(user *User) error {
	cartID := len(d.Carts) + 1
	c := &Cart{
		ID:       cartID,
		UserID:   user.ID,
		Products: map[int][]*CartProduct{},
	}
	d.SetCart(c)
	d.Users[user.ID] = user
	return nil
}

// GetUser by username
func (d *Memory) GetUser(username string) (*User, error) {
	for _, u := range d.Users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

// GetUserByToken finds user associated to token
func (d *Memory) GetUserByToken(token string) (*User, error) {
	for _, u := range d.Users {
		if u.Token == token {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

// GetCart using its id
func (d *Memory) GetCart(id int) (*Cart, error) {
	return d.Carts[id], nil
}

// GetCartByUserID find cart corresponding to user ID
func (d *Memory) GetCartByUserID(userID int) (*Cart, error) {
	for _, c := range d.Carts {
		if c.UserID == userID {
			return c, nil
		}
	}
	return nil, errors.New("cart not found")
}

// SetCart Saves cart data
func (d *Memory) SetCart(cart *Cart) error {
	d.Carts[cart.ID] = cart
	return nil
}
