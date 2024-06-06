-- name: GetUser :one
SELECT id, name, is_active, session, logged_at, created_at
  FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByName :one
SELECT id, name, is_active, session, logged_at, created_at
  FROM users WHERE name = ? LIMIT 1;

-- name: GetUserBySession :one
SELECT id, name, is_active, session, logged_at, created_at 
  FROM users WHERE session = ? and session <> "" and session is not null  LIMIT 1;

-- name: GetUserWithPassword :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByNameWithPassword :one
SELECT * FROM users WHERE name = ? LIMIT 1;

-- name: ListUsers :many
SELECT id, name, is_active, session, logged_at, created_at FROM users ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (name, password, created_at) VALUES (?, ?, ?)
RETURNING id, name, is_active, session, logged_at, created_at;

-- name: UpdateUser :one
UPDATE users set name = ?, password = ?, is_active = ? WHERE id = ?
RETURNING id, name, is_active, session, logged_at, created_at;

-- name: UpdateUserSession :exec
UPDATE users set session = ?, logged_at = ? WHERE id = ?;

-- name: UpdateUserActiveState :one
UPDATE users set is_active = ? WHERE id = ?
RETURNING id, name, is_active, session, logged_at, created_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-----------------------------------------
-----------------------------------------

-- name: GetUserRoles :many
select role from roles join user_roles on roles.id == user_roles.role_id 
 where user_roles.user_id == ?;


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

-- name: ListPokemonsOffset :many
SELECT * FROM pokemons ORDER BY id LIMIT ? OFFSET ?;

-- name: ListPokemonsNames :many
SELECT name FROM pokemons

