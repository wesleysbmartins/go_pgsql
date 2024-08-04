package repository

import (
	"fmt"
	"go_pgsql/internal/adapters/database/sql/operations"
	"go_pgsql/internal/entities"
	"time"
)

type UserRepository struct{}

type IUserRepository interface {
	Create(user entities.User) entities.User
	FindAll(params *UserParams) []entities.User
	FindOne(params UserParams) entities.User
	Update(user UserParams) entities.User
}

func (r *UserRepository) Create(user entities.User) entities.User {
	sql := operations.SqlOperations{}
	newUser := entities.User{}

	row := sql.Insert(user)

	err := row.Scan(&newUser.Id, &newUser.Name, &newUser.Username, &newUser.Email, &newUser.Password, &newUser.Token, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.DeletedAt)

	if err != nil {
		fmt.Println("ERROR\n", err)
	}

	fmt.Println("RESULT: ", newUser)

	return newUser
}

func (r *UserRepository) FindAll(params *UserParams) []entities.User {
	sql := operations.SqlOperations{}
	where := params.handleWhere()

	rows := sql.Select(entities.User{}, where)

	defer rows.Close()

	users := []entities.User{}

	for rows.Next() {
		user := entities.User{}

		err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

		if err != nil {
			fmt.Println("SCAN ERR: ", err)
		}

		users = append(users, user)
	}

	err := rows.Err()

	if err != nil {
		panic(err)
	}

	fmt.Println("SELECT RESULT: ", users)

	return users
}

func (r *UserRepository) FindOne(params UserParams) []entities.User {
	sql := operations.SqlOperations{}
	where := params.handleWhere()

	row := sql.SelectOne(entities.User{}, where)

	users := []entities.User{}

	user := entities.User{}

	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

	if err != nil {
		fmt.Println("SCAN ERR: ", err)
	}

	users = append(users, user)

	fmt.Println("SELECT RESULT: ", users)

	return users
}

func (r *UserRepository) Update(params UserParams) error {
	if params.Id == 0 {
		return fmt.Errorf("Param id is required to Update!")
	}

	where := operations.Where{Equals: &[]operations.Equal{{Column: "id", Value: params.Id}}}

	set := params.handleSet()
	set.Values = append(set.Values, operations.SetParam{Column: "updatedAt", Value: time.Now()})

	sql := operations.SqlOperations{}

	result, err := sql.Update(entities.User{}, set, where)

	if err != nil {
		errmsg := fmt.Sprintf("UPDATE ERROR\n%s", err)
		fmt.Println(errmsg)
		return fmt.Errorf(errmsg)
	} else {
		rows, err := result.RowsAffected()
		fmt.Println("ROWS AFFECTED: ", rows, "ERROR\n", err)
	}

	return err
}

func (r *UserRepository) SoftDelete(id int) error {
	if id == 0 {
		return fmt.Errorf("Param id is required to Update!")
	}

	where := operations.Where{Equals: &[]operations.Equal{{Column: "id", Value: id}}}
	set := operations.Set{Values: []operations.SetParam{{Column: "deletedAt", Value: time.Now()}}}

	sql := operations.SqlOperations{}

	result, err := sql.Update(entities.User{}, set, where)

	if err != nil {
		errmsg := fmt.Sprintf("SOFT DELETE ERROR\n%s", err)
		fmt.Println(errmsg)
		return fmt.Errorf(errmsg)
	} else {
		rows, err := result.RowsAffected()
		fmt.Println("ROWS AFFECTED: ", rows, "ERROR\n", err)
	}

	return err
}

func (r *UserRepository) Delete(id int) error {
	if id == 0 {
		return fmt.Errorf("Param id is required to Update!")
	}

	sql := operations.SqlOperations{}

	result, err := sql.Delete(entities.User{}, id)

	if err != nil {
		errmsg := fmt.Sprintf("DELETE ERROR\n%s", err)
		fmt.Println(errmsg)
		return fmt.Errorf(errmsg)
	} else {
		rows, err := result.RowsAffected()
		fmt.Println("ROWS AFFECTED: ", rows, "ERROR\n", err)
	}

	return err
}
