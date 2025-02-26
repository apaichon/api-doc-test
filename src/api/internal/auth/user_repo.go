package auth

import (
	"api/internal/db"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// UserRepo represents the repository for user operations
type UserRepo struct {
	DB *db.DB
}

// NewUserRepo creates a new instance of UserRepo
func NewUserRepo() *UserRepo {
	db := db.NewDB()
	return &UserRepo{DB: db}
}

// Get Users fetches users from the database with support for text search, limit, and offset
func (cr *UserRepo) GetUsersBySearchText(searchText string, limit, offset int) ([]*User, error) {
	var users []*User

	query := fmt.Sprintf(`
            SELECT * FROM users
             Where user_name like '%%%s%%' OR password like '%%%s%%' OR salt like '%%%s%%'
            LIMIT ? OFFSET ?
        `, searchText, searchText, searchText)

	rows, err := cr.DB.Query(query, limit, offset)

	// fmt.Printf("searchText: %s \n result: %v \n error: %v", searchText, rows, err)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	time.Sleep(30 * time.Second)

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Password,
			&user.Salt,
			&user.CreatedAt,
			&user.CreatedBy,
			&user.StatusID,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// fmt.Println("users: ", len(users))

	return users, nil
}

// Get UserByID retrieves a user by its ID from the database
func (cr *UserRepo) GetUserByID(id int) (*User, error) {
	var user User
	// Execute query to get a user by ID from the database
	row, err := cr.DB.QueryRow("SELECT * FROM users WHERE user_id = ?", id)

	if err != nil {
		return &user, nil
	}

	row.Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.CreatedAt,
		&user.CreatedBy,
	)

	return &user, nil
}

// Get UserByID retrieves a user by its ID from the database
func (cr *UserRepo) GetUserByName(name string) (*User, error) {
	var user User
	// Execute query to get a user by ID from the database
	row, err := cr.DB.QueryRow("SELECT * FROM users WHERE user_name = ?", name)

	if err != nil {
		return nil, fmt.Errorf("error query: %v", err)
	}

	row.Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.StatusID,
	)

	return &user, nil
}

// Get UserByID retrieves a user by its ID from the database
func (cr *UserRepo) ExistsUserByName(name string) (bool, error) {
	// Execute query to get a user by ID from the database
	row, err := cr.DB.QueryRow("SELECT count(*) as count FROM users WHERE user_name = ?", name)

	if err != nil {
		return false, fmt.Errorf("error query: %v", err)
	}

	var count int

	row.Scan(
		&count,
	)

	return count > 0, nil
}

// Insert User inserts a new user into the database
func (cr *UserRepo) InsertUser(user *User) (int64, error) {
	// Execute insert query to insert a new user into the database
	result, err := cr.DB.Insert("INSERT INTO users (user_id,user_name,password,salt,created_at) VALUES ({?,?,?,?,?})",
		user.UserID, user.Username, user.Password, user.Salt, user.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update User updates an existing user in the database
func (cr *UserRepo) UpdateUser(user *User) (int64, error) {
	// Execute update query to update an existing user in the database
	result, err := cr.DB.Update("UPDATE users SET user_id=?,user_name=?,password=?,salt=? where user_id=?",
		user.UserID, user.Username, user.Password, user.Salt, user.UserID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete User deletes a user from the database
func (cr *UserRepo) DeleteUser(id int) (int64, error) {
	// Execute delete query to delete a user from the database
	result, err := cr.DB.Delete("DELETE FROM users WHERE user_id=?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *UserRepo) CreateUser(user *User) error {
	query := `INSERT INTO users (user_name, password, salt, created_at, created_by, status_id) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, user.Username, user.Password, user.Salt, user.CreatedAt, user.CreatedBy, user.StatusID)
	return err
}
