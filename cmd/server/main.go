package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	businessesRouter := router.PathPrefix("/businesses").Subrouter()
	itemsRouter := businessesRouter.PathPrefix("/{id}/items").Subrouter()

	businessesRouter.Methods(http.MethodGet).Path("/{id}").HandlerFunc(getBusiness)
	businessesRouter.Methods(http.MethodPost).Path("").HandlerFunc(createBusiness)
	businessesRouter.Methods(http.MethodPut).Path("/{id}").HandlerFunc(updateBusiness)
	businessesRouter.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(deleteBusiness)
	businessesRouter.Methods(http.MethodGet).Path("").HandlerFunc(searchBusinesses).
		Queries("filterDistance", "{filter-distance}", "orderBy", "{order-by}")

	itemsRouter.Methods(http.MethodGet).Path("/{id}").HandlerFunc(getItem)
	itemsRouter.Methods(http.MethodPost).Path("").HandlerFunc(createItem)
	itemsRouter.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(deleteItem)
	itemsRouter.Methods(http.MethodGet).Path("").HandlerFunc(searchItems).
		Queries("dayOfWeek", "{day-of-week}")

	//port, host, err := config.GetEnvironmentVariables()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("Failed to find Port: %s", port)
	}

	log.Fatal(http.ListenAndServe(":"+port, router))

}

func getBusiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func createBusiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func updateBusiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func deleteBusiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func searchBusinesses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func getItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func createItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func searchItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}
