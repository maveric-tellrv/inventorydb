package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maveric-tellrv/inventoryservice/database"
	"github.com/maveric-tellrv/inventoryservice/product"
	"github.com/maveric-tellrv/inventoryservice/receipt"
)

const apiBasePath = "/api"

func main() {
	database.SetUpDatabase()
	product.SetupRoutes(apiBasePath)
	receipt.SetupRoutes(apiBasePath)
	// http.HandleFunc("/bar", product.Handlebar)
	// http.HandleFunc("/products", product.ProductsHandler)
	// http.HandleFunc("/products/", product.ProductHandler)
	http.ListenAndServe(":5000", nil)
}
