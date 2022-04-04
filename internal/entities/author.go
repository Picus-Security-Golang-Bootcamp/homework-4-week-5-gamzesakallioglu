package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name  string `gorm:"varchar(150)" json:"name"`
	Books []Book
}

type Authors []Author

// Adding Hooks

// overwriting ToString. Author and the books they published
func (a *Author) ToString() string {
	returnStatement := fmt.Sprintf("ID: %v\nAUTHOR NAME: %s\n", a.ID, a.Name)
	for i, v := range a.Books {
		returnStatement += fmt.Sprintf("BOOK #%v: \n\tID: %v\n\tNAME: %s\n\tPAGE NUMBER: %v\n\tTOTAL STOCK: %v\n\tPRICE: %v\n\tSTOCK CODE: %s\n\tISBN: %s\n",
			i+1, v.ID, v.Name, v.TotalPage, v.TotalStock, v.Price, v.StockCode, v.ISBN)
	}
	return returnStatement
}

func (a *Author) BeforeDelete(db *gorm.DB) (err error) {
	fmt.Printf("%s is being deleted", a.Name)
	return nil
}
