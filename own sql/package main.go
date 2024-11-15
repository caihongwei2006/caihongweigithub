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

func init() (err error) {
	dsn := "root:cai2006@tcp(127.0.0.1)/chwuser"
	db, err := sql.Open("mysql", dsn)
	err = db.Ping()
	if err != nil {
		fmt.Println("failed to connect to database", err)
		return
	}
	fmt.Println("connected to database successfully ")
	return
}

func main() {
	err := init()
	if err != nil {
		fmt.Println("failed! failed! failed!")

		return
	}
	fmt.Println("success! now connected")

	db.Exec("insert into chwuser(name) values('五条悟')")

}
