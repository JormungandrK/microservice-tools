package healthcheck

import (
		"context"
		"net/http"
		// "fmt"

		"github.com/goadesign/goa"
	)

func NewCheckMiddleware(healthcheckEndpoint string) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			
			//healthcheck endpoint
			if (req.URL.Path == healthcheckEndpoint) {
				rw.Header().Set("Content-Type", "application/text")
				rw.Write([]byte("OK"))	
				rw.WriteHeader(200)			
			}
			// return h(ctx, rw, req)
			return nil
			}
		}
	}	