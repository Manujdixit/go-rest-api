package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// This is a global variable that is used to store the items in the in-memory store of items.
var items []Item
var idCounter int = 1

func main() {
	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/items/", itemHandler)
	fmt.Println("server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	//get all items
	case http.MethodGet:
		//This line of code sets the Content-Type header of the HTTP response to application/json, indicating that the response body contains data in JSON format.
		w.Header().Set("Content-Type", "application/json")

		//This line of code encodes the items variable into JSON format and writes it to the HTTP response writer w.
		//In the context of main.go:itemsHandler, it sends the in-memory store of items as a JSON response to the client when a GET request is made.
		json.NewEncoder(w).Encode(items)

	//create new item
	case http.MethodPost:
		//This line of code declares a new variable newItem of type Item, which is a struct defined earlier in the code. The Item struct has two fields: ID and Name. This variable is used to hold a new item that will be created or updated in the code that follows.
		var newItem Item
		//to decode the request body into the newItem variable. It uses the json.NewDecoder function to create a new JSON decoder and then uses the Decode method to decode the request body into the newItem variable. If there is an error during the decoding process, it sends an HTTP error response with a status code of 400 (Bad Request) and the error message. The return statement is used to exit the function if there is an error.
		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newItem.ID = idCounter
		idCounter++
		items = append(items, newItem)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newItem)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	//This line of code extracts the ID from the URL path in the itemHandler function. It removes the prefix "/items/" from the URL path and assigns the remaining string to the idStr variable.
	idStr := r.URL.Path[len("/items/"):]
	//converts a string (idStr) to an integer (id) using the strconv.Atoi
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {

	//get item by id
	case http.MethodGet:

		for _, item := range items {
			if item.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(item)
				return
			}
		}
		http.NotFound(w, r)

	//update item by id
	case http.MethodPut:
		var updatedItem Item
		if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i, item := range items {
			if item.ID == id {
				items[i] = updatedItem
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(updatedItem)
				return
			}
		}
		http.NotFound(w, r)

	//delete item by id
	case http.MethodDelete:
		for i, item := range items {
			if item.ID == id {
				items = append(items[:i], items[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.NotFound(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
