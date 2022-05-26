// Package db пакет для работы с базами данных
package database

import (
	"context"
	myerrors "gophkeeper/internal/common/errors"
	"gophkeeper/internal/common/models"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// NewDB функция по созданию нового объекта для работы с базой данных
func NewDB(conn *sqlx.DB) *PGDB {
	return &PGDB{
		conn: conn,
	}
}

// PGDB структура данных для подключения к базе данных
type PGDB struct {
	conn *sqlx.DB
}

// CreateUser функция для создание пользователя по модели пользователя
func (pg *PGDB) CreateUser(ctx context.Context, user models.User) error {
	_, err := pg.conn.ExecContext(ctx, "INSERT INTO users (login, password) VALUES ($1, crypt($2, gen_salt('bf', 8)))",
		user.Login, user.Password)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == pgerrcode.UniqueViolation {
			return myerrors.NewError(err, "user already exists")
		}
	}
	return err
}

// CheckUserPassword функция проверки соответствия пароля, возвращает так же id пользователя
func (pg *PGDB) CheckUserPassword(ctx context.Context, user models.User) (string, error) {
	var result string
	query, err := pg.conn.QueryxContext(ctx, `SELECT id FROM users WHERE login = $1 
                       AND password = crypt($2, password)`,
		user.Login, user.Password)
	if err != nil {
		return result, myerrors.NewError(err, "can't find user")
	}
	query.Next()
	err = query.Scan(&result)
	if err != nil {
		return result, myerrors.NewError(err, "error with getting user_id")
	}
	return result, err
}

// AddSecret функция дл добавления секрета пользователя
func (pg *PGDB) AddNote(ctx context.Context, userID string, note models.Note) (string, error) {
	id := ""
	result, err := pg.conn.QueryxContext(ctx,
		"INSERT INTO user_notes (title, body, user_id) VALUES ($1, $2, $3) RETURNING id",
		note.Title, note.Body, userID)
	if err != nil {
		return id, err
	}
	for result.Next() {
		err = result.Scan(&id)
	}
	return id, err
}

// GetSecrets функция для получение всех секретов пользователя
func (pg *PGDB) GetNotes(ctx context.Context, userID string) ([]models.Note, error) {
	rows, err := pg.conn.QueryxContext(ctx, "SELECT id, user_id,secret_data FROM secrets WHERE user_id=$1 AND deleted_at IS NULL", userID)
	if err != nil {
		return nil, err
	}
	var result []models.Note
	for rows.Next() {
		m := models.Note{}
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, err
}

func (pg *PGDB) GetUpdates(ctx context.Context, userID string, after time.Time) ([]models.Note, error) {
	rows, err := pg.conn.QueryxContext(ctx, `SELECT id, title, body, created_at, updated_at, deleted_at FROM user_notes WHERE user_id = $1 
        AND ( created_at > $2 OR updated_at > $2 OR deleted_at > $2)`, userID, after)
	if err != nil {
		return nil, err
	}
	notes := make([]models.Note, 0)
	for rows.Next() {
		var note models.Note
		var deletedAt *time.Time
		err := rows.Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt, &note.UpdatedAt, &deletedAt)
		if deletedAt != nil {
			note.Deleted = true
			note.DeletedAt = *deletedAt
		}
		if err != nil {
			return notes, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// UpdateNote функция для удаления секрета пользователя
func (pg *PGDB) UpdateNote(ctx context.Context, userID string, note models.Note) error {
	_, err := pg.conn.ExecContext(ctx, `UPDATE user_notes SET title=$1, body=$2, updated_at = current_timestamp WHERE id = $3 AND user_id = $4 AND updated_at < $5`,
		note.Title, note.Body, note.ID, userID, note.UpdatedAt)
	return err
}

// DeleteNote функция для удаления секрета пользователя
func (pg *PGDB) DeleteNote(ctx context.Context, userID string, id string) error {
	_, err := pg.conn.ExecContext(ctx, `UPDATE user_notes SET deleted_at = current_timestamp WHERE id = $1 AND user_id = $2`,
		id, userID)
	return err
}
