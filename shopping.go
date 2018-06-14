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
func (s *Service) calculatePromotions(cart *Cart) ([]*CartProduct, error) {

	promotions, err := s.dao.GetPromotions()
	if err != nil {
		return nil, err
	}
	checkOut := []*CartProduct{}

	checkPromos(promotions, cart.Products)

	for _, productsArray := range cart.Products {
		for _, product := range productsArray {
			if product.Discount {
				product.DiscountedPrice = (product.UnitPrice * product.DiscountPercentage) / 100
				// product.DiscountPercentage = promo.Discount
				product.DiscountAmount = product.UnitPrice - product.DiscountedPrice
			}
			if product.SpecialPrice {
				product.DiscountAmount = product.UnitPrice - product.DiscountedPrice
			}
			checkOut = append(checkOut, product)
		}

	}

	return checkOut, nil
}

func calculateTotalPrice(allProducts []*CartProduct) (float64, float64) {
	var total float64
	var totalDiscount float64

	for _, p := range allProducts {
		total += p.UnitPrice
		totalDiscount += p.DiscountAmount
	}

	return total - totalDiscount, totalDiscount
}

func checkPromos(promotions Promotions, productsMap map[int][]*CartProduct) {
	// productsForDiscount := []int{}
	for productType, products := range productsMap {
		itemsOfType := len(products)
		promo, ok := promotions[productType]
		if ok {
			if promo.SpecialPrice != 0 && itemsOfType >= promo.QuantityForSpecialPrice { // promo
				for k, p := range products {
					if k+1 > promo.QuantityForSpecialPrice {
						p.SpecialPrice = true
						p.DiscountedPrice = promo.SpecialPrice
					}
				}
			}
			if itemsOfType >= promo.QuantityForDiscount {
				for _, v := range promo.ProductsDiscounted {
					for _, product := range productsMap[v] {

						product.Discount = true
						product.DiscountPercentage = promo.Discount

					}
				}
			}
		}

	}

}
