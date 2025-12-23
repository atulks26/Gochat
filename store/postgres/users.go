package postgres

import (
	"chat/store/users"
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) FindClientByUsername(username string) (int64, bool, error) {
	var id int64
	query := `SELECT id FROM users WHERE username = $1`

	err := s.db.QueryRow(query, username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}

	return id, true, nil
}

func (s *Store) FindClientByID(id int64) (*users.UserData, error) {
	var u users.UserData
	query := `SELECT id, username, password_hash FROM users WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.HashedPassword)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Store) CreateUser(username, password string) (*users.UserData, error) {
	var u users.UserData
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, password_hash`

	err := s.db.QueryRow(query, username, password).Scan(&u.ID, &u.Username, &u.HashedPassword)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
