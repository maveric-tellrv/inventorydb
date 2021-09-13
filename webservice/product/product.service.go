package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/maveric-tellrv/inventoryservice/cors"
	"golang.org/x/net/websocket"
)

const productBasePath = "products"

func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(ProductsHandler)
	handleProduct := http.HandlerFunc(ProductHandler)
	handleReport := http.HandlerFunc(handleProductReport)

	http.Handle("/websocket", websocket.Handler(productSocket))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productBasePath), cors.Middleware((handleProducts)))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productBasePath), cors.Middleware(handleProduct))
	http.Handle(fmt.Sprintf("%s/%s/reports", apiBasePath, productBasePath), cors.Middleware(handleReport))
}

func Findproduct() []byte {
	b, _ := json.Marshal(Productlist)
	return b
}

func Handlebar(w http.ResponseWriter, r *http.Request) {
	w.Write(Findproduct())
}
func GetnextId() int {
	highestId := -1
	for _, p := range Productlist {
		if highestId < p.ProductID {
			highestId = p.ProductID
		}
	}
	return highestId + 1
}

func findproductByID(id int) (*Product, int) {
	for i, v := range Productlist {
		if v.ProductID == id {
			return &v, i
		}
	}
	return nil, 0
}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	urlpathsegment := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlpathsegment[len(urlpathsegment)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// product, listItemIndex := findproductByID(productID)
	product, _ := getProduct(productID)
	if product == nil {
		p := Product{}
		x, _ := json.Marshal(p)
		w.Header().Set("Content-Type", "application/json")
		w.Write(x)
		w.WriteHeader(http.StatusNotFound)
		return

	}

	switch r.Method {
	case http.MethodGet:
		// return a single product
		productsJson, err := json.Marshal(product)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)

	case http.MethodPut:
		// update the product
		var updateproduct Product
		bodyByte, error := ioutil.ReadAll(r.Body)
		log.Println(string(bodyByte))
		if error != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := json.Unmarshal(bodyByte, &updateproduct)
		log.Printf("Updated product %v", updateproduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// log.Printf("Updated productID %v %v", productID, listItemIndex)
		if updateproduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product = &updateproduct
		// log.Printf("Product to be updated %v\n", product)
		// log.Printf("Product to be updated %v\n", Productlist[listItemIndex])

		// addOrUpdateProduct(updateProduct)
		err = updateProduct(updateproduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Productlist[listItemIndex] = *product
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		_ = removeProduct(productID)
		w.WriteHeader(http.StatusOK)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Productlist, err := getProducts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ProductsJson, err := json.Marshal(Productlist)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(ProductsJson)

	case http.MethodPost:
		var newProducts Product
		bodyBtes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(bodyBtes, &newProducts)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProducts.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := insertProduct(newProducts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Added new product with Id:", id)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(strconv.Itoa(id)))
		// newProducts.ProductID = GetnextId()
		// Productlist = append(Productlist, newProducts)
		// b, _ := json.Marshal(newProducts)
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}

}
