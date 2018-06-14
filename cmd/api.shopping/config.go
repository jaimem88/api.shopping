package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/jaimemartinez88/api.shopping"
)

// HTTP config
type HTTP struct {
	ListenPort string `json:"listen_port,omitempty"`
}

var config = struct {
	Environment string                `json:"environment,omitempty"`
	HTTP        *HTTP                 `json:"http,omitempty"`
	Products    []*shopping.Product   `json:"products,omitempty"`
	Promotions  []*shopping.Promotion `json:"discounts,omitempty"`
	Users       []*shopping.User      `json:"users,omitempty"`
}{
	Environment: "local",
	HTTP: &HTTP{
		ListenPort: os.Getenv("PORT"),
	},
	Products: []*shopping.Product{
		&shopping.Product{
			ID:    1,
			Type:  1,
			Name:  "Belt",
			Stock: 10,
			Price: 20,
		},
		&shopping.Product{
			ID:    2,
			Type:  2,
			Name:  "Shirt",
			Stock: 5,
			Price: 60,
		},
		&shopping.Product{
			ID:    3,
			Type:  3,
			Name:  "Suit",
			Stock: 2,
			Price: 200,
		},
		&shopping.Product{
			ID:    4,
			Type:  4,
			Name:  "Trouser",
			Stock: 4,
			Price: 70,
		},
		&shopping.Product{
			ID:    5,
			Type:  5,
			Name:  "Shoe",
			Stock: 1,
			Price: 120,
		},
		&shopping.Product{
			ID:    6,
			Type:  6,
			Name:  "Tie",
			Stock: 8,
			Price: 20,
		},
	},
	Promotions: []*shopping.Promotion{
		&shopping.Promotion{
			ProductType:         4, // trouser
			QuantityForDiscount: 2,
			ProductsDiscounted:  []int{1, 5}, // belts + shoes
			Discount:            15,
		},
		&shopping.Promotion{
			ProductType:             2, // shirt
			QuantityForSpecialPrice: 2,
			SpecialPrice:            45,
			ProductsDiscounted:      []int{6}, // shirt
			QuantityForDiscount:     3,
			Discount:                50,
		},
	},
	Users: []*shopping.User{
		&shopping.User{
			ID:       1,
			Username: "test",
			Password: "test",
		},
	},
}

func writeDefaultConfig(location string) {
	f, err := os.Create(location)
	if err != nil {
		log.Fatalln("Couldn't open", location)
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	_, err = f.Write(data)
	if err != nil {
		log.Fatalln("Couldn't write to", location)
	}
}

func loadConfig(location string) {
	raw, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatalln("Couldn't open ", location)
	}

	err = json.Unmarshal(raw, &config)
	if err != nil {
		log.Fatalln("Couldn't understand config in", location, "-", err)
	}
}
