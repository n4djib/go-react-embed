-----------------------------------------------------------------------------
-- Modify the tables directly in the database and reflect it in the schema --
-----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS users(
  -- id         INTEGER PRIMARY KEY AUTOINCREMENT,
  id         INTEGER PRIMARY KEY,
  name       TEXT UNIQUE NOT NULL,
  password   TEXT NOT NULL,
  is_active  BOOLEAN DEFAULT (false),
  session    TEXT,
  logged_at  DATETIME,
  created_at DATETIME
);

INSERT OR IGNORE INTO users (id, name, password, is_active) VALUES 
  (1, "n4djib", "$2a$10$D37JXHtApnRcfq5S77im/OL4/f0GHwDEEMuZGlJbtjX.a15aUx8r6", 0),
  (2, "nad", "$2a$10$D37JXHtApnRcfq5S77im/OL4/f0GHwDEEMuZGlJbtjX.a15aUx8r6", 1);


CREATE TABLE IF NOT EXISTS roles(
  id          INTEGER PRIMARY KEY,
  role        TEXT UNIQUE NOT NULL,
  description TEXT
);

INSERT OR IGNORE INTO roles (id, role) VALUES
  (1, "USER"), (2, "ADMIN"), (3, "MANAGER");


CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    role_id INTEGER REFERENCES roles (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

INSERT OR IGNORE INTO user_roles (user_id, role_id) VALUES
  (3, 2), (3, 3), (2, 1);


CREATE TABLE IF NOT EXISTS pokemons (
  id    INTEGER PRIMARY KEY AUTOINCREMENT,
  name  TEXT UNIQUE NOT NULL,
  image TEXT NOT NULL
);

INSERT OR IGNORE INTO pokemons (image, name, id) VALUES 
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/21.svg', 'spearow', 1),
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/22.svg', 'fearow', 2),
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/1.svg', 'bulbasaur', 3),
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/2.svg', 'ivysaur', 4),
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/4.svg', 'charmander', 5),
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/24.svg', 'arbok', 6),
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/25.svg', 'pikachu', 7),
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/27.svg', 'sandshrew', 8),
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/37.svg', 'vulpix', 9),
  ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/38.svg', 'ninetales', 10);



CREATE TABLE IF NOT EXISTS permissions (
    id          INTEGER PRIMARY KEY,
    permission  TEXT    UNIQUE NOT NULL,
    description TEXT,
    rule        TEXT
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       INTEGER REFERENCES roles (id) NOT NULL,
    permission_id INTEGER REFERENCES permissions (id) NOT NULL,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS permission_child (
    permission_id       INTEGER REFERENCES permissions (id) NOT NULL,
    child_permission_id INTEGER REFERENCES permissions (id) NOT NULL,
    PRIMARY KEY (permission_id, child_permission_id)
);

CREATE TABLE IF NOT EXISTS role_child (
    role_id       INTEGER REFERENCES roles (id) NOT NULL,
    child_role_id INTEGER REFERENCES roles (id) NOT NULL,
    PRIMARY KEY (role_id, child_role_id)
);
