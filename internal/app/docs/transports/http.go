package transports

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/rs/cors"

	httpSwagger "github.com/swaggo/http-swagger"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler() http.Handler { // Zipkin HTTP Server Trace can either be instantiated per endpoint with a
	m := bone.New()
	m.Get("/*", httpSwagger.Handler(
		httpSwagger.URL("https://mask.goodideas-studio.com/docs/doc.json"), //The url pointing to API definition"
	))
	return cors.AllowAll().Handler(m)
}
