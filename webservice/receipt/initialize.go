package receipt

import (
	"fmt"
	"os"
)

func init() {
	fmt.Println("Check If Upload Dir exists ....")
	UploadExists()

}

func UploadExists() {
	_, err := os.Stat(ReceiptDirectory)
	if err == nil {
		fmt.Println(" Upload Dir exists ....")
	}
	if os.IsNotExist(err) {
		fmt.Println(" Upload Dir Does not exists Creaing Now ....")
		err := os.MkdirAll(ReceiptDirectory, 0755)
		if err != nil {
			fmt.Println(" Failed creating Dir ....")
		}
	}

}
