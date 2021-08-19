package api

import (
	"github.com/spf13/cast"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Request struct {
	Method             string
	URL                string
	Body               io.Reader
	ExpectedStatusCode int
}

func Test_createBook(t *testing.T) {
	requests := make([]Request, 2)
	requests[0] = Request{
		"POST",
		"http://localhost:8000/api/books",
		strings.NewReader(`{"Isbn": "0118263347", "Title": "The Alchemist", "Author": {"FirstName": "Paulo", "LastName": "Coelho"}}`),
		http.StatusCreated,
	}
	requests[1] = Request{
		"POST",
		"http://localhost:8000/api/books",
		strings.NewReader(`{"Isbn": "0178293357", "Title": "A Tale of Two Cities", "Author": {"FirstName": "Charles", "LastName": "Dickens"}}`),
		http.StatusCreated,
	}
	processRequest(t, requests)
}

func Test_getBooks(t *testing.T) {
	requests := make([]Request, 1)
	requests[0] = Request{
		"GET",
		"http://localhost:8000/api/books",
		nil,
		http.StatusOK,
	}
	processRequest(t, requests)
}

func Test_getBook(t *testing.T) {
	requests := make([]Request, 2)
	requests[0] = Request{
		"GET",
		"http://localhost:8000/api/books/1",
		nil,
		http.StatusOK,
	}
	requests[1] = Request{
		"GET",
		"http://localhost:8000/api/books/10",
		nil,
		http.StatusNotFound,
	}
	processRequest(t, requests)

}

func Test_updateBook(t *testing.T) {
	requests := make([]Request, 2)
	requests[0] = Request{
		"PUT",
		"http://localhost:8000/api/books/1",
		strings.NewReader(`{"Isbn": "Updated Isbn", "Title": "Updated Title", "Author": {"FirstName": "J. R. R.", "LastName": "Tolkien"}}`),
		http.StatusOK,
	}
	requests[1] = Request{
		"PUT",
		"http://localhost:8000/api/books/100",
		strings.NewReader(`{"Isbn": "Updated Isbn", "Title": "Updated Title", "Author": {"FirstName": "J. R. R.", "LastName": "Tolkien"}}`),
		http.StatusNotFound,
	}
	processRequest(t, requests)
}

func Test_deleteBook(t *testing.T) {
	requests := make([]Request, 2)
	requests[0] = Request{
		"DELETE",
		"http://localhost:8000/api/books/1",
		nil,
		http.StatusNoContent,
	}
	requests[1] = Request{
		"DELETE",
		"http://localhost:8000/api/books/2",
		nil,
		http.StatusNoContent,
	}
	processRequest(t, requests)
}

func processRequest(t *testing.T, requests []Request) {
	for _, req := range requests {
		r, _ := http.NewRequest(req.Method, req.URL, req.Body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		if w.Code != req.ExpectedStatusCode {
			t.Error("Expected status code: " + cast.ToString(req.ExpectedStatusCode) + "\tFound: " + cast.ToString(w.Code) + "\n")
		}
	}
}
