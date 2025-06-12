package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/slangeres/Vypaar/backend_API/internal/storage"
	"github.com/slangeres/Vypaar/backend_API/internal/types"
	"github.com/slangeres/Vypaar/backend_API/internal/util"
)

func PostProduct(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data types.Product

		err := json.NewDecoder(r.Body).Decode(&data)
		if errors.Is(err, io.EOF) {
			util.WriteResponse(w, http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("request body is missing")))
			return
		}
		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("invalid request body: %w", err)))
			return
		}

		//!NEED TO VALIDATE FIRST

		err = validator.New().Struct(data)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				fmt.Println(validationErrors)
				util.WriteResponse(w, http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("validation error occus")))
			} else {
				// Optional: catch other validation-related errors
				util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("validation error")))
			}
			return
		}

		id, err := storage.CreateProduct(data.Name, float32(data.Price), data.Quantity)
		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("failed to create product: %w", err)))
			return
		}

		util.WriteResponse(w, http.StatusCreated, map[string]string{
			"success":    "ok",
			"created_id": strconv.FormatInt(id, 10),
			"message":    "Product created successfully",
		})
	}
}

func GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.WriteResponse(w, http.StatusOK, "Getting the product")
	}
}
func GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.WriteResponse(w, http.StatusOK, "Getting by id ")
	}
}
func DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.WriteResponse(w, http.StatusOK, "Delete Product")
	}
}

func UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.WriteResponse(w, http.StatusOK, "Update product")
	}
}
