package auth

import (
	"api/config"
	"api/internal/handler"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GetRequestID(r *http.Request) string {
	requestID, ok := r.Context().Value("requestContext").(string)
	fmt.Printf("GetRequestID: Got RequestID: %s\n", requestID)
	if !ok {
		cookie, err := r.Cookie("request_id")
		if err == nil {
			requestID = cookie.Value
		} else {
			requestID = ""
		}
	}
	return requestID
}

// LoginHandler godoc
// @Summary User login
// @Description Authenticate a user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} types.ErrorResponse
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	cors(w, r)
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := handler.NewErrorResponse(
			http.StatusBadRequest,
			"Bad Request",
			"INVALID_REQUEST",
			err.Error(),
			GetRequestID(r),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userRepo := NewUserRepo()
	config := config.NewConfig()

	userdb, err := userRepo.GetUserByName(user.Username) // users[user.Username]
	if err != nil {
		resp := handler.NewErrorResponse(
			http.StatusUnauthorized,
			"Unauthorized",
			"INVALID_USERNAME",
			"Invalid username",
			GetRequestID(r),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}
	// fmt.Println("User")
	hash := HashString(user.Username + user.Password + userdb.Salt)
	// fmt.Println("UserDbPassword:", userdb.Password, "Hash:",hash, "UserNameDb:", userdb.Username, "username", user.Username )
	if userdb.Password != hash {
		resp := handler.NewErrorResponse(
			http.StatusUnauthorized,
			"Unauthorized",
			"INVALID_CREDENTIALS",
			"Invalid username or password",
			GetRequestID(r),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// fmt.Println("Duration", time.Duration(config.TokenAge))
	expirationTime := time.Now().Add(time.Duration(config.TokenAge) * time.Minute)

	claims := &JwtClaims{
		UserID:   userdb.UserID,
		Username: userdb.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		resp := handler.NewErrorResponse(
			http.StatusInternalServerError,
			"Internal Server Error",
			"TOKEN_GENERATION_FAILED",
			err.Error(),
			GetRequestID(r),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	// Updated response format
	responseData := JwtToken{
		Token:     tokenString,
		ExpiredAt: expirationTime.Unix(),
	}
	resp := handler.NewResponse(http.StatusOK, "Success", responseData, GetRequestID(r))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)

}

func AuthenticationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the JWT token is present in the request
		tokenString := getTokenFromRequest(r)
		// fmt.Println("tokenString:" + tokenString)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the JWT token
		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, tokenString)
		r = r.WithContext(ctx)

		// Pass the request to the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}

// RegisterHandler godoc
// @Summary Register new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} User
// @Failure 400 {object} types.ErrorResponse
// @Router /register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := handler.NewErrorResponse(
			http.StatusBadRequest,
			"Bad Request",
			"INVALID_REQUEST",
			"Invalid request body",
			GetRequestID(r),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if user.Username == "" || user.Password == "" {
		resp := handler.NewErrorResponse(
			http.StatusBadRequest,
			"Bad Request",
			"MISSING_FIELDS",
			"Username and password are required",
			GetRequestID(r),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userRepo := NewUserRepo()
	exists, err := userRepo.ExistsUserByName(user.Username)
	if err != nil {
		resp := handler.NewErrorResponse(
			http.StatusInternalServerError,
			"Internal Server Error",
			"DATABASE_ERROR",
			"Database operation failed",
			GetRequestID(r),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if exists {
		resp := handler.NewErrorResponse(
			http.StatusConflict,
			"Conflict",
			"USER_EXISTS",
			"Username already exists",
			GetRequestID(r),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	salt := generateSalt()
	hashedPassword := HashString(user.Username + user.Password + salt)

	newUser := User{
		Username:  user.Username,
		Password:  hashedPassword,
		Salt:      salt,
		CreatedAt: time.Now(),
		CreatedBy: user.Username,
		StatusID:  1,
	}

	err = userRepo.CreateUser(&newUser)
	if err != nil {
		resp := handler.NewErrorResponse(
			http.StatusInternalServerError,
			"Internal Server Error",
			"REGISTRATION_FAILED",
			"Failed to create user",
			GetRequestID(r),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	responseData := map[string]string{
		"username": user.Username,
		"message":  "User registered successfully",
	}

	resp := handler.NewResponse(http.StatusCreated, "Created", responseData, GetRequestID(r))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	// Clear the authentication cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour), // Set expiration in the past
	})

	// Create a success response
	resp := handler.NewResponse(http.StatusOK, "Success", map[string]string{"message": "Logged out successfully"}, GetRequestID(r))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)
}

// CreateRoleHandler godoc
// @Summary Create a new role
// @Description Create a new role in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param role body Role true "Role information"
// @Success 201 {object} Role
// @Failure 400 {object} types.ErrorResponse
// @Router /role [post]
func CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	var role Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		resp := handler.NewErrorResponse(http.StatusBadRequest, "Bad Request", "INVALID_REQUEST", "Invalid request body", GetRequestID(r))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	roleRepo := NewRoleRepo()
	roleRepo.InsertRole(&Role{
		RoleName:     role.RoleName,
		RoleDesc:     role.RoleDesc,
		IsSuperAdmin: role.IsSuperAdmin,
		CreatedAt:    time.Now(),
		CreatedBy:    role.CreatedBy,
		StatusID:     role.StatusID,
	})

	responseData := map[string]string{
		"role_name": role.RoleName,
		"message":   "Role created successfully",
	}

	resp := handler.NewResponse(http.StatusCreated, "Created", responseData, GetRequestID(r))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)
}

func CreateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	var userRole UserRoles
	err := json.NewDecoder(r.Body).Decode(&userRole)
	if err != nil {
		resp := handler.NewErrorResponse(http.StatusBadRequest, "Bad Request", "INVALID_REQUEST", "Invalid request body", GetRequestID(r))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userRoleRepo := NewUserRoleRepo()
	userRoleRepo.CreateUserRole(&userRole)

	responseData := map[string]string{
		"user_role_id": strconv.Itoa(userRole.UserRoleID),
		"message":      "User role created successfully",
	}

	resp := handler.NewResponse(http.StatusCreated, "Created", responseData, GetRequestID(r))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	userRepo := NewUserRepo()
	users, err := userRepo.GetUsersBySearchText(r.URL.Query().Get("searchText"), 10, 0)
	if err != nil {
		resp := handler.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", "DATABASE_ERROR", "Database operation failed", GetRequestID(r))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := handler.NewResponse(http.StatusOK, "Success", users, GetRequestID(r))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)
}
