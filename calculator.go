package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort" // Import the sort package
	"strconv"

	"github.com/gorilla/mux"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

var products = []Product{
	{ID: 1, Name: "Product 1", Category: "Electronics", Price: 500},
	{ID: 2, Name: "Product 2", Category: "Clothing", Price: 50},
	{ID: 3, Name: "Product 3", Category: "Electronics", Price: 800},
	// Add more products...
}

func getTopNProductsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryName := vars["categoryname"]
	nStr := r.URL.Query().Get("n")

	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, "Invalid 'n' parameter", http.StatusBadRequest)
		return
	}

	if n <= 0 {
		http.Error(w, "'n' parameter must be greater than zero", http.StatusBadRequest)
		return
	}

	var categoryProducts []Product
	for _, p := range products {
		if p.Category == categoryName {
			categoryProducts = append(categoryProducts, p)
		}
	}

	// Check if there are products in the specified category
	if len(categoryProducts) == 0 {
		http.Error(w, "No products found in the specified category", http.StatusNotFound)
		return
	}

	// Sort products by price
	sort.Slice(categoryProducts, func(i, j int) bool {
		return categoryProducts[i].Price > categoryProducts[j].Price
	})

	if n > len(categoryProducts) {
		n = len(categoryProducts) // Limit 'n' to the number of products available
	}

	topNProducts := categoryProducts[:n]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topNProducts)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/categories/{categoryname}/products", getTopNProductsByCategory).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
