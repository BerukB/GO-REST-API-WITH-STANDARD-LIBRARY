package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cxtFromRequest := r.Context()

		ctx, cancel := context.WithTimeout(cxtFromRequest, 2*time.Second)

		r = r.WithContext(ctx)

		donech := make(chan struct{})

		go func() {
			defer close(donech)
			next.ServeHTTP(w, r)
		}()

		defer cancel()

		select {
		case <-donech:
			return
		case <-r.Context().Done():
			if ctx.Err() == context.DeadlineExceeded {
				w.WriteHeader(http.StatusRequestTimeout)
				json.NewEncoder(w).Encode("408 Status Request Timeout")
			}
		}

	})
}
