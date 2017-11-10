package java

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func AppendJavaDownloadMetadataRouter(r *mux.Router) {
	// Hmm, still missing the type group :  (?:\\.){type:\\w*}
	r.HandleFunc("/{group:.+}/{artifact:.+}/maven-metadata.xml", javaDownloadMetadataHandler)
}

func javaDownloadMetadataHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"metadata": {"group": "`+vars["group"]+`", "artifact": "`+vars["artifact"]+`"}}`)
}
