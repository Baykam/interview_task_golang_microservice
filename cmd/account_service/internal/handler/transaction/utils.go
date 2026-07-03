package transActionHandler

import (
	"encoding/json"
	"net/http"
)

func (h *handler) respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
