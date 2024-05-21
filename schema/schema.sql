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

REPLACE INTO pokemons (image, name, id) VALUES 
 ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/21.svg', 'spearow', 1),
 ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/22.svg', 'fearow', 2),
 ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/1.svg', 'bulbasaur', 3),
 ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/2.svg', 'ivysaur', 4),
 ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/4.svg', 'charmander', 5),
 ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/24.svg', 'arbok', 6),
 ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/25.svg', 'pikachu', 7),
 ('https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/27.svg', 'sandshrew', 8);

