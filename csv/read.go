package csv_operations

import (
	"encoding/csv"
	"os"
	"strconv"

	entities "github.com/gamze.sakallioglu/learningGo/homework-4-week-5-gamzesakallioglu/domain/entities"
	"gorm.io/gorm"
)

func ReadCsv(filename string) ([][]string, error) {

	// Open file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close() // Close the file at the end of the function

	csvReader := csv.NewReader(f)       // Create a csv reader
	records, err := csvReader.ReadAll() // Read everything inside the cvs file  // ReadAll returns a 2D slice line, element
	if err != nil {
		return nil, err
	}
	return records, nil
}

func ReadBooksCsv(filename string) (entities.Books, error) {

	records, err := ReadCsv(filename)
	if err != nil {
		return nil, err
	}

	var books entities.Books // Books -> []Book

	for _, line := range records[1:] { // Read line by line

		// Turn into int
		id, _ := strconv.Atoi(line[0])
		totalPage, _ := strconv.Atoi(line[2])
		totalStock, _ := strconv.Atoi(line[3])
		price, _ := strconv.Atoi(line[4])
		authorId, _ := strconv.Atoi(line[8])

		author := entities.Author{Model: gorm.Model{ID: uint(authorId)}, Name: line[7]}

		books = append(books, entities.Book{
			Model:      gorm.Model{ID: uint(id)},
			Name:       line[1],
			TotalPage:  totalPage,
			TotalStock: totalStock,
			Price:      float32(price),
			StockCode:  line[5],
			ISBN:       line[6],
			AuthorId:   authorId,
			Author:     author,
		})
	}
	return books, nil
}

func ReadAuthorsCsv(filename string) (entities.Authors, error) {

	records, err := ReadCsv(filename)
	if err != nil {
		return nil, err
	}

	var authors entities.Authors // Authors -> []Author

	for _, line := range records[1:] { // Read line by line

		// Turn into int
		id, _ := strconv.Atoi(line[0])

		authors = append(authors, entities.Author{
			Model: gorm.Model{ID: uint(id)},
			Name:  line[1],
		})
	}
	return authors, nil
}
