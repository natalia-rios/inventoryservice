package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/riosnatalia/inventoryservice/cors"
	"golang.org/x/net/websocket"
)

const productsBasePath = "products"

func handleProduct(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		product, err := getProduct(productID)
		if product == nil {
			w.WriteHeader(http.StatusNotFound)
		}

		j, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPut:
		var product Product

		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if product.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = updateProduct(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		removeProduct(productID)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList, err := getProductList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(productList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPost:
		var product Product

		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = insertProduct(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(handleProducts)
	handleProduct := http.HandlerFunc(handleProduct)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsBasePath), cors.Middleware(handleProducts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), cors.Middleware(handleProduct))
	http.Handle("/websocket", websocket.Handler(productSocket))
}
