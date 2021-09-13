package receipt

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const receiptPath = "receipts"

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	// Function to retrive the list of file from upload dir

	switch r.Method {
	case http.MethodGet:
		receipt, err := GetReceipt()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(receipt)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPost:
		// uploadd multi part form type limiting the size
		r.ParseMultipartForm(5 << 20)
		file, handler, err := r.FormFile("receipt") //key for the file upload
		value := r.FormValue("program")
		log.Println("This is the form:-> ", value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer file.Close()
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()
		io.Copy(f, file)

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func handleDownloads(w http.ResponseWriter, r *http.Request) {
	urlpathsegment := strings.Split(r.URL.Path, "receipts/")
	if len(urlpathsegment[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filename := urlpathsegment[1:][0]

	file, err := os.Open(filepath.Join(ReceiptDirectory, filename))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	fHeader := make([]byte, 512)
	file.Read(fHeader)
	// figureout what kind of file is this ?
	fContentType := http.DetectContentType(fHeader)
	// check the file size
	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fSize := strconv.FormatInt(stat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Lenght", fSize)

	file.Seek(0, 0) // set seeker to start of the file
	io.Copy(w, file)

	//

}

// func SetupRoutes(apiBasePath string) {
// 	receiptHandler := http.HandlerFunc(handleReceipts)
// 	downloadHandler := http.HandlerFunc(handleDownloads)
// 	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.Middleware((receiptHandler)))
// 	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptPath), cors.Middleware((downloadHandler)))
// }
