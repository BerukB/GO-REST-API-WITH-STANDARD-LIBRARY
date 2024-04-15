package middleware

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/common"
)

func SetCtxRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cxtFromRequest := r.Context()

		RequestID := strconv.Itoa(rand.Intn(100000000))

		ctxWithReqID := context.WithValue(cxtFromRequest, common.RequestIDKey, RequestID)

		r = r.WithContext(ctxWithReqID)

		next.ServeHTTP(w, r)
	})
}
