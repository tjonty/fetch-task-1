package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Store receipt IDs and points
var receipts = make(map[string]int)
var pointsData = make(map[string]int)

// Receipt represents the structure of the JSON receipt data
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Item represents the structure of an item in the receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}
type PointsResponse struct {
	Points int `json:"points"`
}

func generateIdAndPoints(receipt Receipt) (string, int, error) {
	points := 0

	// Rule: One point for every alphanumeric character in the retailer name
	alphanumericRegex := regexp.MustCompile("[a-zA-Z0-9]")
	retailerName := receipt.Retailer
	alphanumericCount := len(alphanumericRegex.FindAllString(retailerName, -1))
	points += alphanumericCount

	// Rule: 50 points if the total is a round dollar amount with no cents
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// Rule: 25 points if the total is a multiple of 0.25
	totalValue := parseFloat(receipt.Total)
	if math.Mod(totalValue, 0.25) == 0 {
		points += 25
	}

	// Rule: 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// Rule: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			itemPrice := parseFloat(item.Price)
			points += int(math.Ceil(itemPrice * 0.2))
		}
	}

	// Rule: 6 points if the day in the purchase date is odd
	// Rule: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return "", points, fmt.Errorf("Invalid purchase date")
	} else {
		if purchaseDate.Day()%2 != 0 {
			points += 6
		}

		purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
		if err != nil {
			return "", points, fmt.Errorf("Invalid purchase time")
		} else {
			if purchaseTime.After(time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)) &&
				purchaseTime.Before(time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)) {
				points += 10
			}
		}
	}

	// Use a UUID library to generate IDs
	id := uuid.New().String()
	pointsData[id] = points

	return id, points, nil
}

// parseFloat is a helper function to convert string to float64
func parseFloat(s string) float64 {
	var value float64
	_, _ = fmt.Sscanf(s, "%f", &value)
	return value
}

// ---------- router ----------

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")
	return router
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {

	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, _, err := generateIdAndPoints(receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	points, found := pointsData[id]
	if !found {
		http.NotFound(w, r)
		return
	}

	response := PointsResponse{Points: points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ---------- main ----------

func main() {
	router := setupRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
