package http

import (
	"net/http"

	"back/internal"
	"back/usecase"
)

func SetupRoutes(api *usecase.Usecase, mux *http.ServeMux) {
	// Создаем цепочку миддлварей и передаем API через замыкание
	// нужно создать композитный хендлер-миддлварь в котором я буду использовать свитч по методу у ручки чтобы
	// выбирать обработку по нужной миддлеварине
	// т.е нужна миддлеварь распределитель
	handler := internal.ChainMiddleware(
		AddressHandler(api),
		// internal.GetMethodMiddleware,
	)

	mux.Handle("/api/address/", handler)
}
