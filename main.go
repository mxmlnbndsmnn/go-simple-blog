package main

import (
	"log"
	"net/http"
	"simple-blog/database"
	"simple-blog/handlers"
)

func main() {
	database.InitDatabase()
	defer database.CloseDatabase()

	// define routes and handlers
	http.HandleFunc("/blogs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetBlogs(w, r)
		case http.MethodPost:
			handlers.CreateBlog(w, r)
		case http.MethodDelete:
			handlers.DeleteBlog(w, r)
		case http.MethodPut:
			handlers.UpdateBlog(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// start the HTTP server
	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
