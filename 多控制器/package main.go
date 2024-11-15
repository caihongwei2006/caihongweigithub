package main

import (
	"fmt"
	"net/http"
)

type Myhandler struct{}

func (m *Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("返回了数据喔！"))
	userAgent := r.Header.Get("User-Agent")
	fmt.Fprintln(w, "User-Agent:", userAgent)

}

func main() {
	fmt.Println("Hello, World!")
	h := Myhandler{}
	http.Handle("/ghost", &h)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
