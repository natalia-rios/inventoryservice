package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/riosnatalia/inventoryservice/database"
	"github.com/riosnatalia/inventoryservice/product"
	"github.com/riosnatalia/inventoryservice/receipt"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
