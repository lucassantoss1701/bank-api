package responses

import (
	"encoding/json"
	"log"
	"lucassantoss1701/bank/internal/entity"
	"net/http"
)

func encode(w http.ResponseWriter, data interface{}) {
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Fatal(err)
		}
	}
}

func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	encode(w, data)
}

func Err(w http.ResponseWriter, err error) {
	statusCode := 500

	if errorHandler, ok := err.(*entity.ErrorHandler); ok {
		switch errorHandler.TypeError {
		case entity.INTERNAL_ERROR:
			statusCode = http.StatusInternalServerError
		case entity.ENTITY_ERROR:
			statusCode = http.StatusUnprocessableEntity
		case entity.NOT_FOUND_ERROR:
			statusCode = http.StatusNotFound
		case entity.UNAUTHORIZED_ERROR:
			statusCode = http.StatusUnauthorized
		case entity.NOT_ALLOWED_ERROR:
			statusCode = http.StatusMethodNotAllowed
		case entity.BAD_REQUEST:
			statusCode = http.StatusBadRequest
		case entity.CONFLICT_ERROR:
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	encode(w, struct {
		Erro string `json:"error"`
	}{
		Erro: err.Error(),
	})

}
