package handlers

import (
	// _ "database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"simple-blog/database"
	"simple-blog/models"
	"strconv"
	"strings"
	"time"
)

// helper function to write JSON responses
func jsonRespose(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func GetBlogs(w http.ResponseWriter, r *http.Request) {
	// collect filters from query parameters
	filters := []string{}
	args := []interface{}{}

	// TODO to be replaced by author_id
	authorFilter := r.URL.Query().Get("author")
	if authorFilter != "" {
		filters = append(filters, "author = ?")
		args = append(args, authorFilter)
	}

	titleFilter := r.URL.Query().Get("title")
	if titleFilter != "" {
		filters = append(filters, "title LIKE ?")
		args = append(args, "%" + titleFilter + "%")
	}

	textFilter := r.URL.Query().Get("text")
	if textFilter != "" {
		filters = append(filters, "text LIKE ?")
		args = append(args, "%" + textFilter + "%")
	}

	startTimeFilter := r.URL.Query().Get("start_time")
	endTimeFilter := r.URL.Query().Get("end_time")
	if startTimeFilter != "" && endTimeFilter != "" {
		filters = append(filters, "creation_time BETWEEN ? AND ?")
		args = append(args, startTimeFilter, endTimeFilter)
	}

	// TODO can add more filters

	query := "SELECT id, author, title, text, creation_time FROM blog"

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	// limit is not part of WHERE clause
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		query += " LIMIT ?"
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		args = append(args, limitInt)
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Failed to fetch blogs: " + err.Error(), http.StatusInternalServerError)
		return
	}

	var blogs []models.Blog
	for rows.Next() {
		var blog models.Blog
		err := rows.Scan(&blog.Id, &blog.Author, &blog.Title, &blog.Text, &blog.CreationTime)
		if err != nil {
			http.Error(w, "Failed to scan blog", http.StatusInternalServerError)
		}
		blogs = append(blogs, blog)
	}
	jsonRespose(w, blogs)
}

// deprecated
func GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, author, title, text, creation_time FROM blog")
	if err != nil {
		http.Error(w, "Failed to fetch blogs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var blog models.Blog
		err := rows.Scan(&blog.Id, &blog.Author, &blog.Title, &blog.Text, &blog.CreationTime)
		if err != nil {
			http.Error(w, "Failed to scan blog", http.StatusInternalServerError)
		}
		blogs = append(blogs, blog)
	}
	jsonRespose(w, blogs)
}

func CreateBlog(w http.ResponseWriter, r *http.Request) {
	var newBlog models.Blog
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &newBlog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newBlog.CreationTime = time.Now().Format("2006-01-02 15:04:05")

	statement, err := database.DB.Prepare("INSERT INTO blog (author, title, text, creation_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Failed to prepare insert statement", http.StatusInternalServerError)
		return
	}

	_, err = statement.Exec(newBlog.Author, newBlog.Title, newBlog.Text, newBlog.CreationTime)
	if err != nil {
		http.Error(w, "Failed to execute insert statement", http.StatusInternalServerError)
		return
	}

	jsonRespose(w, newBlog)
}

func UpdateBlog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing Id", http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	var updatedBlog models.Blog
	err = json.Unmarshal(reqBody, &updatedBlog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statement, err := database.DB.Prepare("UPDATE blog SET author = ?, title = ?, text = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "Failed to prepare update statement", http.StatusInternalServerError)
		return
	}

	_, err = statement.Exec(updatedBlog.Author, updatedBlog.Title, updatedBlog.Text, id)
	if err != nil {
		http.Error(w, "Failed to execute update statement", http.StatusInternalServerError)
		return
	}

	jsonRespose(w, updatedBlog)
}

func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing Id", http.StatusBadRequest)
		return
	}

	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM blog WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		http.Error(w, "Failed to check if entry exists", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Blog entry not found", http.StatusNotFound)
		return
	}

	statement, err := database.DB.Prepare("DELETE FROM blog WHERE id = ?")
	if err != nil {
		http.Error(w, "Failed to prepare delete statement", http.StatusInternalServerError)
		return
	}

	_, err = statement.Exec(id)
	if err != nil {
		http.Error(w, "Failed to execute delete statement", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Deleted blog with Id %s", id)
}
