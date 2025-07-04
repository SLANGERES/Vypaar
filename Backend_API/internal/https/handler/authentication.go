package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"time"

	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/slangeres/Vypaar/backend_API/internal/https/middleware"
	RabbitMQ "github.com/slangeres/Vypaar/backend_API/internal/rabbitMQ"
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

func SignupUser(storage storage.UserStorage, rds *redis.Client) http.HandlerFunc {
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
		//Verigy gmail
		// Generate a 5-digit random number between 10000 and 99999
		otp := rand.Intn(90000) + 10000

		//!Adding otp in redis
		key := "otp:" + user.Email
		rds.Set(context.Background(), key, otp, 5*time.Minute)

		//push this otp in redis ....to check while verifying
		RabbitMQ.SendMail(user.Email, otp)

		util.WriteResponse(w, http.StatusCreated, map[string]any{
			"sucess":       "True",
			"message":      "Sign up sucessful",
			"user_created": id,
		})
	}
}

func VerifyUser(storage storage.UserStorage, rds *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			Email string `json:"email"`
			Otp   string `json:"otp"`
		}

		// Parse request body
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			util.WriteResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		key := "otp:" + reqBody.Email // fixed typo: "opt" -> "otp"

		// Get OTP from Redis
		otp, err := rds.Get(context.Background(), key).Result()
		if err == redis.Nil {
			util.WriteResponse(w, http.StatusRequestTimeout, "OTP expired or not found")
			//
			return
		} else if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, "Redis error: "+err.Error())
			return
		}
		if reqBody.Otp != otp {
			util.WriteResponse(w, http.StatusUnauthorized, "Invalid OTP")
			return
		}

		// Delete OTP from Redis
		rds.Del(context.Background(), key)

		// Mark user as verified in DB
		if _, err := storage.VerifyEmail(reqBody.Email); err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, "DB error: "+err.Error())
			return
		}

		util.WriteResponse(w, http.StatusOK, map[string]string{
			"message": "Account successfully verified",
		})
	}
}
func GenerateOTP(rds *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claim, ok := middleware.UserClaimsFromContext(r.Context())
		if !ok {
			slog.Error("Unable to extract JWT claims")
			util.WriteResponse(w, http.StatusUnauthorized, "Unauthorized request")
			return
		}

		email := claim.Subject
		otp := rand.Intn(90000) + 10000 // Generates a 5-digit OTP

		// Save OTP to Redis with TTL
		key := "otp:" + email
		err := rds.Set(context.Background(), key, otp, 5*time.Minute).Err()
		if err != nil {
			slog.Error("Failed to store OTP in Redis", "error", err)
			util.WriteResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		// Send OTP via RabbitMQ
		RabbitMQ.SendMail(email, otp)

		util.WriteResponse(w, http.StatusOK, "OTP sent successfully")
	}
}
