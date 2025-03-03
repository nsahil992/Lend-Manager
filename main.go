// main.go
package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Friend represents a friend in the system
type Friend struct {
	ID   int
	Name string
}

// Item represents an item lent to a friend
type Item struct {
	ID       int
	Name     string
	FriendID int
}

func main() {
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

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Check connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
	fmt.Println("Successfully connected to database!")

	// Create tables if they don't exist
	createTables(db)

	reader := bufio.NewReader(os.Stdin)

	// Main program loop
	for {
		fmt.Println("\nWhat do you want to do? (takeback/give/newfriend/quit)")
		userAction, _ := reader.ReadString('\n')
		userAction = strings.TrimSpace(userAction)

		switch userAction {
		case "quit":
			fmt.Println("Goodbye!")
			return

		case "takeback":
			handleTakeback(db, reader)

		case "give":
			handleGive(db, reader)

		case "newfriend":
			handleNewFriend(db, reader)

		default:
			fmt.Println("Sorry, I didn't understand that. (Valid choices: give/takeback/newfriend/quit)")
		}
	}
}

func createTables(db *sql.DB) {
	// Create friends table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS friends (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		log.Fatal("Failed to create friends table: ", err)
	}

	// Create items table
	_, err = db.Exec(`
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

func handleTakeback(db *sql.DB, reader *bufio.Reader) {
	// Get all friends
	rows, err := db.Query("SELECT id, name FROM friends ORDER BY name")
	if err != nil {
		fmt.Println("Error fetching friends: ", err)
		return
	}
	defer rows.Close()

	// Display friends
	fmt.Println("These are your friends:")
	var friends []Friend
	for rows.Next() {
		var friend Friend
		err := rows.Scan(&friend.ID, &friend.Name)
		if err != nil {
			fmt.Println("Error scanning friend: ", err)
			return
		}
		friends = append(friends, friend)
		fmt.Println(friend.Name)
	}

	if len(friends) == 0 {
		fmt.Println("You don't have any friends in the system yet. Add a friend first.")
		return
	}

	// Ask which friend
	fmt.Print("Which friend did you lend to? ")
	friendName, _ := reader.ReadString('\n')
	friendName = strings.TrimSpace(friendName)

	// Find friend
	var friendID int
	err = db.QueryRow("SELECT id FROM friends WHERE name = $1", friendName).Scan(&friendID)
	if err != nil {
		fmt.Println("Sorry, I didn't find that friend.")
		return
	}

	// Get items lent to this friend
	rows, err = db.Query("SELECT id, name FROM items WHERE friend_id = $1", friendID)
	if err != nil {
		fmt.Println("Error fetching items: ", err)
		return
	}
	defer rows.Close()

	// Display items
	var items []Item
	fmt.Printf("This is what you gave to %s:\n", friendName)
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			fmt.Println("Error scanning item: ", err)
			return
		}
		items = append(items, item)
		fmt.Println(item.Name)
	}

	if len(items) == 0 {
		fmt.Printf("You haven't given anything to %s\n", friendName)
		return
	}

	// Ask which item to take back
	fmt.Printf("What did you take back from %s? ", friendName)
	itemName, _ := reader.ReadString('\n')
	itemName = strings.TrimSpace(itemName)

	// Find item
	var itemID int
	err = db.QueryRow("SELECT id FROM items WHERE name = $1 AND friend_id = $2", itemName, friendID).Scan(&itemID)
	if err != nil {
		fmt.Println("Sorry, I didn't find that item.")
		return
	}

	// Delete item
	_, err = db.Exec("DELETE FROM items WHERE id = $1", itemID)
	if err != nil {
		fmt.Println("Error deleting item: ", err)
		return
	}

	fmt.Printf("Alright, I'll remember that you took %s from %s\n", itemName, friendName)
}

func handleGive(db *sql.DB, reader *bufio.Reader) {
	// Get all friends
	rows, err := db.Query("SELECT id, name FROM friends ORDER BY name")
	if err != nil {
		fmt.Println("Error fetching friends: ", err)
		return
	}
	defer rows.Close()

	// Display friends
	fmt.Println("These are your friends:")
	var friends []Friend
	for rows.Next() {
		var friend Friend
		err := rows.Scan(&friend.ID, &friend.Name)
		if err != nil {
			fmt.Println("Error scanning friend: ", err)
			return
		}
		friends = append(friends, friend)
		fmt.Println(friend.Name)
	}

	if len(friends) == 0 {
		fmt.Println("You don't have any friends in the system yet. Add a friend first.")
		return
	}

	// Ask which friend
	fmt.Print("Which friend did you lend to? ")
	friendName, _ := reader.ReadString('\n')
	friendName = strings.TrimSpace(friendName)

	// Find friend
	var friendID int
	err = db.QueryRow("SELECT id FROM friends WHERE name = $1", friendName).Scan(&friendID)
	if err != nil {
		fmt.Println("Sorry, I didn't find that friend.")
		return
	}

	// Ask which item
	fmt.Printf("What did you lend to %s? ", friendName)
	itemName, _ := reader.ReadString('\n')
	itemName = strings.TrimSpace(itemName)

	// Insert item
	_, err = db.Exec("INSERT INTO items (name, friend_id) VALUES ($1, $2)", itemName, friendID)
	if err != nil {
		fmt.Println("Error adding item: ", err)
		return
	}

	fmt.Printf("Got it! You lent %s to %s.\n", itemName, friendName)
}

func handleNewFriend(db *sql.DB, reader *bufio.Reader) {
	// Ask friend name
	fmt.Print("Who is your new friend? ")
	friendName, _ := reader.ReadString('\n')
	friendName = strings.TrimSpace(friendName)

	// Insert friend
	_, err := db.Exec("INSERT INTO friends (name) VALUES ($1)", friendName)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			fmt.Println("That friend already exists.")
		} else {
			fmt.Println("Error adding friend: ", err)
		}
		return
	}

	fmt.Printf("Great! I've added %s as your friend.\n", friendName)
}
