package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gamze.sakallioglu/learningGo/homework-4-week-5-gamzesakallioglu/domain/entities"
	repo "github.com/gamze.sakallioglu/learningGo/homework-4-week-5-gamzesakallioglu/domain/repositories"
	"github.com/gorilla/mux"
)

type BookHandler struct {
	bookRepo *repo.BookRepository
}

func (b *BookHandler) NewBookHandler(bookRepo repo.BookRepository) {
	b.bookRepo = &bookRepo
}

func (b *BookHandler) BooksHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		books, err := b.bookRepo.GetBooksWithAuthor()
		if err != nil {
			log.Fatal(err)
		}

		d := ApiResponse{
			Data: books,
		}

		resp, _ := json.Marshal(d)
		w.Write(resp)
	}
}

func (b *BookHandler) BookCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book entities.Book

		if r.Header.Get("Content-Type") != "application/json" {
			w.Write([]byte("Only json format is accepted"))
			return
		}

		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			w.Write([]byte("400: Bad Request"))
			return
		}

		err = b.bookRepo.InsertOneData(book)
		if err != nil {
			w.Write([]byte("Data cannot inserted:("))
			return
		}

		bookData, err := json.Marshal(book)
		if err != nil {
			w.Write([]byte("400: Bad Request"))
			return
		}
		w.Write(bookData)
	}
}

func (b *BookHandler) BookIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//r.URL.Query().Get("param")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		id, _ := strconv.Atoi(vars["id"])
		book, err := b.bookRepo.GetById(uint(id))
		if err != nil {
			log.Fatal(err)
		}
		d := ApiResponse{
			Data: book,
		}

		resp, _ := json.Marshal(d)
		w.Write(resp)
	}
}

func (b *BookHandler) BookBuyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//r.URL.Query().Get("param")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		id, _ := strconv.Atoi(vars["id"])
		quantity, _ := strconv.Atoi(vars["quantitiy"])
		message, err := b.bookRepo.BuyBook(uint(id), quantity)
		if err != nil {
			log.Fatal(err)
		}
		d := ApiResponse{
			Data: message,
		}

		resp, _ := json.Marshal(d)
		w.Write(resp)
	}
}

func (b *BookHandler) BookNameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//r.URL.Query().Get("param")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		books, err := b.bookRepo.FindByName(vars["name"])
		if err != nil {
			log.Fatal(err)
		}
		d := ApiResponse{
			Data: books,
		}

		resp, _ := json.Marshal(d)
		w.Write(resp)
	}

}
