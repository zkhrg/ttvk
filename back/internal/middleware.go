package internal

import (
	"net/http"
)

func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// Проходим по всем миддлварям в обратном порядке, чтобы
	// первый миддлварь был самым внешним, а последний - самым внутренним
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
