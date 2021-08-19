package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/hmsayem/go-rest-api/authentication"
	"log"
	"net/http"
	"strconv"
)

type JwtToken struct {
	Token string `json:"token"`
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Book struct {
	ID     int     `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var books []Book
var jwtKey = []byte("secret_key")
var router = mux.NewRouter()

// Sign In
func login(w http.ResponseWriter, r *http.Request) {

	tokenString, err := authentication.GenerateJWT("admin", w, r)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenString,
		//Expires: expirationTime,
	})

	err = json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	if err != nil {
		log.Println(err)
	}
}

func getID() int {
	if len(books) == 0 {
		return 0
	}
	return books[len(books)-1].ID + 1
}

//Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}

// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, item := range books {
		if item.ID == id {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
	}
	book.ID = getID()
	books = append(books, book)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		log.Println(err)
	}
}

// Update a book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
	}
	for index, item := range books {
		if item.ID == id {
			book.ID = books[index].ID
			books[index] = book
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(book)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

// Delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for index, item := range books {
		if item.ID == id {
			w.WriteHeader(http.StatusNoContent)
			books = append(books[:index], books[index+1:]...)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func createDB() {
	//Mock Data
	books = append(books,
		Book{
			ID:     1,
			Isbn:   "0618260307",
			Title:  "The Hobbit",
			Author: &Author{FirstName: "J. R. R.", LastName: "Tolkien"},
		},
		Book{
			ID:     2,
			Isbn:   "0908606664",
			Title:  "Slinky Malinki",
			Author: &Author{FirstName: "Lynley", LastName: "Dodd"},
		},
		Book{
			ID:     3,
			Isbn:   "0908783116",
			Title:  "Mechanical Harry",
			Author: &Author{FirstName: "Bob", LastName: "Kerr"},
		},
	)
}
func handleRequest() {
	//Endpoints
	router.HandleFunc("/api/login", authentication.BasicAuthentication(login)).Methods("GET")
	router.HandleFunc("/api/books", authentication.JWTAuthentication(getBooks)).Methods("GET")
	router.HandleFunc("/api/books/{id}", authentication.JWTAuthentication(getBook)).Methods("GET")
	router.HandleFunc("/api/books", authentication.JWTAuthentication(createBook)).Methods("POST")
	router.HandleFunc("/api/books/{id}", authentication.JWTAuthentication(updateBook)).Methods("PUT")
	router.HandleFunc("/api/books/{id}", authentication.JWTAuthentication(deleteBook)).Methods("DELETE")
}
func Run(p string) {
	port := ":" + p
	log.Fatal(http.ListenAndServe(port, router))
}
func init() {
	createDB()
	handleRequest()
}
