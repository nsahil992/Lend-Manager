// api.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Friend represents a friend in the system
type Friend struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Item represents an item lent to a friend
type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FriendID int    `json:"friendId"`
}

// App represents the application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Get database credentials from environment variables
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	a.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Check connection
	err = a.DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
	log.Println("Successfully connected to database!")

	// Create tables if they don't exist
	a.createTables()

	// Initialize router
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run starts the app on the specified address
func (a *App) Run(addr string) {
	log.Printf("Server started on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, a.enableCORS(a.Router)))
}

// enableCORS enables CORS for the router
func (a *App) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// createTables creates the necessary tables
func (a *App) createTables() {
	// Create friends table
	_, err := a.DB.Exec(`
		CREATE TABLE IF NOT EXISTS friends (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		log.Fatal("Failed to create friends table: ", err)
	}

	// Create items table
	_, err = a.DB.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			friend_id INTEGER REFERENCES friends(id)
		)
	`)
	if err != nil {
		log.Fatal("Failed to create items table: ", err)
	}
}

// initializeRoutes initializes the routes
func (a *App) initializeRoutes() {
	// Define API routes
	api := a.Router.PathPrefix("/api").Subrouter()

	// Friend routes
	api.HandleFunc("/friends", a.getFriends).Methods("GET")
	api.HandleFunc("/friends", a.createFriend).Methods("POST")
	api.HandleFunc("/friends/{id}", a.getFriend).Methods("GET")
	api.HandleFunc("/friends/{id}", a.deleteFriend).Methods("DELETE")

	// Item routes
	api.HandleFunc("/friends/{id}/items", a.getFriendItems).Methods("GET")
	api.HandleFunc("/items", a.createItem).Methods("POST")
	api.HandleFunc("/items/{id}", a.getItem).Methods("GET")
	api.HandleFunc("/items/{id}", a.deleteItem).Methods("DELETE")

	// Serve static files
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
}

// respondWithError responds with an error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON responds with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// getFriends gets all friends
func (a *App) getFriends(w http.ResponseWriter, r *http.Request) {
	// Get all friends from database
	rows, err := a.DB.Query("SELECT id, name FROM friends ORDER BY name")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving friends")
		return
	}
	defer rows.Close()

	// Create friends slice
	friends := []Friend{}
	for rows.Next() {
		var friend Friend
		if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning friend")
			return
		}
		friends = append(friends, friend)
	}

	// Respond with friends
	respondWithJSON(w, http.StatusOK, friends)
}

// getFriend gets a specific friend
func (a *App) getFriend(w http.ResponseWriter, r *http.Request) {
	// Get friend ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid friend ID")
		return
	}

	// Get friend from database
	var friend Friend
	err = a.DB.QueryRow("SELECT id, name FROM friends WHERE id = $1", id).Scan(&friend.ID, &friend.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Friend not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving friend")
		}
		return
	}

	// Respond with friend
	respondWithJSON(w, http.StatusOK, friend)
}

// createFriend creates a new friend
func (a *App) createFriend(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var friend Friend
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&friend); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate friend name
	if friend.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Friend name is required")
		return
	}

	// Insert friend into database
	err := a.DB.QueryRow("INSERT INTO friends (name) VALUES ($1) RETURNING id", friend.Name).Scan(&friend.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating friend")
		return
	}

	// Respond with created friend
	respondWithJSON(w, http.StatusCreated, friend)
}

// deleteFriend deletes a friend
func (a *App) deleteFriend(w http.ResponseWriter, r *http.Request) {
	// Get friend ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid friend ID")
		return
	}

	// Delete friend from database
	_, err = a.DB.Exec("DELETE FROM friends WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting friend")
		return
	}

	// Respond with success
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// getFriendItems gets all items for a specific friend
func (a *App) getFriendItems(w http.ResponseWriter, r *http.Request) {
	// Get friend ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid friend ID")
		return
	}

	// Get items from database
	rows, err := a.DB.Query("SELECT id, name, friend_id FROM items WHERE friend_id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving items")
		return
	}
	defer rows.Close()

	// Create items slice
	items := []Item{}
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.FriendID); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning item")
			return
		}
		items = append(items, item)
	}

	// Respond with items
	respondWithJSON(w, http.StatusOK, items)
}

// getItem gets a specific item
func (a *App) getItem(w http.ResponseWriter, r *http.Request) {
	// Get item ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	// Get item from database
	var item Item
	err = a.DB.QueryRow("SELECT id, name, friend_id FROM items WHERE id = $1", id).Scan(&item.ID, &item.Name, &item.FriendID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Item not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving item")
		}
		return
	}

	// Respond with item
	respondWithJSON(w, http.StatusOK, item)
}

// createItem creates a new item
func (a *App) createItem(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var item Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate item
	if item.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Item name is required")
		return
	}

	if item.FriendID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Friend ID is required")
		return
	}

	// Check if friend exists
	var friendExists bool
	err := a.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM friends WHERE id = $1)", item.FriendID).Scan(&friendExists)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking friend existence")
		return
	}

	if !friendExists {
		respondWithError(w, http.StatusBadRequest, "Friend does not exist")
		return
	}

	// Insert item into database
	err = a.DB.QueryRow("INSERT INTO items (name, friend_id) VALUES ($1, $2) RETURNING id", item.Name, item.FriendID).Scan(&item.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating item")
		return
	}

	// Respond with created item
	respondWithJSON(w, http.StatusCreated, item)
}

// deleteItem deletes an item
func (a *App) deleteItem(w http.ResponseWriter, r *http.Request) {
	// Get item ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	// Delete item from database
	_, err = a.DB.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting item")
		return
	}

	// Respond with success
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func main() {
	app := App{}
	app.Initialize()
	app.Run(":8081")
}
