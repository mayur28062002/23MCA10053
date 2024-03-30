package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type NumbersRequest struct {
	Qualifiers []string `json:"qualifiers"`
	WindowSize int      `json:"window_size"`
}

type NumbersResponse struct {
	Numbers []int `json:"numbers"`
}

func fetchNumbersFromServer(url string) ([]int, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var data map[string][]int
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}

	primeNumbers, ok := data["numbers"]
	if !ok {
		return nil, errors.New("prime numbers not found in response")
	}

	return primeNumbers, nil
}

func applyQualifiers(numbers []int, qualifiers []string) []int {
	filteredNumbers := make([]int, 0)
	for _, num := range numbers {
		for _, qualifier := range qualifiers {
			switch qualifier {
			case "p": // Prime number qualifier
				if isPrime(num) {
					filteredNumbers = append(filteredNumbers, num)
				}
			case "f": // Fibonacci number qualifier
				if isFibonacci(num) {
					filteredNumbers = append(filteredNumbers, num)
				}
			case "e": // Even number qualifier
				if num%2 == 0 {
					filteredNumbers = append(filteredNumbers, num)
				}
			case "r": // Random number qualifier
				filteredNumbers = append(filteredNumbers, num)
			}
		}
	}
	return filteredNumbers
}

func isPrime(num int) bool {
	if num <= 1 {
		return false
	}
	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func isFibonacci(num int) bool {
	a, b := 0, 1
	for b <= num {
		if b == num {
			return true
		}
		a, b = b, a+b
	}
	return false
}

func calculateAverageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	numberID := vars["numberid"]

	var req NumbersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Construct the URL based on the number ID provided
	url := "http://20.244.56.144/test/" + numberID

	// Fetch numbers from the server
	numbers, err := fetchNumbersFromServer(url)
	if err != nil {
		http.Error(w, "Failed to fetch numbers from server", http.StatusInternalServerError)
		return
	}

	// Apply qualifiers and window size to the fetched numbers
	filteredNumbers := applyQualifiers(numbers, req.Qualifiers)

	// Take only unique numbers
	uniqueNumbers := make([]int, 0)
	uniqueMap := make(map[int]bool)
	for _, num := range filteredNumbers {
		if !uniqueMap[num] {
			uniqueMap[num] = true
			uniqueNumbers = append(uniqueNumbers, num)
		}
	}

	// Apply window size if specified
	if req.WindowSize > 0 && req.WindowSize < len(uniqueNumbers) {
		uniqueNumbers = uniqueNumbers[:req.WindowSize]
	}

	resp := NumbersResponse{Numbers: uniqueNumbers}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/numbers/{numberid}", calculateAverageHandler).Methods("POST")

	log.Println("Starting HTTP server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

