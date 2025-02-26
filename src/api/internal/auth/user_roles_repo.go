package auth

import "api/internal/db"

// UserRepo represents the repository for user operations
type UserRoleRepo struct {
	DB *db.DB
}

// NewUserRepo creates a new instance of UserRepo
func NewUserRoleRepo() *UserRoleRepo {
	db := db.NewDB()
	return &UserRoleRepo{DB: db}
}


func (r *UserRoleRepo) CreateUserRole(userRole *UserRoles) error {
	query := `INSERT INTO user_roles (role_id, user_id, created_at, created_by, status_id) VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, userRole.RoleID, userRole.UserID, userRole.CreatedAt, userRole.CreatedBy, userRole.StatusID)
	return err
}

