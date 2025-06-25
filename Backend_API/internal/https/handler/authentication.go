package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/slangeres/Vypaar/backend_API/internal/storage"
	"github.com/slangeres/Vypaar/backend_API/internal/token"
	"github.com/slangeres/Vypaar/backend_API/internal/types"
	"github.com/slangeres/Vypaar/backend_API/internal/util"
)

func LoginUser(storage storage.UserStorage, jwtMaker *token.JwtMaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userCred types.Login

		// Decode JSON body
		err := json.NewDecoder(r.Body).Decode(&userCred)
		if errors.Is(err, io.EOF) {
			util.WriteResponse(w, http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("request body is empty")))
			return
		}
		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request format"))
			return
		}

		// Validate struct
		err = validator.New().Struct(userCred)
		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, "Validation Error")
			return
		}

		// Authenticate user
		shopID, err := storage.Login(userCred.Email, userCred.Password)
		if err != nil {
			util.WriteResponse(w, http.StatusUnauthorized, util.ErrorResponse(fmt.Errorf("invalid credentials")))
			return
		}

		// Generate JWT Token
		tokenString, err := jwtMaker.GenerateToken(userCred.Email, shopID, time.Hour*24) // token valid for 24 hours
		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, util.ErrorResponse(fmt.Errorf("failed to generate token")))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:   "token",
			Value:  tokenString,
			MaxAge: int((24 * time.Hour).Seconds()),
		})

		// Success response
		util.WriteResponse(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "Login successful",
		})
	}
}

func SignupUser(storage storage.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user types.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if errors.Is(err, io.EOF) {
			//resppnse req body is missing
			util.WriteResponse(w, http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("response body is missing")))
			return
		}
		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, fmt.Errorf("request error"))
			return
		}

		err = validator.New().Struct(user)
		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, "Validation Error Occur")
			return
		}
		newShopID := uuid.NewString()
		id, err := storage.Signup(user.Name, user.Email, user.Password, newShopID)

		if err != nil {
			util.ErrorResponse(fmt.Errorf("unable to sign up"))
			return

		}
		util.WriteResponse(w, http.StatusCreated, map[string]any{
			"sucess":       "True",
			"message":      "Sign up sucessful",
			"user_created": id,
		})
	}
}
