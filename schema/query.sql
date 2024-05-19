-- name: GetUser :one
SELECT id, name, is_active, created_at FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByName :one
SELECT id, name, is_active, created_at FROM users WHERE name = ? LIMIT 1;

-- name: GetUserWithPassword :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT id, name, is_active, created_at FROM users ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (name, password) VALUES (?, ?)
RETURNING id, name, is_active, created_at;

-- name: UpdateUser :one
UPDATE users set name = ?, password = ?, is_active = ? WHERE id = ?
RETURNING id, name, is_active, created_at;

-- name: UpdateUserActiveState :one
UPDATE users set is_active = ? WHERE id = ?
RETURNING id, name, is_active, created_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-----------------------------------------
-----------------------------------------

-- name: GetPokemon :one
SELECT * FROM pokemons WHERE id = ? LIMIT 1;

-- name: GetPokemonByName :one
SELECT * FROM pokemons WHERE name = ? LIMIT 1;

-- name: GetPokemonWithPassword :one
SELECT * FROM pokemons WHERE id = ? LIMIT 1;

-- name: ListPokemons :many
SELECT * FROM pokemons ORDER BY id;


