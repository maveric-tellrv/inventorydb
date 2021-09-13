package receipt

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

var ReceiptDirectory string = filepath.Join("uploads")

type Receipt struct {
	ReceiptName string    `json:"name"`
	UploadTime  time.Time `json:"uploadDate"`
}

func GetReceipt() ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	// reading the directory files
	file, err := ioutil.ReadDir(ReceiptDirectory)
	if err != nil {
		return nil, err
	}
	for _, f := range file {
		receipts = append(receipts, Receipt{ReceiptName: f.Name(), UploadTime: f.ModTime()})

	}

	return receipts, nil
}
