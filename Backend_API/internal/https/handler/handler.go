package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/slangeres/Vypaar/backend_API/internal/https/middleware"
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

		//!jwt claims

		claims, ok := middleware.UserClaimsFromContext(r.Context())
		if !ok {
			slog.Error("unable to get the jwt claims")
		}

		shopID := claims.ShopID

		id, err := storage.CreateProduct(data.Name, float32(data.Price), data.Quantity, shopID)
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

func GetProduct(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claim, ok := middleware.UserClaimsFromContext(r.Context())
		if !ok {
			slog.Info("unable to get the jwt")
		}
		products, err := storage.GetAllProduct(claim.ShopID)
		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("unable to get the product")))
			return
		}

		// Construct a generic map[string]interface{} response
		util.WriteResponse(w, http.StatusOK, map[string]interface{}{
			"success":  true,
			"products": products,
		})
	}
}

func GetProductById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var product types.Product
		id := r.PathValue("id")
		if id == "" {
			util.ParameterMissing(w, http.StatusInternalServerError)
		}
		newId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("unable to parse id")))
			return
		}
		claim, ok := middleware.UserClaimsFromContext(r.Context())
		if !ok {
			slog.Info("Unable to get the claim")
		}
		product, err = storage.GetUserById(newId, claim.ShopID)

		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("unable to get the product")))
			return
		}

		util.WriteResponse(w, http.StatusOK, product)
	}
}
func DeleteProduct(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			util.ParameterMissing(w, http.StatusInternalServerError)
		}

		newId, err := util.ParseInt(id)
		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("unable to parse the id")))
			return
		}
		claim, ok := middleware.UserClaimsFromContext(r.Context())
		if !ok {
			slog.Info("Unable to get the claim")
		}
		err = storage.DeleteUser(newId,claim.ShopID)
		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("unable to delete product")))
			return
		}

		util.WriteResponse(w, http.StatusOK, map[string]interface{}{
			"message": "Product deleted successfully",
		})
	}
}

func UpdateProduct(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			util.ParameterMissing(w, http.StatusInternalServerError)
		}

		newId, err := util.ParseInt(id)
		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("unable to parse the id")))
			return
		}

		var product types.Product

		json.NewDecoder(r.Body).Decode(&product)
		claim, ok := middleware.UserClaimsFromContext(r.Context())
		if !ok {
			slog.Info("Unable to get the claim")
		}

		response, err := storage.UpdateProduct(newId, product.Name, float32(product.Price), int(product.Price),claim.ShopID)

		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("unable to update the product")))
			return
		}

		util.WriteResponse(w, http.StatusOK, map[string]interface{}{
			"sucess":   true,
			"response": response,
		})

	}
}
