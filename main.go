package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type MenuItem struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Order struct {
	Date     string `json:"date"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	OrderID  int
}

var (
	db *sql.DB
)

func main() {
	var err error

	// Get the MySQL connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	// Create a new MySQL configuration
	mysqlConfig := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  "tcp",
		Addr:                 dbHost,
		AllowNativePasswords: true,
	}

	// Create a new MySQL connection
	db, err = sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		log.Fatal("Failed to create MySQL connection:", err)
	}
	defer db.Close()

	// Create the database if it does not exist
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", DB_NAME)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}

	// Select the database
	_, err = db.Exec(fmt.Sprintf("USE %s", DB_NAME))
	if err != nil {
		log.Fatal("Failed to select database:", err)
	}

	// Create tables if they does not exist
	if err := createTables(db); err != nil {
		log.Fatalf("failed to create tables: %v", err)
	}

	// Populate menu
	if err := populateMenu(db); err != nil {
		log.Fatalf("failed to populate menu: %v", err)
	}

	log.Println("DB Schema created successfully")

	http.HandleFunc("/api/menu", handleMenu)
	http.HandleFunc("/api/orders", handleOrder)

	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("API Server is listening on 8080")
}

func createTables(db *sql.DB) error {

	// Check if the "menu" table exists

	rows, err := db.Query("SHOW TABLES LIKE 'menu'")
	if err != nil {
		return fmt.Errorf("failed to check if menu table exists: %v", err)
	}
	defer rows.Close()

	// If the table already exists, return without creating it again
	if !rows.Next() {
		// return nil
		// Create the "menu" table
		_, err = db.Exec(`
			CREATE TABLE menu (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				price DECIMAL(10,2) NOT NULL
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create menu table: %v", err)
		}
	}

	// Check if the "orders" table exists
	_, err = db.Exec("USE food_ordering")

	orderRows, err1 := db.Query("SHOW TABLES LIKE 'orders'")
	if err1 != nil {
		return fmt.Errorf("failed to check if orders table exists: %v", err1)
	}
	defer orderRows.Close()

	// If the table already exists, return without creating it again
	if !orderRows.Next() {
		log.Printf("Orders table already exist")
		// return nil
		log.Printf("Creating Orders table")
		// Create the "orders" table
		_, err = db.Exec(`
			CREATE TABLE orders (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name varchar(255) NOT NULL,
				quantity INT NOT NULL,
				order_date date
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create orders table: %v", err1)
		}
	}

	return nil
}

func populateMenu(db *sql.DB) error {
	// Check if menu items already exist
	rows, err := db.Query("SELECT COUNT(*) FROM menu")
	if err != nil {
		return fmt.Errorf("failed to query menu items: %v", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to scan menu item count: %v", err)
		}
	}

	// If menu items exist, return without populating again
	if count > 0 {
		return nil
	}

	// Populate menu items
	menuItems := []struct {
		name  string
		price float64
	}{
		{name: "Pizza", price: 10},
		{name: "Burger", price: 8},
		{name: "Salad", price: 6},
		{name: "Cofee", price: 4},
		{name: "Tea", price: 3},
	}

	// Prepare the insert statement
	stmt, err := db.Prepare("INSERT INTO menu (name, price) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	// Insert menu items
	for _, item := range menuItems {
		_, err = stmt.Exec(item.name, item.price)
		if err != nil {
			return fmt.Errorf("failed to insert menu item: %v", err)
		}
	}

	return nil
}

func handleMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// mshahid-tod-temporary allowing the address http access , we need to remove ths line once move to https access
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Fetch the menu items from the database
		rows, err := db.Query("SELECT * FROM menu")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		menuItems := make([]MenuItem, 0)

		// Iterate over the rows and populate the menu items slice
		for rows.Next() {
			var item MenuItem
			err := rows.Scan(&item.ID, &item.Name, &item.Price)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			menuItems = append(menuItems, item)
		}

		// Encode the menu items as JSON and send the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(menuItems)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var order Order
		// mshahid-tod-temporary allowing the address http access , we need to remove ths line once move to https access
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error decoding order JSON: %s", err)
			return
		}

		// Insert the order into the database
		result, err := db.Exec("INSERT INTO orders (name, quantity,order_date) VALUES (?, ?, ?)", order.Name, order.Quantity, order.Date)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the last inserted order ID
		orderID, err := result.LastInsertId()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		order.OrderID = int(orderID)

		// Encode the order details as JSON and send the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	fmt.Fprintf(w, "Order submitted successfully")
}
