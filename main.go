package main

import (
	"fmt"
	// "log"
	"net/http"
)

func main() {
	go func() {
		http.Handle("/", http.FileServer(http.Dir("./static/hello.html")))
		err := http.ListenAndServe(":8080", nil)
		// fmt.Println(err)
		fmt.Println(err)
	}()
	fmt.Println("running in 8080")

	// requestServer()
}

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("You cool"))
// }
//
// func requestServer() {
// 	resp, err := http.Get("http://localhost:8080")
// 	fmt.Println(err)
// 	defer resp.Body.Close()
// }
