package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/riosnatalia/inventoryservice/cors"
)

const receiptPath = "receipts"

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		receiptList, err := GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(receiptList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPost:
		r.ParseMultipartForm(5 << 20)
		file, handler, err := r.FormFile("receipt")
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
		w.WriteHeader(http.StatusCreated)

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filename := urlPathSegments[1:][0]
	file, err := os.Open(filepath.Join(ReceiptDirectory, filename))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)

	fileContentType := http.DetectContentType(fileHeader)
	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fileSize := strconv.FormatInt(stat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	file.Seek(0, 0)
	io.Copy(w, file)
}

func SetupRoutes(apiBasePath string) {
	receiptHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownload)
	http.Handle(fmt.Sprintf("%s%s", apiBasePath, receiptPath), cors.Middleware(receiptHandler))
	http.Handle(fmt.Sprintf("%s%s/", apiBasePath, receiptPath), cors.Middleware(downloadHandler))
}
