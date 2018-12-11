package main

import (
	"log"
	"net/http"

	"cartutils"

	_ "github.com/go-sql-driver/mysql"
	mux "github.com/gorilla/mux"
)

func main() {

	db := cartutils.CartDB{
		Driver:     "mysql",
		DataSource: "root:root@tcp(shopping_cart:3306)/shopping_cart",
	}

	if err := db.OpenCartDB(); err != nil {
		log.Fatal(err)
	}
	defer db.CloseCartDB()
	setRouts()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setRouts() {
	r := mux.NewRouter()
	r.HandleFunc("/", cartutils.HomeHandler).
		Methods("GET").
		Name("home")
	r.HandleFunc("/carts", cartutils.CartsHandler).
		Methods("GET", "POST").
		Name("carts")
	cartSubrouter := r.PathPrefix("/carts/{cartID:[0-9]+}").Subrouter()
	cartSubrouter.HandleFunc("/", cartutils.CartHandler).
		Methods("GET", "DELETE").
		Name("cart_details")
	cartSubrouter.HandleFunc("/items", cartutils.CartItemsHandler).
		Methods("GET", "POST", "DELETE").
		Name("items")
	cartSubrouter.HandleFunc("/items/{itemID:[0-9]+}", cartutils.CartItemHandler).
		Methods("GET", "DELETE", "PATCH").
		Name("item_details")
	http.Handle("/", r)
}
