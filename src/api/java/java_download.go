package api

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func javaDownloadMetadataHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"metadata": true}`)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{group:.+}/{artifact:.+}/maven-metadata.xml{type:(\\.)?(\\w)*}", javaDownloadMetadataHandler)
	http.ListenAndServe(":8080", nil)
}
