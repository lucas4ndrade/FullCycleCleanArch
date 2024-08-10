package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lucas4ndrade/FullcycleCleanArch/internal/usecase"
)

func CreateOrderHandler(uc *usecase.CreateOrderUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto usecase.CreateOrderInputDTO
		err := json.NewDecoder(r.Body).Decode(&dto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		output, err := uc.Execute(dto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ListOrderHandler(uc *usecase.ListOrderUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := parseListQuery(r)

		output, err := uc.Execute(dto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func parseListQuery(r *http.Request) (dto usecase.ListOrderInputDTO) {
	dto = usecase.GetDefaultOrderInputDTO()
	queries := r.URL.Query()

	if fromQuery := queries.Get("from"); fromQuery != "" {
		fromQueryInt, err := strconv.Atoi(fromQuery)
		if err == nil {
			dto.From = int64(fromQueryInt)
		}
	}

	if sizeQuery := queries.Get("size"); sizeQuery != "" {
		fromQueryInt, err := strconv.Atoi(sizeQuery)
		if err == nil {
			dto.Size = int64(fromQueryInt)
		}
	}
	return
}
