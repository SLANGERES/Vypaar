package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/slangeres/Vypaar/backend_API/internal/storage"
	"github.com/slangeres/Vypaar/backend_API/internal/types"
	"github.com/slangeres/Vypaar/backend_API/internal/util"
)

func LoginUser(storage storage.UserStorage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var userCred types.Login

		err := json.NewDecoder(r.Body).Decode(&userCred)

		if errors.Is(err, io.EOF) {
			util.WriteResponse(w, http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("response body is missing")))
			return
		}
		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, fmt.Errorf("request error"))
			return
		}
		err = validator.New().Struct(userCred)

		if err != nil {
			//!validation error
			util.WriteResponse(w, http.StatusBadRequest, "Validation Error Occur")
			return
		}
		id, err := storage.Login(userCred.Email, userCred.Password)

		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("404 user not found ")))
			return
		}

		util.WriteResponse(w, http.StatusOK, map[string]any{
			"sucess":  "true",
			"message": "Login sucessful",
			"Id":      id,
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
		id, err := storage.Signup(user.Name, user.Email, user.Password)

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
