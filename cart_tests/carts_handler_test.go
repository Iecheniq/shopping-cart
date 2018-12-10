package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iecheniq/cart/cartutils"
)

func TestCartsHandler(t *testing.T) {
	testCases := []struct {
		name   string
		method string
	}{
		{name: "list all carts", method: "GET"},
		{name: "create new cart", method: "POST"},
	}
	db := cartutils.CartDB{
		Driver:     "mysql",
		DataSource: "iecheniq:HoUsE22$@tcp(localhost:3306)/shopping_cart_test",
	}
	if err := db.OpenCartDB(); err != nil {
		t.Fatalf(err.Error())
	}
	defer db.CloseCartDB()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "localhost:8080/carts", nil)
			if err != nil {
				t.Fatalf("Could not create request")
			}
			rec := httptest.NewRecorder()
			cartutils.CartsHandler(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("Expected Status OK, got %v", res.StatusCode)
			}
		})
	}
}
