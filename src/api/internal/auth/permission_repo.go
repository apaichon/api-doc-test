package auth

import (
	"api/internal/db"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type UserPermissionRepo struct {
	DB *db.DB
}

// NewUserPermissionRepo creates a new instance of UserPermissionRepo
func NewUserPermissionRepo() *UserPermissionRepo {
	db := db.NewDB()
	return &UserPermissionRepo{DB: db}
}

func (cr *UserPermissionRepo) GetUserPermissionView(userID int) ([]*UserPermissionView, error) {
	var users []*UserPermissionView

	query := fmt.Sprintf(`
            SELECT * FROM vw_user_permissions
             Where user_id = ?
        `, userID)

	rows, err := cr.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user UserPermissionView
		err := rows.Scan(
			&user.UserID,
			&user.RoleID,
			&user.RoleName,
			&user.IsSuperAdmin,
			&user.RolePermissionID,
			&user.ResourceTypeID,
			&user.ResourceName,
			&user.CanExecute,
			&user.CanRead,
			&user.CanWrite,
			&user.CanDelete,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (cr *UserPermissionRepo) IsSuperAdmin(userID int) (bool, error) {
	query := `select count(*) as is_super_admin from user_roles ur left join roles r 
			 on ur.role_id = r.role_id
              where ur.user_id = ? and r.is_super_admin =1;`

	row, err := cr.DB.QueryRow(query, userID)
	if err != nil {
		// fmt.Printf("count query err:%v", err)
		return false, err
	}

	var isSuperAdmin int
	err = row.Scan(&isSuperAdmin)
	if err != nil {
		fmt.Printf("scan err:%v", err)
		return false, err
	}
	// fmt.Printf("isSuperAdmin:%v", isSuperAdmin)
	return isSuperAdmin > 0, nil
}

func (cr *UserPermissionRepo) GetUserApiPermissionView(userID int, resourceName string) (*UserPermissionView, error) {

	query := `SELECT sum(can_execute) as can_execute, sum(can_read) as can_read, sum(can_write) as can_write, sum(can_delete) as can_delete FROM vw_user_permissions
             Where user_id = ? AND resource_type_id = 1 AND resource_name = ?`

	row, err := cr.DB.QueryRow(query, userID, resourceName)

	if err != nil {
		// fmt.Printf("query err:%v", err)
		return nil, err
	}

	var user UserPermissionView

	err = row.Scan(
		&user.CanExecute,
		&user.CanRead,
		&user.CanWrite,
		&user.CanDelete,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
