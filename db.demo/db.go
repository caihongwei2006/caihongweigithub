package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Define the data source name (DSN)
	dsn := "root:cai2006@tcp(127.0.0.1:3306)/test" // username:password@protocol(address)/dbname

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Failed to open a connection to the MySQL database:", err)
		return
	}

	// Ping the database to ensure the connection is valid
	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to the MySQL database:", err)
		return
	}

	fmt.Println("Connected to the MySQL database!")

	// Insert "只因你太美" into the database
	_, err = db.Exec("INSERT INTO t_test (name) VALUES ('只因你太美')")
	if err != nil {
		fmt.Println("Failed to insert data into the database:", err)
		return
	}

	fmt.Println("Successfully inserted '只因你太美' into the database!")

	// Close the database connection when it's no longer needed
	defer db.Close()
}
