package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/maveric-tellrv/inventoryservice/database"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("Loading products.....")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Product loaded Successfully \n")
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	// Check if file exists
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exists", fileName)
	}

	// read the file byte
	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)

	// create a slice of product list
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	for _, v := range productList {
		prodMap[v.ProductID] = v
	}

	return prodMap, nil

}

func getProduct(productId int) (*Product, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	productId, 
	manufacturer, 
	sku, 
	upc, 
	quantityOnHand, 
	pricePerUnit, 
	productName 
	FROM products
	WHERE productId = ?`, productId)

	product := &Product{}
	err := row.Scan(&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.QuantityOnHand,
		&product.PricePerUnit,
		&product.ProductName,
	)
	if err == sql.ErrNoRows {
		log.Println("No Rows found for Prodduct ID: ", productId)
		return nil, nil
	}
	return product, nil

	// productMap.RLock()
	// defer productMap.RUnlock()
	// if product, ok := productMap.m[productId]; ok {
	// 	return &product
	// }
	// return nil
}

func removeProduct(productId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE 
	from products 
	where 
	productId = ?`,
		productId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
	// productMap.RLock()
	// defer productMap.RUnlock()
	// if _, ok := productMap.m[productId]; ok {
	// 	delete(productMap.m, productId)
	// 	log.Printf("Delted Product with ID: %v", productId)
	// }

}

func getProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// connection to database
	result, err := database.DbConn.QueryContext(ctx, `SELECT 
	productId, 
	manufacturer, 
	sku, 
	upc, 
	quantityOnHand, 
	pricePerUnit, 
	productName 
	FROM products`)
	if err != nil {
		log.Println(err)
		return nil, err

	}
	defer result.Close()

	products := make([]Product, 0)
	for result.Next() {
		var product Product
		result.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.QuantityOnHand,
			&product.PricePerUnit,
			&product.ProductName,
		)
		products = append(products, product)
	}
	return products, nil

	// productMap.RLock()
	// var productList = make([]Product, 0, len(productMap.m))
	// for _, v := range productMap.m {
	// 	productList = append(productList, v)
	// }
	// productMap.RUnlock()
	// return productList

}

func GetTopTenProducts() ([]Product, error) {
	// Return only top 10 prodcuts for websocket function
	log.Println("GetToTenProduct Called from the WebSocket...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// connection to database
	result, err := database.DbConn.QueryContext(ctx, `SELECT 
	productId, 
	manufacturer, 
	sku, 
	upc, 
	quantityOnHand, 
	pricePerUnit, 
	productName 
	FROM products ORDER BY quantityOnHand DESC LIMIT 10`)
	if err != nil {
		log.Println(err)
		return nil, err

	}
	defer result.Close()

	products := make([]Product, 0)
	for result.Next() {
		var product Product
		result.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.QuantityOnHand,
			&product.PricePerUnit,
			&product.ProductName,
		)
		products = append(products, product)
	}
	return products, nil

}

func getProductIds() []int {

	productMap.RLock()
	productids := []int{}
	for key := range productMap.m {
		productids = append(productids, key)
	}
	productMap.RUnlock()
	sort.Ints(productids)
	return productids

}

func getNextProductId() int {
	productId := getProductIds()
	return productId[len(productId)-1] + 1

}

func updateProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(ctx, `UPDATE products SET
	manufacturer=?, 
	sku=?, 
	upc=?, 
	quantityOnHand=?, 
	pricePerUnit=?, 
	productName=?
	WHERE productId =? `,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.QuantityOnHand,
		product.PricePerUnit,
		product.ProductName,
		product.ProductID)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func insertProduct(product Product) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO products
	(manufacturer, 
	sku, 
	upc, 
	quantityOnHand, 
	pricePerUnit, 
	productName) VALUES (?, ?, ?, ?, ?, ?) `,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.QuantityOnHand,
		product.PricePerUnit,
		product.ProductName)
	if err != nil {
		log.Println("Inside Insert-> ", err)
		return 0, nil
	}
	inserID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	log.Printf("Added product With ID %v", inserID)
	return int(inserID), nil
}

func addOrUpdateProduct(product Product) (int, error) {

	addOrUpdateId := -1
	if product.ProductID > 0 {
		oldProduct, _ := getProduct(product.ProductID)
		if oldProduct == nil {
			return 0, fmt.Errorf("Product with PID %v does not exists \n", product.ProductID)
		}
		addOrUpdateId = product.ProductID
	} else {
		addOrUpdateId := getNextProductId()
		product.ProductID = addOrUpdateId
	}
	productMap.Lock()
	productMap.m[product.ProductID] = product
	productMap.Unlock()
	return addOrUpdateId, nil

}

// function in product report

func searchForProductData(productFilter ProductReportFilter) ([]Product, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var queryArgs = make([]interface{}, 0)

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`SELECT
	productId, 
	manufacturer, 
	sku, 
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products WHERE `)

	if productFilter.NameFilter != "" {
		queryBuilder.WriteString(`productName LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.NameFilter)+"%")
	}
	if productFilter.ManufacturerFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString("AND")
		}
		queryBuilder.WriteString(`manufacturer LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.ManufacturerFilter)+"%")
	}

	if productFilter.SkuFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString("AND")
		}
		queryBuilder.WriteString(`sku LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.SkuFilter)+"%")
	}

	log.Println("\nQuery recieved->\n\n", queryBuilder.String())
	log.Println("\n Query Args -> \n\n", queryArgs)
	result, err := database.DbConn.QueryContext(ctx, queryBuilder.String(), queryArgs...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer result.Close()
	products := make([]Product, 0)

	for result.Next() {
		var product Product

		result.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)

		products = append(products, product)
		log.Println(products)

	}
	return products, nil
}
