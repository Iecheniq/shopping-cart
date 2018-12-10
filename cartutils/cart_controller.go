package cartutils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	mux "github.com/gorilla/mux"
)

//HomeHandler handles the root URL
//GET: Description of the shopping cart API
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if _, err := w.Write([]byte(fmt.Sprintf("Welcome to the shopping cart API"))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//CartsHandler handles all existing carts
//GET: List all carts
//POST: Create a new cart
func CartsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		carts, err := GetAllCarts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Jcarts, err := json.Marshal(carts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(Jcarts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else if r.Method == "POST" {

		if err := CreateCart(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write([]byte(fmt.Sprintf("Cart has been created"))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Allowed methods: GET, POST", http.StatusMethodNotAllowed)
		return
	}
}

//CartHandler handles a cart.
//GET: Get a specific cart
//DELETE: Delete the cart
func CartHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cartID, _ := strconv.Atoi(vars["cartID"])
	cart := Cart{
		Id:    cartID,
		Items: make([]Item, 0),
	}
	if r.Method == "GET" {
		if err := cart.GetCartById(); err != nil {
			if err == sql.ErrNoRows {
				if _, err := w.Write([]byte(fmt.Sprintf("Cart with ID %v does not exist", vars["cartID"]))); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

				}
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		Jcart, err := json.Marshal(cart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(Jcart); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "DELETE" {

		if err := cart.DeleteCartById(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write([]byte(fmt.Sprintf("Cart with ID %v has been deleted", vars["cartID"]))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Allowed methods: GET, DELETE", http.StatusMethodNotAllowed)
		return
	}
}

//CartItemsHandler handles the items of a cart.
//GET: Get all items of a cart
//POST: Add a new item to the cart
//DELETE: Delete all items of a cart
func CartItemsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID, _ := strconv.Atoi(vars["cartID"])
	cart := Cart{
		Id:    cartID,
		Items: make([]Item, 0),
	}
	if r.Method == "GET" {
		if err := cart.GetAllCartItems(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jItems, err := json.Marshal(cart.Items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(jItems); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ID := r.Form.Get("id")
		title := r.Form.Get("title")
		price := r.Form.Get("price")
		if ID == "" || title == "" || price == "" {
			if _, err := w.Write([]byte(fmt.Sprintf("One or more fileds are empty, you must enter title and price"))); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			ID, _ := strconv.Atoi(ID)
			price, _ := strconv.ParseFloat(price, 64)
			item := Item{
				Id:    ID,
				Title: title,
				Price: price,
			}
			if err := cart.CreateCartItem(item); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := w.Write([]byte(fmt.Sprintf("Item created in  cart %v", cart.Id))); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else if r.Method == "DELETE" {
		if err := cart.DeleteAllCartItems(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write([]byte(fmt.Sprintf("All items deleted from Cart %v", cart.Id))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Allowed methods: GET, POST, DELETE", http.StatusMethodNotAllowed)
		return
	}
}

//CartItemHandler handles a specific item in a cart
//GET: Get a specific item
//DELETE: Delete item
//PATCH: Modify item price
func CartItemHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cartID, _ := strconv.Atoi(vars["cartID"])
	cart := Cart{
		Id:    cartID,
		Items: make([]Item, 0),
	}
	itemID, _ := strconv.Atoi(vars["itemID"])

	if r.Method == "GET" {
		item, err := cart.GetCartItem(itemID)
		if err != nil {
			if err == sql.ErrNoRows {
				if _, err := w.Write([]byte(fmt.Sprintf("Cart %v has no item with ID %v", cart.Id, itemID))); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Jitem, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(Jitem); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "DELETE" {
		if err := cart.DeleteCartItem(itemID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write([]byte(fmt.Sprintf("Item with ID %v has been deleted from Cart %v", itemID, cartID))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "PATCH" {
		r.ParseForm()
		price := r.Form.Get("price")
		if price == "" {
			if _, err := w.Write([]byte(fmt.Sprintf("Price must not be empty and must be a float"))); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		itemPrice, err := strconv.ParseFloat(price, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cart.UpdateCartItemPrice(itemPrice, itemID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write([]byte(fmt.Sprintf("Item price has been updated"))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Allowed methods: GET, DELETE, PATCH", http.StatusMethodNotAllowed)
		return
	}
}
