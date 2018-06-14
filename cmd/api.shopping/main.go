package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	shopping "github.com/jaimemartinez88/api.shopping"
	log "github.com/sirupsen/logrus"
)

const (
	httpServerReadTimeout  = 3 * time.Second
	httpServerWriteTimeout = 120 * time.Second
)

var (
	version            string
	corsAllowedHeaders = handlers.AllowedHeaders([]string{"*"})
	corsAllowedDomains = handlers.AllowedOrigins([]string{
		"*",
	})
	corsAllowedMethods = handlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions})
)

func main() {
	// flag handling
	defaultLocation := flag.String("default", "", "location to write a default configuration to (this will overwrite an existing file at this location)")
	configLocation := flag.String("config", "", "JSON config file to load")

	flag.Parse()
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	log.SetLevel(log.DebugLevel)
	if *defaultLocation == "" && *configLocation == "" {
		log.Println("Using default config:")
		data, _ := json.MarshalIndent(config, "", "  ")
		io.Copy(os.Stdout, bytes.NewReader(data))
		fmt.Printf("\n")
	} else if *defaultLocation != "" {
		writeDefaultConfig(*defaultLocation)
		os.Exit(0)
	} else if *configLocation != "" {
		loadConfig(*configLocation)
	}

	dao := shopping.NewMemory()
	service := shopping.New(
		shopping.SetDAO(dao),
	)

	if err := service.InitInventory(config.Products); err != nil {
		log.WithError(err).Fatalln("failed to set inventory")
	}
	if err := service.InitPromotions(config.Promotions); err != nil {
		log.WithError(err).Fatalln("failed to set discounts")
	}
	if err := service.InitUsers(config.Users); err != nil {
		log.WithError(err).Fatalln("failed to set discounts")
	}
	service.PrintDAO()
	r := mux.NewRouter()
	r.Use(service.LoggingMiddleware)
	r.NotFoundHandler = http.HandlerFunc(service.HandleNotFound)

	r.HandleFunc("/healthcheck", service.HandleHealthcheck).Methods(http.MethodGet)

	v1 := r.PathPrefix("/v1/shopping").Subrouter()

	v1.HandleFunc("/login", service.HandleLogin).Methods(http.MethodPost)
	v1.HandleFunc("/products", service.HandleGetProducts).Methods(http.MethodGet)
	v1.HandleFunc("/promotions", service.HandleGetPromotions).Methods(http.MethodGet)
	v1Secure := v1.NewRoute().Subrouter()
	v1Secure.Use(service.ValidateAccessToken)
	v1Secure.HandleFunc("/cart", service.HandleGetCart).Methods(http.MethodGet)
	v1Secure.HandleFunc("/cart/add", service.HandleCartAddItem).Methods(http.MethodPost)
	v1Secure.HandleFunc("/cart/remove", service.HandleCartRemoveItem).Methods(http.MethodPost)
	v1Secure.HandleFunc("/cart/clear", service.HandleCartClear).Methods(http.MethodPost)
	v1Secure.HandleFunc("/cart/checkout", service.HandleCartCheckout).Methods(http.MethodGet)
	v1Secure.HandleFunc("/cart/buy", service.HandleCartBuy).Methods(http.MethodPost)

	handler := handlers.CORS(corsAllowedHeaders, corsAllowedDomains, corsAllowedMethods)(r)
	port := config.HTTP.ListenPort
	if port == "" { // default to port 8080
		port = "8080"
	}
	httpServer := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  httpServerReadTimeout,
		WriteTimeout: httpServerWriteTimeout,
		Handler:      handler,
	}
	log.Infof("Server listening on port :%s", port)
	if err := httpServer.ListenAndServe(); nil != err {
		log.Fatalln("Failed to start server", err)
	}

}
