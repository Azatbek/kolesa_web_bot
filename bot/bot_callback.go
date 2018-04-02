package bot

import (
	"../db"
	"fmt"
)

func getCategories() []db.Categories {
	categories, err := db.GetCategories()

	if err != nil {
		fmt.Println(err)
	}

	return categories
}