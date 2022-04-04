package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name       string `gorm:"varchar(150)"`
	TotalPage  int
	TotalStock int
	Price      float32
	StockCode  string `gorm:"varchar(50)"`
	ISBN       string `gorm:"varchar(50)"`
	Author     Author `gorm:"foreignKey:AuthorId;references:ID"`
	AuthorId   int
}

type Books []Book

// Adding Hooks

// overwriting ToString. Books and the author
func (b *Book) ToString() string {
	returnStatement := fmt.Sprintf("ID: %v\nNAME: %s\nPAGE NUMBER: %v\nTOTAL STOCK: %v\nPRICE: %v\nSTOCK CODE: %s\nISBN: %s\n",
		b.ID, b.Name, b.TotalPage, b.TotalStock, b.Price, b.StockCode, b.ISBN)

	if b.AuthorId > 0 {
		returnStatement += fmt.Sprintf("AUTHOR INFO: \n\tID: %v \n\tNAME: %s", b.Author.ID, b.Author.Name)
	}
	return returnStatement
}

func (b *Book) BeforeDelete(db *gorm.DB) (err error) {
	fmt.Printf("%s is being deleted", b.Name)
	return nil
}
