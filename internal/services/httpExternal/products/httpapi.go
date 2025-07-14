package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Chandra5468/cfp-Products-Service/internal/types"
)

func ValidateCart(item *types.Item) *types.PurchasedProduct {
	// POST request

	URL := fmt.Sprintf("http://%s/v1/api/cart/%v/buy", os.Getenv("PRODUCTS_URL"), item.ProductId)

	data := &types.Item{
		ProductId: item.ProductId,
		Quantity:  item.Quantity,
	}
	// Marashall the struct to json byte
	jsonBody, err := json.Marshal(data)

	if err != nil {
		log.Printf("There is an error while sending purchase request to products %v", err)
		return nil
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Printf("There is an error while sending purchase request to products %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	// defer req.Body.Close() Do we need to do or will it either way happen by default

	// Create a http client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending the request %v", err)
		return nil
	}
	defer resp.Body.Close()

	purchasedProduct := &types.PurchasedProduct{}

	err = json.NewDecoder(resp.Body).Decode(purchasedProduct)
	if err != nil {
		log.Printf("Error unmarshalling response from products %v", err)
		return nil
	}
	log.Println("This is purchased product-----", purchasedProduct)
	return purchasedProduct
}
