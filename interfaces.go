package shopping

// DAO is the Data Access Object interface
// that will describe the way to read/write from a data service
//e.g. a Database or in-memory
type DAO interface {
	SetProductInventory([]*Product) error
	SetPromotions([]*Promotion) error
	SetUsers([]*User) error

	GetProducts() (Products, error)
	UpdateProducts(Products) error

	GetPromotions() (Promotions, error)

	SetUser(user *User) error
	GetUser(username string) (*User, error)
	GetUserByToken(token string) (*User, error)

	GetCart(id int) (*Cart, error)
	GetCartByUserID(userID int) (*Cart, error)
	SetCart(cart *Cart) error
}
