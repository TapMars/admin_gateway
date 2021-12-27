package main

import (
	"TapMars/admin_proxy/pkg/productManager"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/health").HandlerFunc(healthCheck)

	var addr *string
	var opts []grpc.DialOption

	//Create a sub-router for every connected service
	productManagerRouter := router.PathPrefix("/product-manager").Subrouter()

	pm, err := productManager.NewProductManager(addr, opts)
	if err != nil {
		log.Fatalf("Product Manager connection: %v", err)
	}
	defer pm.Close()
	pm.RegisterHandlers(productManagerRouter)

	//port, host, err := config.GetEnvironmentVariables()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("Failed to find Port: %s", port)
	}

	log.Fatal(http.ListenAndServe(":"+port, router))

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}
