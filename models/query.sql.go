// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package models

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, password) VALUES (?, ?)
RETURNING id, name, is_active, created_at
`

type CreateUserParams struct {
	Name     string `db:"name" json:"name" validate:"required"`
	Password string `db:"password" json:"password" validate:"required"`
}

type CreateUserRow struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name" validate:"required"`
	IsActive  *bool      `db:"is_active" json:"is_active"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Name, arg.Password)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsActive,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getPokemon = `-- name: GetPokemon :one

SELECT id, name, image FROM pokemons WHERE id = ? LIMIT 1
`

// ---------------------------------------
// ---------------------------------------
func (q *Queries) GetPokemon(ctx context.Context, id int64) (Pokemon, error) {
	row := q.db.QueryRowContext(ctx, getPokemon, id)
	var i Pokemon
	err := row.Scan(&i.ID, &i.Name, &i.Image)
	return i, err
}

const getPokemonByName = `-- name: GetPokemonByName :one
SELECT id, name, image FROM pokemons WHERE name = ? LIMIT 1
`

func (q *Queries) GetPokemonByName(ctx context.Context, name string) (Pokemon, error) {
	row := q.db.QueryRowContext(ctx, getPokemonByName, name)
	var i Pokemon
	err := row.Scan(&i.ID, &i.Name, &i.Image)
	return i, err
}

const getPokemonWithPassword = `-- name: GetPokemonWithPassword :one
SELECT id, name, image FROM pokemons WHERE id = ? LIMIT 1
`

func (q *Queries) GetPokemonWithPassword(ctx context.Context, id int64) (Pokemon, error) {
	row := q.db.QueryRowContext(ctx, getPokemonWithPassword, id)
	var i Pokemon
	err := row.Scan(&i.ID, &i.Name, &i.Image)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, name, is_active, created_at FROM users WHERE id = ? LIMIT 1
`

type GetUserRow struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name" validate:"required"`
	IsActive  *bool      `db:"is_active" json:"is_active"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

func (q *Queries) GetUser(ctx context.Context, id int64) (GetUserRow, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i GetUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsActive,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT id, name, is_active, created_at FROM users WHERE name = ? LIMIT 1
`

type GetUserByNameRow struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name" validate:"required"`
	IsActive  *bool      `db:"is_active" json:"is_active"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

func (q *Queries) GetUserByName(ctx context.Context, name string) (GetUserByNameRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, name)
	var i GetUserByNameRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsActive,
		&i.CreatedAt,
	)
	return i, err
}

const getUserWithPassword = `-- name: GetUserWithPassword :one
SELECT id, name, password, is_active, created_at FROM users WHERE id = ? LIMIT 1
`

func (q *Queries) GetUserWithPassword(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserWithPassword, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Password,
		&i.IsActive,
		&i.CreatedAt,
	)
	return i, err
}

const listPokemons = `-- name: ListPokemons :many
SELECT id, name, image FROM pokemons ORDER BY id
`

func (q *Queries) ListPokemons(ctx context.Context) ([]Pokemon, error) {
	rows, err := q.db.QueryContext(ctx, listPokemons)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pokemon
	for rows.Next() {
		var i Pokemon
		if err := rows.Scan(&i.ID, &i.Name, &i.Image); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, is_active, created_at FROM users ORDER BY id
`

type ListUsersRow struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name" validate:"required"`
	IsActive  *bool      `db:"is_active" json:"is_active"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUsersRow
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsActive,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users set name = ?, password = ?, is_active = ? WHERE id = ?
RETURNING id, name, is_active, created_at
`

type UpdateUserParams struct {
	Name     string `db:"name" json:"name" validate:"required"`
	Password string `db:"password" json:"password" validate:"required"`
	IsActive *bool  `db:"is_active" json:"is_active"`
	ID       int64  `db:"id" json:"id"`
}

type UpdateUserRow struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name" validate:"required"`
	IsActive  *bool      `db:"is_active" json:"is_active"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Name,
		arg.Password,
		arg.IsActive,
		arg.ID,
	)
	var i UpdateUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsActive,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserActiveState = `-- name: UpdateUserActiveState :one
UPDATE users set is_active = ? WHERE id = ?
RETURNING id, name, is_active, created_at
`

type UpdateUserActiveStateParams struct {
	IsActive *bool `db:"is_active" json:"is_active"`
	ID       int64 `db:"id" json:"id"`
}

type UpdateUserActiveStateRow struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name" validate:"required"`
	IsActive  *bool      `db:"is_active" json:"is_active"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

func (q *Queries) UpdateUserActiveState(ctx context.Context, arg UpdateUserActiveStateParams) (UpdateUserActiveStateRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserActiveState, arg.IsActive, arg.ID)
	var i UpdateUserActiveStateRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsActive,
		&i.CreatedAt,
	)
	return i, err
}
