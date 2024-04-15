package handlers

import (
	"context"
	"net/http"

	"github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/common"
)

func HandleRequestError(w http.ResponseWriter, r *http.Request) {

	// Retrieve the error from the context
	if err, ok := r.Context().Value(common.ErrorKey).(*common.CustomError); ok {

		http.Error(w, err.Message, err.StatusCode)
		return
	}

	// Default case
	http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
}

func HandleError(w http.ResponseWriter, r *http.Request, errType string, message string, statusCode int) {

	customErr := &common.CustomError{
		Type:       errType,
		Message:    message,
		StatusCode: statusCode,
	}

	ctxWithError := context.WithValue(r.Context(), common.ErrorKey, customErr)
	r = r.WithContext(ctxWithError)

	HandleRequestError(w, r)
}
