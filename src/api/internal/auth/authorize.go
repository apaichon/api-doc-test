package auth

import (
	"api/config"
	"errors"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

var roleRepo RoleRepo
var appConfig config.Config

type ContextKey string

const userKey = ContextKey("user")

func init() {
	roleRepo = *NewRoleRepo()
	appConfig = *config.NewConfig()
}

// AuthorizeWorkflow is a struct that holds the result and errors of the authorization workflow
type AuthorizeWorkflow struct {
	errors []error     // List of errors encountered during the workflow
	result interface{} // The result of the workflow operations
}

// IsSuperAdmin checks if the user is a super admin
func (auth *AuthorizeWorkflow) IsSuperAdmin(userId int) *AuthorizeWorkflow {
	// roleRepo := NewRoleRepo()
	isSuperAdmin, err := roleRepo.GetUserIsSuperAdminByUserID(userId)
	if err != nil {
		auth.addError(err)
	}
	auth.setResult(isSuperAdmin)
	return auth
}

// GetUserIDFromToken retrieves the user ID from the token
func (auth *AuthorizeWorkflow) GetUserIDFromToken(p graphql.ResolveParams) *AuthorizeWorkflow {
	userKey := ContextKey("user")
	tokenString, _ := p.Context.Value(userKey).(string)
	claims, err := DecodeJWTToken(tokenString, appConfig.SecretKey)

	if err != nil {
		auth.addError(errors.New("token expired"))
	}
	auth.setResult(claims)
	return auth
}

// GetResult returns the result of the authorization workflow
func (auth *AuthorizeWorkflow) GetResult() interface{} {
	return auth.result
}

// GetError returns the errors of the authorization workflow
func (auth *AuthorizeWorkflow) GetError() interface{} {
	joinedErr := errors.Join(auth.errors...)
	return joinedErr
}

// addError adds an error to the authorization workflow
func (auth *AuthorizeWorkflow) addError(err error) *AuthorizeWorkflow {
	if err != nil {
		auth.errors = append(auth.errors, err)
	}
	return auth
}

// setResult sets the result of the authorization workflow
func (auth *AuthorizeWorkflow) setResult(res interface{}) *AuthorizeWorkflow {
	auth.result = res
	return auth
}

// GetUserName retrieves the user name from the token
func GetUserName(p graphql.ResolveParams) (JwtClaims, error) {
	tokenString, _ := p.Context.Value(userKey).(string)
	config := config.NewConfig()
	claim, err := DecodeJWTToken(tokenString, config.SecretKey)
	return *claim, err
}

// GetUserPermission retrieves the user permission from the token
func GetUserPermission(r *http.Request) ([]*UserPermissionView, error) {
	tokenString := getTokenFromRequest(r)
	claims, err := DecodeJWTToken(tokenString, config.SecretKey)
	if err != nil {
		return nil, err
	}
	userPermissionRepo := NewUserPermissionRepo()
	userPermission, err := userPermissionRepo.GetUserPermissionView(claims.UserID)
	if err != nil {
		return nil, err
	}
	return userPermission, nil
}

// HasUserApiPermission checks if the user has the API permission
func HasUserApiPermission(r *http.Request) (bool, error) {
	tokenString := getTokenFromRequest(r)
	claims, err := DecodeJWTToken(tokenString, config.SecretKey)
	if err != nil {
		return false, err
	}

	userPermissionRepo := NewUserPermissionRepo()

	// fmt.Println("claims", claims.UserID)

	isSuperAdmin, err := userPermissionRepo.IsSuperAdmin(claims.UserID)
	fmt.Printf("isSuperAdmin:%v, err:%v", isSuperAdmin, err)
	if err != nil {
		fmt.Printf("err:%v", err)
		return false, err
	}

	if isSuperAdmin {
		return true, nil
	}

	userPermission, err := userPermissionRepo.GetUserApiPermissionView(claims.UserID, r.URL.Path)
	// fmt.Printf("userPermission:%v, err:%v", userPermission, err)

	if err != nil {
		return false, err
	}
	method := r.Method

	// fmt.Printf("userPermission:%v", userPermission)

	switch method {
	case http.MethodGet:
		return userPermission.CanRead, nil
	case http.MethodPost:
		return userPermission.CanWrite, nil
	case http.MethodPut:
		return userPermission.CanWrite, nil
	case http.MethodDelete:
		return userPermission.CanDelete, nil
	}
	return false, nil
}

// AuthorizeUserMiddleware is a middleware function that checks if the user has the API permission
func AuthorizeUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors(w, r)

		authorized, err := HasUserApiPermission(r)
		if err != nil || !authorized {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("authorized", authorized)

		next.ServeHTTP(w, r)
	})
}
