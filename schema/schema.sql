CREATE TABLE IF NOT EXISTS users(
  -- id INTEGER PRIMARY KEY AUTOINCREMENT,
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  is_active  BOOLEAN DEFAULT (false),
  created_at DATETIME DEFAULT (CURRENT_TIMESTAMP) 
);

REPLACE INTO users (id, name, password) values (1, "n4djib", "long one");
REPLACE INTO users (id, name, password) values (2, "nad", "ddddd");

CREATE TABLE IF NOT EXISTS pokemons (
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    name  TEXT UNIQUE   NOT NULL,
    image TEXT    NOT NULL
);

REPLACE INTO pokemons (id, name, image) VALUES 
  (1, 'spearow', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/21.svg'),
  (2, 'fearow', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/22.svg');

