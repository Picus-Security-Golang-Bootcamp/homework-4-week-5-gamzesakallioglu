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

type AuthorHandler struct {
	authorRepo *repo.AuthorRepository
}

// ApiResponse struct for response body
type ApiResponse struct {
	Data interface{} `json:"data"`
}

func (a *AuthorHandler) NewAuthorHandler(authorRepo repo.AuthorRepository) {
	a.authorRepo = &authorRepo
}

func (a AuthorHandler) AuthorsHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		books, err := a.authorRepo.GetAuthorsWithBook()
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

func (a AuthorHandler) AuthorIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//r.URL.Query().Get("param")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		id, _ := strconv.Atoi(vars["id"])
		book, err := a.authorRepo.GetById(uint(id))
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

func (a AuthorHandler) AuthorCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var author entities.Author

		if r.Header.Get("Content-Type") != "application/json" {
			w.Write([]byte("Only json format is accepted"))
			return
		}

		err := json.NewDecoder(r.Body).Decode(&author)
		if err != nil {
			w.Write([]byte("400: Bad Request"))
			return
		}

		err = a.authorRepo.InsertOneData(author)
		if err != nil {
			w.Write([]byte("Data cannot inserted:("))
			return
		}

		bookData, err := json.Marshal(author)
		if err != nil {
			w.Write([]byte("400: Bad Request"))
			return
		}
		w.Write(bookData)
	}
}
