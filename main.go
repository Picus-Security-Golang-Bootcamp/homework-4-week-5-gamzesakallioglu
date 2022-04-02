package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	postgres "github.com/gamze.sakallioglu/learningGo/homework-4-week-5-gamzesakallioglu/common/db"
	csv_operations "github.com/gamze.sakallioglu/learningGo/homework-4-week-5-gamzesakallioglu/csv"
	custom_handlers "github.com/gamze.sakallioglu/learningGo/homework-4-week-5-gamzesakallioglu/domain/handlers"
	repo "github.com/gamze.sakallioglu/learningGo/homework-4-week-5-gamzesakallioglu/domain/repositories"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize the app with a router and gorm db. Also added logging middleware and authentication middleware
func (a *App) Initialize() error {

	err := godotenv.Load()
	if err != nil {
		return err
	}

	db, err := postgres.NewPsqlDB()
	if err != nil {
		return err
	}

	// Initialize DB and Router
	a.DB = db
	a.Router = mux.NewRouter()

	handlers.AllowedOrigins([]string{"https://www.example.com"})
	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	handlers.AllowedMethods([]string{"POST", "GET", "PUT", "PATCH"})

	a.Router.Use(loggingMiddleware)
	//a.Router.Use(authenticationMiddleware)

	amw := AuthenticationMiddleware{make(map[string]string)}
	amw.Populate()
	a.Router.Use(amw.Middleware)

	return nil

}

// Run the app with given port
func (a *App) Run(addr string) {

	server := &http.Server{
		Addr:         "localhost:" + addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ShutdownServer(server, time.Second*10)

}

// ApiResponse struct for response body
type ApiResponse struct {
	Data interface{} `json:"data"`
}

// AuthenticationMiddleware struct for declearing the token and users
type AuthenticationMiddleware struct {
	tokenUsers map[string]string
}

func (amw *AuthenticationMiddleware) Populate() {
	amw.tokenUsers["00000000"] = "user0"
	amw.tokenUsers["aaaaaaaa"] = "userA"
	amw.tokenUsers["05f717e5"] = "randomUser"
	amw.tokenUsers["deadbeef"] = "user0"
}

// This is middleware function. In order to call this for every request, calling this in App's Initialize method via Use()
func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func main() {

	// Read csv Files
	file_book := "books.csv"
	file_author := "authors.csv"
	bookList, err := csv_operations.ReadBooksCsv(file_book)
	if err != nil {
		fmt.Println(err)
	}

	authorList, err := csv_operations.ReadAuthorsCsv(file_author)
	if err != nil {
		fmt.Println(err)
	}

	// Initialize the App
	app := App{}
	app.Initialize()

	// Get the db from App
	db := app.DB
	// Book and Author repo. And inserting the datas inside of the csv files
	bookRepo := repo.NewBookRepository(db)
	bookRepo.Migrations()
	bookRepo.InsertDatas(bookList)
	authorRepo := repo.NewAuthorRepository(db)
	authorRepo.Migrations()
	authorRepo.InsertDatas(authorList)

	author_handler := custom_handlers.AuthorHandler{}
	author_handler.NewAuthorHandler(*authorRepo)

	book_handler := custom_handlers.BookHandler{}
	book_handler.NewBookHandler(*bookRepo)

	// API - Books

	// Get - for list all the books with author
	app.Router.HandleFunc("/books", book_handler.BooksHandler())

	// Get - for search by the name
	app.Router.PathPrefix("/books").Subrouter().HandleFunc("/{name}", book_handler.BookNameHandler())

	// Get - for search by id
	app.Router.PathPrefix("/books").Subrouter().HandleFunc("/{id}", book_handler.BookIdHandler())

	// Patch - for buy given amount of book with given id
	app.Router.PathPrefix("/books/buy").Subrouter().HandleFunc("/{id}/{quantitiy}", book_handler.BookBuyHandler()).Methods(http.MethodPatch)

	// Post - for create a new book
	app.Router.PathPrefix("/books").Subrouter().HandleFunc("/", book_handler.BookCreateHandler()).Methods(http.MethodPost)
	// Books

	// API - Authors

	// Get - for list all the authors with books they published
	app.Router.HandleFunc("/authors", author_handler.AuthorsHandler())

	// Get - for search by id
	app.Router.PathPrefix("/authors").Subrouter().HandleFunc("/{id}", author_handler.AuthorIdHandler())

	// Post - for create a new author
	app.Router.PathPrefix("/authors").Subrouter().HandleFunc("/", author_handler.AuthorCreateHandler()).Methods(http.MethodPost)
	// Authors

	// Run the app
	app.Run("8090")

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func ShutdownServer(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
