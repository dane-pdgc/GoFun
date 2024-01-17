package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/big"
	"math/rand"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Number struct to hold the random number
type Number struct {
	RandomNumber int `json:"randomNumber"`
}

// getRandomNumber handles the GET request and returns a random number
func getRandomNumber(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Generate a random number between 0 and 100
	randomNumber := rand.Intn(101) // No need to seed in Go 1.21, it's handled automatically

	// Create a Number struct with the random number
	number := Number{RandomNumber: randomNumber}

	// Encode the Number struct to JSON and send as response
	json.NewEncoder(w).Encode(number)
}

// User struct to hold user data from the database
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email_address"`
}

// getUsers handles the GET request and returns all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Data source name
	dsn := "root:toor@tcp(127.0.0.1:3306)/todo"

	// Open a new connection to the MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Query the database
	rows, err := db.Query("SELECT id, username, email_address FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Encode users to JSON and send as response
	json.NewEncoder(w).Encode(users)
}

// Request struct to hold the incoming request data
type FactorialRequest struct {
	Number int `json:"number"`
}

// Response struct to hold the factorial result
type FactorialResponse struct {
	Result string `json:"result"`
}

// calculateFactorial calculates the factorial of a given number
func calculateFactorial(n int) *big.Int {
	result := big.NewInt(1)
	for i := 1; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}

// factorialHandler handles the POST request and returns the factorial of the number
func factorialHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request
	var req FactorialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate the factorial
	result := calculateFactorial(req.Number)

	// Create a response struct with the result
	response := FactorialResponse{Result: result.String()}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the response struct to JSON and send as response
	json.NewEncoder(w).Encode(response)
}

func main() {

	// Register the URL path and handler function
	http.HandleFunc("/getRandomNumber", getRandomNumber)
	http.HandleFunc("/getUsers", getUsers)
	http.HandleFunc("/factorial", factorialHandler)

	// Start the server on port 8080 and handle requests
	http.ListenAndServe(":8080", nil)
}
