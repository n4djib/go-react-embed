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
        // DB = nil
        fmt.Println("Error opening database:", err)
		return err
	}
	DB = db
	return nil
}

func CloseDatabase() {
    DB.Close()
}

func CreateDbTables(file string) error {
    // Open a connection to the SQLite database
    err := ConnectDatabase(file)
	if err != nil {
        return err
    }

    // Create a new table
    createUsersTableSQL := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT unique NOT NULL,
            password TEXT NOT NULL
        );
    `
    _, err = DB.Exec(createUsersTableSQL)
	if err != nil {
        return err
    }

    createPokemonsTableSQL := `
        CREATE TABLE IF NOT EXISTS pokemons (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            image TEXT NOT NULL
        );
    `
    _, err = DB.Exec(createPokemonsTableSQL)
    return err
}

// TODO open and close the db
func GetPokemons() ([]Pokemon, error) {
    rows, err := DB.Query("SELECT id, name, image FROM pokemons")
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
    defer stmt.Close()

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

///////

func GetUsers() ([]User, error) {
    rows, err := DB.Query("SELECT id, name, password FROM users")
    if err != nil {
        return nil, err
    }

    users := []User{}
    for rows.Next() {
        user := User{}
        err = rows.Scan(&user.Id, &user.Name, &user.Password)
        if err != nil {
			return nil, err
		}
        users = append(users, user)
    }
    err = rows.Err()

	if err != nil {
		return nil, err
	}
	return users, err
}

func GetUser(name string) (User, error) {
    stmt, err := DB.Prepare("SELECT id, name, image FROM users WHERE name=?")
    if err != nil {
        return User{}, err
    }
    defer stmt.Close()

    user := User{}
    sqlErr := stmt.QueryRow(name).Scan(&user.Id, &user.Name, &user.Password)
    if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, sqlErr
	}
	return user, nil
}

func CreateUser(user User) (bool, error) {
    trans, err := DB.Begin()
    if err != nil {
        return false, err
    }

    stmt, err := trans.Prepare("INSERT INTO users (name, password) VALUES (?, ?)")
    if err != nil {
		return false, err
	}
    defer stmt.Close()

    _, err = stmt.Exec(user.Name, user.Password)
    if err != nil {
		return false, err
	}

    trans.Commit()
    return true, nil
}