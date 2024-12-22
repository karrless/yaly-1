package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"yaly-1/pkg/errs"
)

type CalcService interface {
	Calc(ex string) (float64, error)
}

type CalcHandlers struct {
	service CalcService
}

func NewCalcHandlers(service CalcService) *CalcHandlers {
	return &CalcHandlers{
		service: service,
	}
}

type calcRequest struct {
	Expression string `json:"expression"`
}

type calcResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (ch *CalcHandlers) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req calcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := ch.service.Calc(req.Expression)
	if err != nil {
		if errors.Is(err, errs.ErrExpressionNotValid) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			err = errs.ErrInternal
		}

		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calcResponse{Result: result})

}
