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
SELECT name FROM pokemons;

-----------------------------------------
-----------------------------------------
-- RBAC

-- name: GetRolesParents :many
SELECT r.id, r.role, rc.role_id AS parent 
  FROM roles r 
  LEFT JOIN role_child rc ON r.id = rc.child_role_id;

-- name: GetPermissionsParents :many
SELECT p.id, p.permission, p.rule, pc.permission_id AS parent
    FROM permissions p 
    LEFT JOIN permission_child pc ON p.id = pc.child_permission_id;

-- name: GetPermissionsRoles :many
SELECT p.id, p.permission, rp.role_id AS role_id, r.role 
    FROM permissions p
    JOIN role_permissions rp ON p.id = rp.permission_id
    JOIN roles r ON rp.role_id = r.id;

-- name: GetUsersRoles :many
SELECT u.id, u.name, r.id AS role_id, r.role 
    FROM users u  
    JOIN user_roles ur ON u.id = ur.user_id
    JOIN roles r ON r.id = ur.role_id;
