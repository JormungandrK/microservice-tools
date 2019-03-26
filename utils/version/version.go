package version

import (
		"context"
		"net/http"
		"encoding/json"

		"github.com/goadesign/goa"
	)

type Read struct {
	Version string
}

func NewVersionMiddleware(version, versionEndpoint string) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {

			endpoints := Read{version}
			js, err := json.Marshal(endpoints)
			
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return err
			}	

			//endpoint that returns the microservice version
			if (req.URL.Path == versionEndpoint) {
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(js)		
				rw.WriteHeader(200)	

				return nil
				}

			return h(ctx, rw, req)
			}
		}
	}
	