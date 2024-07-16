package helper

import (
	"database/sql"
	"errors"
	"learn-go-chi/database"
	"learn-go-chi/models"
	"log"
)

func Get(msg models.Message) models.MessageOutput {
	var (
		rows *sql.Rows
		err  error
	)

	if msg.URLParam != 0 {
		rows, err = database.DB.Query(msg.Query, msg.URLParam)
	} else {
		rows, err = database.DB.Query(msg.Query)
	}

	if err != nil {
		log.Printf("Error querying product: %v", err)
		return models.MessageOutput{
			Data:  nil,
			Error: err,
		}
	}

	products := []models.Product{}

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Created_at)
		if err != nil {
			log.Printf("Error scanning product: %v", err)
			return models.MessageOutput{
				Data:  nil,
				Error: err,
			}
		}
		products = append(products, product)
	}
	rows.Close()

	return models.MessageOutput{
		Data:  products,
		Error: nil,
	}
}

func Post(msg models.Message) models.MessageOutput {
	product := msg.Body

	_, err := database.DB.Exec(msg.Query, product.Name, product.Price, product.Quantity, product.Description)

	if err != nil {
		log.Printf("Error creating product: %v", err)
		return models.MessageOutput{
			Error: err,
		}
	}

	return models.MessageOutput{}
}

func Put(msg models.Message) models.MessageOutput {
	product := msg.Body

	result, err := database.DB.Exec(msg.Query, product.Name, product.Price, product.Quantity, product.Description, msg.URLParam)

	if err != nil {
		log.Printf("Error updating product: %v", err)
		return models.MessageOutput{
			Error: err,
		}
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		log.Printf("Error updating product: %v", err)
		return models.MessageOutput{
			Error: err,
		}
	}
	if rowAffected == 0 {
		log.Println("Error: no row affected")
		return models.MessageOutput{
			Error: errors.New("no row affected"),
		}
	}

	return models.MessageOutput{}
}

func Del(msg models.Message) models.MessageOutput {
	result, err := database.DB.Exec(msg.Query, msg.URLParam)

	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return models.MessageOutput{
			Error: err,
		}
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		log.Printf("Error updating product: %v", err)
		return models.MessageOutput{
			Error: err,
		}
	}
	if rowAffected == 0 {
		log.Println("Error: no row affected")
		return models.MessageOutput{
			Error: errors.New("no row affected"),
		}
	}

	return models.MessageOutput{}
}
