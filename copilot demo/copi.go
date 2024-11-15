package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var count int = 1

func handleconnection(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("现在立刻马上，我在连接  For the %dth time", count)
	t, err := template.ParseFiles("html/substance.html")
	if err != nil {
		http.Error(w, "模板解析错误", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
	count++
}

func main() {
	fmt.Println("Setting up static file server...")
	http.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	fmt.Println("Setting up handle connection...")
	http.HandleFunc("/handle", handleconnection)

	fmt.Println("Starting server on :8090...")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}
