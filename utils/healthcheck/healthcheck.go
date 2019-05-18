package healthcheck

import (
	"context"
	"net/http"

	"github.com/keitaroinc/goa"
)

func NewCheckMiddleware(healthcheckEndpoint string) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			// Healthcheck endpoint
			if req.URL.Path == healthcheckEndpoint {
				rw.Header().Set("Content-Type", "application/text")
				rw.Write([]byte("OK"))
				rw.WriteHeader(200)
				return nil
			}
			return h(ctx, rw, req)
		}
	}
}
