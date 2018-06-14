package shopping

import log "github.com/sirupsen/logrus"

// Option ...
type Option func(*Service) error

// New returns a new server instance
func New(options ...Option) *Service {
	s := &Service{}

	for _, option := range options {
		if err := option(s); err != nil {
			log.Fatalln("server: Exiting - Failed to apply option -", err)
		}
	}

	return s
}

// SetDAO allows you to attach a DAO to the Service
func SetDAO(d DAO) Option {
	return func(s *Service) error {

		s.dao = d
		return nil
	}
}

// InitInventory calls DAO to load current stock from config
func (s *Service) InitInventory(products []*Product) error {
	return s.dao.SetProductInventory(products)

}

// InitPromotions calls DAO to load current discounts from config
func (s *Service) InitPromotions(discounts []*Promotion) error {
	return s.dao.SetPromotions(discounts)
}

// InitUsers calls DAO to load current registered users from config
func (s *Service) InitUsers(users []*User) error {
	return s.dao.SetUsers(users)
}

// PrintDAO printing data in DAO
func (s *Service) PrintDAO() {
	products, _ := s.dao.GetProducts()
	promos, _ := s.dao.GetPromotions()
	user, _ := s.dao.GetUser("test")

	log.Infof("Products: %+v\nPromotions: %+v\nUser:%+v\n", products, promos, user)
}
func (s *Service) calculatePromotions(cart *Cart) (map[int][]*CartProduct, error) {
	// var totalPrice float64
	// var discount int
	// var unitPrice map[int]int
	promotions, err := s.dao.GetPromotions()
	if err != nil {
		return nil, err
	}
	processedItems := map[int][]*CartProduct{}
	for productType, cartProducts := range cart.Products {
		itemsOfType := len(cartProducts)
		if promo, ok := promotions[productType]; ok { // is there a promotion associated to this type of product

			// calculate special price
			for k, cartProduct := range cartProducts {
				processedProduct := cartProduct
				if promo.SpecialPrice != 0 && itemsOfType >= promo.QuantityForSpecialPrice { // promo is a special price after
					// set correct price for the next product
					if k > promo.QuantityForSpecialPrice {
						processedProduct.DiscountedPrice = promo.SpecialPrice
					}
				}

				processedItems[productType] = append(processedItems[productType], processedProduct)
			}

			// calculate discounts

			for _, v := range promo.ProductsDiscounted {
				for _, productToDiscount := range cart.Products[v] {
					discountedProduct := productToDiscount
					if promo.Discount != 0 && itemsOfType >= promo.QuantityForDiscount {
						discountedProduct.DiscountedPrice = (productToDiscount.UnitPrice * promo.Discount) / 100
						discountedProduct.DiscountPercentage = promo.Discount
						discountedProduct.DiscountAmount = productToDiscount.UnitPrice - discountedProduct.DiscountedPrice
					}
					processedItems[productType] = append(processedItems[productType], discountedProduct)
				}
			}
		}

	}
	return processedItems, nil
}

func calculateTotalPrice(allProducts map[int][]*CartProduct) (float64, float64) {
	var total float64
	var totalWithDiscount float64
	for _, products := range allProducts {
		for _, p := range products {
			total += p.UnitPrice
			totalWithDiscount += p.DiscountedPrice
		}
	}
	return total, totalWithDiscount
}
