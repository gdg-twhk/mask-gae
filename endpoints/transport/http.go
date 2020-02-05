package transport

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/rs/cors"

	"github.com/cage1016/mask/endpoint"
)

func MakeHandler() http.Handler {
	mux := bone.New()

	mux.PostFunc("/stores", endpoint.StoresHandler)
	mux.PostFunc("/sync", endpoint.SyncHandler)
	mux.GetFunc("/", endpoint.HomeHandler)

	return cors.Default().Handler(mux)
}
