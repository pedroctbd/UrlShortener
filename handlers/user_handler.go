package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"shorturl.com/entities"
	"shorturl.com/utils"
)

// CreateUser godoc
// @Summary      Create a new user
// @Description  Creates a new user with email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      entities.CreateUserInput  true  "User input"
// @Success      201   {object}  map[string]string          "ID of the created user"
// @Router       /user [post]
func CreateUser(db *sql.DB) http.HandlerFunc {

	query := `INSERT INTO users (id, email, password_hash, created_at, updated_at)
	VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
	RETURNING id;`

	fmt.Println("ENTERED CREATE USER HANDLER")

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		input, err := utils.Decode[entities.CreateUserInput](r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}

		var newUserId string
		if err := db.QueryRowContext(ctx, query, input.Email, string(hashedPassword)).Scan(&newUserId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		resp := map[string]interface{}{
			"id": newUserId,
		}

		if err := utils.Encode(w, r, http.StatusCreated, resp); err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

// ListUsers godoc
// @Summary      List all users
// @Description  Retrieves a list of all registered users
// @Tags         users
// @Produce      json
// @Success      200  {array}   entities.User
// @Router       /users [get]
func ListUsers(db *sql.DB) http.HandlerFunc {
	query := `SELECT * FROM users`
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		rows, err := db.QueryContext(ctx, query)
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users := []entities.User{}

		for rows.Next() {
			var user entities.User

			err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			users = append(users, user)
		}

		if err := utils.Encode(w, r, http.StatusOK, users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}
