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
			slog.Warn("unable to get the jwt from context")
			util.WriteResponse(w, http.StatusUnauthorized, util.ErrorResponse(fmt.Errorf("unauthorized")))
			return
		}

		query := r.URL.Query()

		// Defaults
		page := int64(1)
		limit := int64(5)

		// Parse page
		if query.Has("page") {
			if parsedPage, err := util.ParseInt(query.Get("page")); err == nil && parsedPage > 0 {
				page = parsedPage
			} else {
				slog.Warn("invalid page, using default: 1")
			}
		}

		// Parse limit
		if query.Has("limit") {
			if parsedLimit, err := util.ParseInt(query.Get("limit")); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			} else {
				slog.Warn("invalid limit, using default: 10")
			}
		}

		offset := (page - 1) * limit

		//! sort

		sortField := "id" // default sort field
		sortOrder := "ASC"

		if query.Has("sortField") {
			sortField = query.Get("sortField")
		}

		if query.Has("sortOrder") {
			sortOrder = query.Get("sortOrder")
		}

		newSortfield := types.ValidSortFields[sortField]

		products, err := storage.GetAllProduct(claim.ShopID, offset, limit, sortOrder, newSortfield)
		if err != nil {
			slog.Error("unable to get products", "error", err)
			util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("unable to get products")))
			return
		}

		util.WriteResponse(w, http.StatusOK, map[string]interface{}{
			"success":   true,
			"products":  products,
			"page":      page,
			"limit":     limit,
			"sortOrder": sortOrder,
			"sortField": sortField,
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
		err = storage.DeleteUser(newId, claim.ShopID)
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

		response, err := storage.UpdateProduct(newId, product.Name, float32(product.Price), int(product.Price), claim.ShopID)

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
