package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"back/domain/entities"
	"back/usecase"
)

func AddressHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetAddresses(uc, w, r)
		case "POST":
			fmt.Println("POST METHOD")
			CreateAddress(uc, w, r)
		case "PATCH":
			EditAddress(uc, w, r)
		default:
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		}
	}
}

func GetAddresses(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {
	entity, err := uc.GetFullInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func EditAddress(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {
	var newAddressInfo entities.EntityRequest
	if err := json.NewDecoder(r.Body).Decode(&newAddressInfo); err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	answer, err := uc.EditAddressByIP(newAddressInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}

func CreateAddress(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {
	var newAddressInfo entities.EntityRequest
	fmt.Println("handler", newAddressInfo)
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&newAddressInfo); err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("handler2", newAddressInfo)
	defer r.Body.Close()

	answer, err := uc.CreateAddress(newAddressInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}
