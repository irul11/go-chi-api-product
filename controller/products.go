package controller

import (
	"encoding/json"
	"fmt"
	"learn-go-chi/helper"
	"learn-go-chi/models"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

var Channel chan models.Message

func Exec(wg *sync.WaitGroup, i int) {
	defer wg.Done()
	for msg := range Channel {
		fmt.Println("worker", i, "is working")
		var output models.MessageOutput

		if msg.Method == "get" {
			output = helper.Get(msg)
		} else if msg.Method == "post" {
			output = helper.Post(msg)
		} else if msg.Method == "put" {
			output = helper.Put(msg)
		} else if msg.Method == "delete" {
			output = helper.Del(msg)
		}

		// For timing testing
		// random := 3*time.Second + time.Duration(rand.Intn(5)*int(time.Second))
		// fmt.Println("Time duration is:", random)
		// time.Sleep(random)

		msg.OutputChannel <- output

		fmt.Println("worker", i, "is Done")
	}
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, price, quantity, created_at FROM products`

	outputChannel := make(chan models.MessageOutput)
	Channel <- models.Message{
		Query:         query,
		Method:        "get",
		OutputChannel: outputChannel,
	}
	defer close(outputChannel)

	products := <-outputChannel
	if products.Error != nil {
		log.Printf("Error querying product: %v", products.Error)
		http.Error(w, "Error querying product data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(products.Data)
	if err != nil {
		log.Printf("Error marshalling product to JSON: %v", err)
		http.Error(w, "Error processing product data", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
	// Another method to send response json
	// err = json.NewEncoder(w).Encode(products)
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	productIdParam := chi.URLParam(r, "productId")
	productId, err := strconv.Atoi(productIdParam)

	if err != nil {
		// Handle invalid ID error
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, price, quantity, created_at FROM products WHERE id=$1`

	outputChannel := make(chan models.MessageOutput)
	Channel <- models.Message{
		Query:         query,
		Method:        "get",
		URLParam:      productId,
		OutputChannel: outputChannel,
	}
	defer close(outputChannel)

	output := <-outputChannel
	err = output.Error
	products := output.Data

	if err != nil {
		log.Printf("Error querying product: %v", err)
		http.Error(w, "Error querying product data", http.StatusInternalServerError)
		return
	}

	if len(products) < 1 {
		log.Println("Product not found")
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(products[0])
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO products (name, price, quantity, description)
				VALUES ($1, $2, $3, $4)
	`

	outputChannel := make(chan models.MessageOutput)
	Channel <- models.Message{
		Query:         query,
		Method:        "post",
		Body:          product,
		OutputChannel: outputChannel,
	}
	defer close(outputChannel)

	output := <-outputChannel
	err = output.Error

	if err != nil {
		// Handle error
		http.Error(w, "Failed insert data", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode("Insert data success")
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Parsing id from URL
	productIdParam := chi.URLParam(r, "productId")
	productId, err := strconv.Atoi(productIdParam)

	if err != nil {
		// Handle invalid ID error
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parsing request body
	var product models.Product
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query
	query := `UPDATE products SET name=$1, price=$2, quantity=$3, description=$4 WHERE id=$5`

	outputChannel := make(chan models.MessageOutput)
	Channel <- models.Message{
		Query:         query,
		Method:        "put",
		Body:          product,
		URLParam:      productId,
		OutputChannel: outputChannel,
	}
	defer close(outputChannel)

	output := <-outputChannel
	err = output.Error
	if err != nil {
		// Handle user not found error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode("Data with id:" + strconv.Itoa(productId) + " succesfully updated")
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productIdParam := chi.URLParam(r, "productId")
	productId, err := strconv.Atoi(productIdParam)

	if err != nil {
		// Handle invalid ID error
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM products WHERE id=$1`

	outputChannel := make(chan models.MessageOutput)
	Channel <- models.Message{
		Query:         query,
		Method:        "delete",
		URLParam:      productId,
		OutputChannel: outputChannel,
	}
	defer close(outputChannel)

	output := <-outputChannel
	err = output.Error
	if err != nil {
		// Handle user not found error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode("Data with id: " + strconv.Itoa(productId) + " successfully deleted")
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}
