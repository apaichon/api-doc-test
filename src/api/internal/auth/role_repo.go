package auth

import (
	"api/internal/db"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// RoleRepo represents the repository for role operations
type RoleRepo struct {
	DB *db.DB
}

// NewRoleRepo creates a new instance of RoleRepo
func NewRoleRepo() *RoleRepo {
	db := db.NewDB()
	return &RoleRepo{DB: db}
}

// InsertRole inserts a new role into the database
func (rr *RoleRepo) InsertRole(role *Role) (int64, error) {
	result, err := rr.DB.Exec("INSERT INTO role (role_name, role_desc, is_super_admin, created_at, created_by, status_id) VALUES (?, ?, ?, ?, ?, ?)",
		role.RoleName, role.RoleDesc, role.IsSuperAdmin, role.CreatedAt, role.CreatedBy, role.StatusID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetRoleByID retrieves a role by its ID
func (rr *RoleRepo) GetRoleByID(roleID int) (*Role, error) {
	row, err := rr.DB.QueryRow("SELECT * FROM role WHERE role_id = ?", roleID)

	if err != nil {
		return nil, fmt.Errorf("error query: %v", err)
	}

	var role Role
	err = row.Scan(
		&role.RoleID,
		&role.RoleName,
		&role.RoleDesc,
		&role.IsSuperAdmin,
		&role.CreatedAt,
		&role.CreatedBy,
		&role.StatusID,
	)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// UpdateRole updates an existing role
func (rr *RoleRepo) UpdateRole(role *Role) (int64, error) {
	result, err := rr.DB.Exec("UPDATE role SET role_name = ?, role_desc = ?, is_super_admin = ?, created_at = ?, created_by = ?, status_id = ? WHERE role_id = ?",
		role.RoleName, role.RoleDesc, role.IsSuperAdmin, role.CreatedAt, role.CreatedBy, role.StatusID, role.RoleID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// DeleteRole deletes a role by its ID
func (rr *RoleRepo) DeleteRole(roleID int) (int64, error) {
	result, err := rr.DB.Exec("DELETE FROM role WHERE role_id = ?", roleID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertUserRole inserts a new user-role relationship
func (rr *RoleRepo) InsertUserRole(userRole *UserRoles) (int64, error) {
	result, err := rr.DB.Exec("INSERT INTO user_roles (role_id, user_id, created_at, created_by) VALUES (?, ?, ?, ?)",
		userRole.RoleID, userRole.UserID, time.Now(), userRole.CreatedBy)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserRoleByID retrieves a user role by its ID
func (rr *RoleRepo) GetUserRoleByID(userRoleID int) (*UserRoles, error) {
	row, err := rr.DB.QueryRow("SELECT * FROM user_roles WHERE user_role_id = ?", userRoleID)

	if err != nil {
		return nil, err
	}

	var userRole UserRoles
	err = row.Scan(
		&userRole.UserRoleID,
		&userRole.RoleID,
		&userRole.UserID,
		&userRole.CreatedAt,
		&userRole.CreatedBy,
	)

	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

// GetUserRoleByID retrieves a user role by its ID
func (rr *RoleRepo) GetUserRoleByUserID(userID int) ([]*UserRoles, error) {
	rows, err := rr.DB.Query("SELECT * FROM user_roles WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var userRoles []*UserRoles
	for rows.Next() {
		var userRole UserRoles
		err := rows.Scan(
			&userRole.UserRoleID,
			&userRole.RoleID,
			&userRole.UserID,
			&userRole.CreatedAt,
			&userRole.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, &userRole)
	}
	return userRoles, nil
}

// GetUserRoleResolveByID retrieves a user role by its ID
func (rr *RoleRepo) GetUserIsSuperAdminByUserID(userID int) (bool, error) {
	row, err := rr.DB.QueryRow(`select sum(is_super_admin) as total_super_admin
	from vw_user_role_permissions 
	where user_id =?`, userID)

	isSuperAdmin := false
	if err != nil {
		return isSuperAdmin, err
	}

	totalSuperAdmin := 0

	err = row.Scan(
		&totalSuperAdmin,
	)

	if err != nil {
		return isSuperAdmin, err
	}

	if totalSuperAdmin > 0 {
		isSuperAdmin = true
	}

	return isSuperAdmin, nil
}

