package cartutils

import (
	"log"
	"strconv"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type CartDB struct {
	Driver     string
	DataSource string
}

//Item is the model used for articles in a cart
type Item struct {
	Id     int
	Title  string
	Price  float64
	CartID int
}

//Cart is the model that describes a shopping cart
type Cart struct {
	Id    int
	Items []Item
}

func (database *CartDB) OpenCartDB() error {
	db, _ = sql.Open(database.Driver, database.DataSource)

	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (database *CartDB) CloseCartDB() {
	db.Close()
}

func GetAllCarts() (map[string]Cart, error) {
	carts := make(map[string]Cart)
	rows, err := db.Query("SELECT * FROM carts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		cart := Cart{
			Id:    0,
			Items: make([]Item, 0),
		}
		err := rows.Scan(&cart.Id)
		if err != nil {
			return nil, err
		}
		carts["Cart "+strconv.Itoa(cart.Id)] = cart
	}
	rows, err = db.Query("SELECT * FROM items")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Item{
			Id:     0,
			Title:  "",
			Price:  0.0,
			CartID: 0,
		}
		rows.Scan(&item.Id, &item.Title, &item.Price, &item.CartID)
		cart := carts["Cart "+strconv.Itoa(int(item.CartID))]
		cart.Items = append(cart.Items, item)
		carts["Cart "+strconv.Itoa(int(item.CartID))] = cart
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func CreateCart() error {
	res, err := db.Exec("INSERT INTO carts() VALUES()")
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)
	return nil
}

func (cart *Cart) GetCartById() error {
	item := Item{
		Id:     0,
		Title:  "",
		Price:  0.0,
		CartID: 0,
	}
	err := db.QueryRow("SELECT * FROM carts WHERE id = ? ", cart.Id).Scan(&cart.Id)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("SELECT * FROM items WHERE cart_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(cart.Id)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&item.Id, &item.Title, &item.Price, &item.CartID)
		if err != nil {
			return err
		}
		cart.Items = append(cart.Items, item)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (cart *Cart) DeleteCartById() error {
	res, err := db.Exec("DELETE FROM carts WHERE id = ? ", cart.Id)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("%d rows affected\n", rowCnt)
	return nil
}

func (cart *Cart) GetAllCartItems() error {
	item := Item{
		Id:     0,
		Title:  "",
		Price:  0.0,
		CartID: 0,
	}
	stmt, err := db.Prepare("SELECT * FROM items WHERE cart_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(cart.Id)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&item.Id, &item.Title, &item.Price, &item.CartID)
		cart.Items = append(cart.Items, item)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (cart *Cart) CreateCartItem(item Item) error {
	stmt, err := db.Prepare("INSERT INTO items (id, title, price, cart_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(item.Id, item.Title, item.Price, cart.Id)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("affected = %d\n", rowCnt)
	return nil
}

func (cart *Cart) DeleteAllCartItems() error {
	stmt, err := db.Prepare("DELETE FROM items WHERE cart_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(cart.Id)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("affected = %d\n", rowCnt)
	return nil
}

func (cart *Cart) GetCartItem(id int) (Item, error) {
	item := Item{
		Id:     0,
		Title:  "",
		Price:  0.0,
		CartID: 0,
	}
	err := db.QueryRow("SELECT * FROM items WHERE id = ? AND cart_id = ?", id, cart.Id).
		Scan(&item.Id, &item.Title, &item.Price, &item.CartID)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (cart *Cart) DeleteCartItem(id int) error {
	stmt, err := db.Prepare("DELETE FROM items WHERE id = ? AND cart_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(id, cart.Id)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("affected = %d\n", rowCnt)
	return nil
}

func (cart *Cart) UpdateCartItemPrice(price float64, id int) error {
	stmt, err := db.Prepare("UPDATE items SET price = ? WHERE id = ? AND cart_id = ? ")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(price, id, cart.Id)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("affected = %d\n", rowCnt)
	return nil
}
