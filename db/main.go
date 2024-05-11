package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

var DB *sql.DB

func ConnectDatabase(file string) error {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
        DB = nil
        fmt.Println("Error opening database:", err)
		return err
	}
	DB = db
	return nil
}

func CloseDatabase() {
    DB.Close()
}

func CreateDb(file string) error {
    // Open a connection to the SQLite database
    err := ConnectDatabase(file)
	if err != nil {
        return err
    }
    // defer DB.Close() 

    // Create a new table
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        password TEXT NOT NULL
    );
    `
    _, err = DB.Exec(createTableSQL)
	if err != nil {
        return err
    }

    createTableSQL2 := `
        CREATE TABLE IF NOT EXISTS pokemons (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            image TEXT NOT NULL
        );
    `
    _, err = DB.Exec(createTableSQL2)
    return err
}

// TODO open and close the db
func GetPokemons() ([]Pokemon, error) {
    rows, err := DB.Query("SELECT id, name, image from pokemons")
    if err != nil {
        return nil, err
    }

    pokemons := []Pokemon{}
    for rows.Next() {
        pokemon := Pokemon{}
        err = rows.Scan(&pokemon.Id, &pokemon.Name, &pokemon.Image)
        if err != nil {
			return nil, err
		}
        pokemons = append(pokemons, pokemon)
    }
    err = rows.Err()

	if err != nil {
		return nil, err
	}

	return pokemons, err
}

func GetPokemon(id int) (Pokemon, error) {
    stmt, err := DB.Prepare("SELECT id, name, image FROM pokemons WHERE id=?")
    if err != nil {
        return Pokemon{}, err
    }

    pokemon := Pokemon{}
    sqlErr := stmt.QueryRow(id).Scan(&pokemon.Id, &pokemon.Name, &pokemon.Image)
    if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Pokemon{}, nil
		}
		return Pokemon{}, sqlErr
	}

	return pokemon, nil
} 

// func check(e error) {
//     if e != nil {
//         panic(e)
//     }
// }
