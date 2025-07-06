package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type Invoice struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	DueDate  string  `json:"due_date"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var inv Invoice
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	var errs []string
	if inv.Amount <= 0 {
		errs = append(errs, "amount must be greater than 0")
	}

	if len(inv.Currency) != 3 {
		errs = append(errs, "currency must be a 3-letter code")
	}

	if _, err := time.Parse("2006-01-02", inv.DueDate); err != nil {
		errs = append(errs, "due_date must be in YYYY-MM-DD format")
	}

	if len(errs) > 0 {
		resp, _ := json.Marshal(map[string]any{
			"valid":  false,
			"errors": errs,
		})
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(resp)
		return
	}

	// resp, _ := json.Marshal(map[string]bool{"valid": true})
	json.NewEncoder(w).Encode(map[string]bool{"valid": true})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(resp)
}
