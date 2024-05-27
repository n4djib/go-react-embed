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
