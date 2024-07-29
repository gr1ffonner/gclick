package db

import (
	"database/sql"
	"fmt"

	"github.com/gr1ffonner/fintest/internal/config"
	"github.com/gr1ffonner/fintest/pkg/logging"
	"github.com/pkg/errors"
)

type DB struct {
	conn *sql.DB
}

// NewDB creates a new instance of DB with the provided connection string.
func NewDB(cfg config.Config, logger logging.Logger) (*DB, error) {
	psqlconn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db conn db.go")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping db conn db.go")
	}

	return &DB{conn: db}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// Ping pings the database to check the connection.
func (db *DB) Ping() error {
	return db.conn.Ping()
}

func (db *DB) GetMaterialByUUID(id string) (Material, error) {
	var m Material
	getStmt := `SELECT * FROM "material" WHERE "uuid" = $1`
	row := db.conn.QueryRow(getStmt, id)
	err := row.Scan(&m.UUID, &m.MaterialType, &m.PublicationStatus, &m.Title, &m.Content, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return m, errors.Wrap(err, "failed to scan row in GetMaterialByUUID")
	}
	return m, nil
}

func (db *DB) CreateMaterial(m Material) error {
	insertStmt := `
		INSERT INTO "material" ("uuid", "material_type", "publication_status", "title", "content", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.conn.Exec(insertStmt, m.UUID, m.MaterialType, m.PublicationStatus, m.Title, m.Content, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert in CreateMaterial")
	}
	return nil
}

func (db *DB) UpdateMaterial(m Material) error {
	updateStmt := `
		UPDATE "material" SET "material_type" = $2, "publication_status" = $3, "title" = $4, "content" = $5, "updated_at" = $6
		WHERE "uuid" = $1
	`
	_, err := db.conn.Exec(updateStmt, m.UUID, m.MaterialType, m.PublicationStatus, m.Title, m.Content, m.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to execute update in UpdateMaterial")
	}
	return nil
}

func (db *DB) GetAllMaterials() ([]Material, error) {
	var materials []Material
	getAllStmt := `SELECT * FROM "material"`
	rows, err := db.conn.Query(getAllStmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query in GetAllMaterials")
	}
	defer rows.Close()
	for rows.Next() {
		var m Material
		err := rows.Scan(&m.UUID, &m.MaterialType, &m.PublicationStatus, &m.Title, &m.Content, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row in GetAllMaterials")
		}
		materials = append(materials, m)
	}
	return materials, nil
}
