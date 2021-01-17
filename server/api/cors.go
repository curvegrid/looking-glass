package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var domains []string

// CORSMiddleware handles the Multibaas app's CORS options.
func CORSMiddleware() echo.MiddlewareFunc {
	allowOriginFunc := func(origin string) (bool, error) {
		// in theory this check should be done in the CORS middleware but it appears not to be
		if origin == "" {
			return false, nil
		}

		// try to match the allowed origins
		for _, domain := range domains {
			if matchCORSOrigin(origin, domain) {
				return true, nil
			}
		}
		return false, nil
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc:  allowOriginFunc,
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete,
			http.MethodPut, http.MethodHead, http.MethodOptions},
	})
}

// matchCORSOrigin returns true if the provided origin matches the provided allowed origin.
// The allowed origin supports up to one wildcard `*` character.
func matchCORSOrigin(origin, allowed string) bool {
	// if the allowed origin contains a wildcard, we split it in
	// two: start and end string without the * and match each parts.
	if i := strings.IndexByte(allowed, '*'); i >= 0 {
		prefix := allowed[0:i]
		suffix := allowed[i+1:]
		return len(origin) >= len(prefix)+len(suffix) && strings.HasPrefix(origin, prefix) && strings.HasSuffix(origin, suffix)
	}
	return origin == allowed
}
