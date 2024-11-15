package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	id   int
	name string
	age  int
}

var db *sql.DB

func QueryRowDemo() {
	sqlStr := "select id, name,age from chwuser where id =?"
	var u user
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Println("failed in querryrow")
		return
	}
	fmt.Printf("id:%d,name:%s,age:%d\n", u.id, u.name, u.age)

}

func initDB() (err error) {
	dsn := "root:cai2006@tcp(127.0.0.1)/test"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Failed to open a connection to the MySQL database:", err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error, Erroe,Error\n", err)
		return
	}

	fmt.Println("Opened successfully\n", dsn)
	return

}

func main() {
	err := initDB()
	if err != nil {
		fmt.Println("failed failed failed to connect to database\n")
		return
	}
	defer db.Close()

	QueryRowDemo()

}
