package healthcheck

import (
	"io"
	"net/http"
)

// HandlerHealthcheck : A healthcheck to confirm the underlying system is performing as expected.
func HandlerHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"alive": true}`)
}
