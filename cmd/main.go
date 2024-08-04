package main

import (
	"fmt"
	"go_pgsql/internal/adapters/database/sql/postgres"
	"go_pgsql/internal/entities"
	"time"
)

func init() {
	postgres := postgres.Postgres{}
	postgres.Connect()
}

func CreateUserTx(user entities.User) error {
	tx, err := postgres.Client.Begin()
	if err != nil {
		fmt.Println("TRANSACTION INTANCE ERROR\n", err)
	}

	defer tx.Rollback()
	defer tx.Commit()

	stmt, err := tx.Prepare("INSERT INTO <nome do banco>.users (name, username, email, password, createdat) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR\n", err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(user.Name, user.Username, user.Email, user.Password, time.Now())
	if err != nil {
		fmt.Println("INSERT ERROR\n", err)
		return err
	} else {
		rows, _ := result.RowsAffected()
		fmt.Println("ROWS AFFECTED: ", rows)
	}

	return nil
}

func CreateUser(user entities.User) (entities.User, error) {

	row := postgres.Client.QueryRow("INSERT INTO teste.users (name, username, email, password, token, createdat, updatedat) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *", user.Name, user.Username, user.Email, user.Password, user.Token, time.Now(), time.Now())

	newUser := entities.User{}

	err := row.Scan(&newUser.Id, &newUser.Name, &newUser.Username, &newUser.Email, &newUser.Password, &newUser.Token, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.DeletedAt)

	if err != nil {
		fmt.Println("SCAN ERROR\n", err)
	}

	return newUser, err
}

func CreateAndFindUser(user entities.User) (entities.User, error) {
	name := user.Name
	exists := false

	stmt, err := postgres.Client.Prepare("SELECT * FROM teste.users WHERE name = $1")
	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR\n", err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(name)

	if err != nil {
		fmt.Println("STATEMENT READ ERROR\n", err)
	} else {
		defer rows.Close()

		for rows.Next() {
			currentUser := entities.User{}
			err = rows.Scan(&currentUser.Id, &currentUser.Name, &currentUser.Username, &currentUser.Email, &currentUser.Password, &currentUser.Token, &currentUser.CreatedAt, &currentUser.UpdatedAt, &currentUser.DeletedAt)

			if err != nil {
				fmt.Println("SCAN ERROR\n", err)
			}

			if currentUser.Name == name {
				exists = true
			}
		}
	}

	newUser := entities.User{}

	if !exists {
		CreateUser(user)
	} else {
		err := fmt.Errorf("USER EXISTS")
		return newUser, err
	}

	rows, err = stmt.Query(name)

	if err != nil {
		fmt.Println("STATEMENT READ ERROR\n", err)
	} else {
		defer rows.Close()

		for rows.Next() {
			currentUser := entities.User{}
			err = rows.Scan(&currentUser.Id, &currentUser.Name, &currentUser.Username, &currentUser.Email, &currentUser.Password, &currentUser.Token, &currentUser.CreatedAt, &currentUser.UpdatedAt, &currentUser.DeletedAt)

			if err != nil {
				fmt.Println("SCAN ERROR\n", err)
			}

			if currentUser.Name == name {
				newUser = currentUser
			}
		}
	}

	return newUser, err
}

func main() {

	user := entities.User{
		Name:      "TESTE",
		Username:  "TESTE",
		Email:     "TESTE",
		Password:  "TESTE",
		CreatedAt: time.Now(),
	}

	u, err := CreateAndFindUser(user)

	fmt.Println("USERS\n", u, "ERROR\n", err)

	// tx, err := postgres.Client.Begin()
	// if err != nil {
	// 	fmt.Println("TRANSACTION INTANCE ERROR\n", err)
	// }

	// defer tx.Rollback()
	// defer tx.Commit()

	// stmt, _ := postgres.Client.Prepare("SELECT * FROM teste.users")
	// stmt = tx.Stmt(stmt)

	// defer stmt.Close()

	// operations := operations.SqlOperations{}

	// rows, _ := operations.ReadPreparedTx(tx, stmt)

	// defer rows.Close()

	// users := []entities.User{}

	// for rows.Next() {
	// 	user := entities.User{}

	// 	err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

	// 	if err != nil {
	// 		fmt.Println("SCAN ERR: ", err)
	// 	}

	// 	users = append(users, user)
	// }

	// err = rows.Err()

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("SELECT RESULT: ", users)

}
