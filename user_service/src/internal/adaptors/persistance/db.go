package persistance

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"user_service/internal/core/users"

	_ "github.com/lib/pq"
)

type DBStore struct {
	DB *sql.DB
}

func NewDBStore(dbSource string) (*DBStore, error) {
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}
	log.Println("Database connection successful.")
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("could not run migrations: %w", err)
	}
	return &DBStore{DB: db}, nil
}

// to execute sql file.
func runMigrations(db *sql.DB) error {
	migration, err := ioutil.ReadFile("migrations/init.sql")
	if err != nil {
		return fmt.Errorf("could not read migration file: %w", err)
	}
	_, err = db.Exec(string(migration))
	if err != nil {
		return fmt.Errorf("could not execute migration: %w", err)
	}
	log.Println("Database migration successful.")
	return nil
}

func (store *DBStore) CreateUser(user *users.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	err := store.DB.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	return nil
}

func (store *DBStore) GetUserByEmail(email string) (*users.User, error) {
	user := &users.User{}
	query := `SELECT id, username, email, password_hash FROM users WHERE email = $1`
	err := store.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("could not get user by email: %w", err)
	}
	return user, nil
}

func (store *DBStore) GetUserByID(id int) (*users.User, error) {
	user := &users.User{}
	query := `SELECT id, username, email FROM users WHERE id = $1`
	err := store.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("could not get user by id: %w", err)
	}
	return user, nil
}
