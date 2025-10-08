package utils
import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}

func ExtractID(path, prefix string) int {
	idStr := strings.TrimPrefix(path, prefix)
	id, _ := strconv.Atoi(idStr)
	return id
}
