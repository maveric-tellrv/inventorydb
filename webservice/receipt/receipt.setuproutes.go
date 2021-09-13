package receipt

import (
	"fmt"
	"net/http"

	"github.com/maveric-tellrv/inventoryservice/cors"
)

func SetupRoutes(apiBasePath string) {
	receiptHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownloads)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.Middleware((receiptHandler)))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptPath), cors.Middleware((downloadHandler)))
}
