package shopping

const ()

// Service config for the API, different components can be attached to it
type Service struct {
	environment string
	dao         DAO
}

// Product describes the representation of products to be available in the shop
type Product struct {
	ID    int     `json:"id,omitempty"` // a number works now but it may be better to use UUIDs
	Type  int     `json:"type,omitempty"`
	Name  string  `json:"name,omitempty"`
	Stock int     `json:"stock,omitempty"`
	Price float64 `json:"price,omitempty"` //float64 can be used safely assuming the precision in the nnumber representation is not crucial
}

// Promotion describes model to have different promotions with discounts/special price
type Promotion struct {
	ProductType             int     `json:"product_type,omitempty"`
	QuantityForDiscount     int     `json:"quantity_needed,omitempty"`
	ProductsDiscounted      []int   `json:"products_discounted,omitempty"` // array of product types
	Discount                float64 `json:"discount,omitempty"`            // from 0 to 100
	QuantityForSpecialPrice int     `json:"quantity_for_special_price,omitempty"`
	SpecialPrice            float64 `json:"special_price,omitempty"` // for cases where extra items cost less
}

// Cart describes the shopping cart of a user
type Cart struct {
	ID            int                    `json:"id"`
	UserID        int                    `json:"user_id"`
	Products      map[int][]*CartProduct `json:"products"`
	Checkout      []*CartProduct         `json:"checkout,omitempty"`
	TotalPrice    float64                `json:"total_price,omitempty"`
	TotalDiscount float64                `json:"total_discount,omitempty"`
}

// User model
type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

// CartProduct describes the products in a cart
// where the key is the product ID and the value is the quantity
type CartProduct struct {
	ProductType        int     `json:"product_type,omitempty"`
	ProductName        string  `json:"product_name,omitempty"`
	UnitPrice          float64 `json:"unit_price,omitempty"`
	Discount           bool    `json:"discount,omitempty"`
	SpecialPrice       bool    `json:"special_price,omitempty"`
	DiscountPercentage float64 `json:"discount_percentage,omitempty"`
	DiscountAmount     float64 `json:"discount_amount,omitempty"`
	DiscountedPrice    float64 `json:"discounted_price,omitempty"`
}

// Products map to hold all products by id
type Products map[int]*Product

// Promotions map to hold all promotions by product type
type Promotions map[int]*Promotion

// Users map to hold all users by id
type Users map[int]*User
