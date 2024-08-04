# Golang PostgreSQL
Este repositório didático explora as possibilidades existentes em uma aplicação Golang integrada a um banco de dados relacional PostgreSQL, intermediada pelo driver ou lib **database/sql** do próprio Go.
O pacote uma interface genérica para trabalhar com bancos de dados SQL. Ele permite que você execute operações de banco de dados de forma consistente, independentemente do banco de dados específico que você está usando (como PostgreSQL, MySQL, SQLite, etc.). A biblioteca abstrai os detalhes específicos do banco de dados e oferece um conjunto de funções e tipos para executar consultas, gerenciar transações e manipular conexões.

<details>
    <summary>Criação de Client (Conexão).</summary>

## Client
Um client ou uma conexão refere-se a uma instância que gerencia a comunicação com o banco de dados e a execução de operações SQL.
Na abordagem utilizado criamos um Singleton do client, onde será centralizado apenas uma conexão com o banco de dados, e a partir desta conexão poderemos executar diversas operações.
```go
package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	host     string
	port     int
	user     string
	password string
	database string
}

type IPostgres interface {
	Connect()
}

var Client *sql.DB

func (p *Postgres) Connect() {

	credentials := Postgres{
		host:     "localhost",
		port:     5432,
		user:     "postgres",
		password: "1234",
		database: "teste",
	}

	connStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", credentials.host, credentials.port, credentials.user, credentials.password, credentials.database)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("POSTGRES CONNECTION SUCCESS!")

	Client = db
}
```
</details>

<details>
    <summary>Entidade</summary>

## Users
A nossa aplicação realiza operações com base na tabela Users no banco de dados, que tem suas colunas equivalentes a nossa struct User.
```go
package entities

import "time"

type User struct {
	Id        int
	Name      string
	Username  string
	Email     string
	Password  string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
```
</details>

<details>
    <summary>CRUD (Create, Read, Update, Delete).</summary>

## Leituras com Query e QueryRow
A leitura de linhas em tabelas do banco pode ser abordada de duas formas, leituras de apenas uma linha e de mais de uma, o pacote sql permite efetuar duas operações como o método Query para trazer diversas linhas e o QueryRow para apenas uma.

### Query
O método query trás como retorno as rows ou linhas de resultado da consulta e um error, caso não tenha ocorrido nenhum erro o valor será **nil**.
Em caso de sucesso, será necessário iterar as linhas retornadas e converter os valores de cada coluna da linha nos valores da nossa struct User.
```go
func FindUsers() ([]entities.User, error) {
    users := []entities.User{}

	rows, err := postgres.Client.Query("SELECT * FROM <nome do banco>.users")
    defer rows.Close()

	if err != nil {
		fmt.Println("ERROR\n", err)
	} else {
        for rows.Next() {
            user := entities.User{}
		    err = rows.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

            if err != nil {
                fmt.Println("SCAN ERROR\n", err)
            }

            users = append(users, user)
        }
    }

	return users, err
}
```
E se você precisar adicionar filtros a sua consulta?
O método Query aceita argumentos como parametros, sendo que no seu where voce deve substituir o valor real por um cifrão e o indice do parametro, por exemplo:

```go
id := 1
rows, err := postgres.Client.Query("SELECT * FROM <nome do banco>.users WHERE id = $1", id)
```
```go
ids := []int{1,2}
rows, err := postgres.Client.Query("SELECT * FROM <nome do banco>.users WHERE id IN ($1,$2)", ids...)
```

### QueryRow
O método **QueryRow** trás como retorno a row ou a linha de resultado da consulta.
E como na operação anterior será necessário iterar a linha retornada e converter os valores de cada coluna nos valores da nossa struct User.
QueryRow também aceita argumentas como a consulta anterior.
```go
func FindUserById(id int) (entities.User, error) {
	row := postgres.Client.QueryRow("SELECT * FROM <nome do banco>.users WHERE id = $1", id)
	
    user := entities.User{}

    err = row.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
    if err != nil {
        fmt.Println("SCAN ERROR\n", err)
    }

	return user, err
}
```

## Escrevendo ou criando, atualizando e removendo com Exec
A escrita de linhas em tabelas do banco pode ser abordada de duas formas, inserção, e update, o pacote sql permite efetuar duas operações como o método Query para trazer diversas linhas e o QueryRow para apenas uma.

### Exec
O método **Exec** trás como retorno o result que contém a informação de quantas linhas foram afetadas e um error, caso não tenha ocorrido nenhum erro o valor será **nil**.

### Insert
```go
func CreateUser(user entities.User) error {
	
	result, err := postgres.Client.Exec("INSERT INTO <nome do banco>.users (name, username, email, password, createdat) VALUES ($1, $2, $3, $4, $5)", user.Name, user.Username, user.Email, user.Password, time.Now())
    if err != nil {
		fmt.Println("INSERT ERROR\n", err)
		return err
	} else {
		rows, _ := result.RowsAffected()
		fmt.Println("ROWS AFFECTED: ", rows)
	}

    return nil
}
```

### Update
```go
func UpdateUser(id int, name string, email string) error {
	
	result, err := postgres.Client.Exec("UPDATE <nome do banco de dados>.users SET name = $1, email = $2 WHERE id = $3", name, email, id)
    if err != nil {
		fmt.Println("UPDATE ERROR\n", err)
		return err
	} else {
		rows, _ := result.RowsAffected()
		fmt.Println("ROWS AFFECTED: ", rows)
	}

    return nil
}
```

### Delete
```go
func DeleteUser(id int) error {
	
	result, err := postgres.Client.Exec("DELETE FROM <nome do banco de dados>.users WHERE id = $1", id)
    if err != nil {
		fmt.Println("DELETE ERROR\n", err)
		return err
	} else {
		rows, _ := result.RowsAffected()
		fmt.Println("ROWS AFFECTED: ", rows)
	}

    return nil
}
```

### Escrevendo com QueryRow
E se eu quiser escrever e ter o retorno da linha que eu escrevi?

Voce pode usar o QueryRow e adicionar ao final de sua operação o comando "RETURNING *" para retornar a linha

Isso servirá para operações de **INSERT** e **UPDATE**.
```go
func CreateUser(user entities.User) (entities.User, error) {

	row := postgres.Client.QueryRow("INSERT INTO teste.users (name, username, email, password, createdat) VALUES ($1, $2, $3, $4, $5) RETURNING *", user.Name, user.Username, user.Email, user.Password, time.Now())

	newUser := entities.User{}

	err := row.Scan(&newUser.Id, &newUser.Name, &newUser.Username, &newUser.Email, &newUser.Password, &newUser.Token, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.DeletedAt)

	if err != nil {
		fmt.Println("SCAN ERROR\n", err)
	}

	return newUser, err
}
```
</details>

<details>
    <summary>Prepared Statements</summary>

## Prepared Statements
Consultas preparadas, ou prepared statements, são uma forma eficiente e segura de executar consultas repetidas em um banco de dados. Ao usar consultas preparadas, você pode separar a compilação da consulta SQL da execução dos dados, o que pode melhorar o desempenho e a segurança.

As consultas preparadas são compiladas uma vez pelo servidor de banco de dados e podem ser executadas várias vezes com diferentes parâmetros sem precisar ser recompiladas.

Usar consultas preparadas ajuda a prevenir injeções de SQL, já que os parâmetros são tratados separadamente do comando SQL.

Segue um exemplo onde para criar um usuário antes é feito uma busca de usuários pelo nome, caso não exista deverá ser criado, então é feito a busca novamente para retornar o valor do usuários criado, senão é retornado um usuário vazio e uma mensagem de erro.
```go
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
```
Neste exemplo, a consulta é reutilizada, os argumentos passados tem os mesmos valores porém, se fossem diferentes, a consulta seria executada novamente com sucesso.
A execução da query de um Prepared Statement pode ser tanto um **Query**, **QueryRow**, ou **Exec**.
</details>

<details>
    <summary>Transactions</summary>

## Transactions
Transações são um conceito fundamental em bancos de dados, incluindo PostgreSQL, que permitem agrupar uma ou mais operações SQL em uma única unidade de trabalho. As transações garantem que todas as operações dentro da transação sejam concluídas com sucesso ou nenhuma delas seja aplicada, mantendo a consistência e integridade dos dados.

As transações seguem as propriedades **ACID**:

**Atomicidade:** Todas as operações dentro da transação são completadas com sucesso ou nenhuma delas é aplicada.

**Consistência:** As transações levam o banco de dados de um estado consistente a outro estado consistente.

**Isolamento:** As transações são isoladas umas das outras, garantindo que os resultados de uma transação não sejam visíveis para outras transações até que sejam finalizadas.

**Durabilidade:** Uma vez que uma transação é confirmada (committed), seus efeitos persistem no banco de dados mesmo que haja uma falha no sistema.

**Locks:** Durante a utilização de transactions é necessário ter cuidado, pois elas geram locks e podem comprometer a performance das operações no seu banco de dados, bloqueando linhas ou até mesmo tabelas inteiras de serem lidas ou escritas até que sua execução termine.

Para usar uma transação em uma ou mais operações você deve xriar a transação e a partir dela executar suas querys, caso uma das operações de erro você deve fazer o rollback desta operação, ou seja, se voce estiver inserindo, atualizando ou deletando alguma linha, esta operação não será efetuada e a sua tabela permanecerá no estado anterior as operações, e em caso de sucesso, você deve efetuar o commit da transação para que as operações sejam efetuadas, assim alterando seu banco de dados definitivamente.
```go
func CreateUser(user entities.User) error {
	tx, err := postgres.Client.Begin()
	if err != nil {
		fmt.Println("TRANSACTION INTANCE ERROR\n", err)
	}

	defer tx.Rollback()
	defer tx.Commit()

	result, err := tx.Exec("INSERT INTO <nome do banco>.users (name, username, email, password, createdat) VALUES ($1, $2, $3, $4, $5)", user.Name, user.Username, user.Email, user.Password, time.Now())
    if err != nil {
		fmt.Println("INSERT ERROR\n", err)
		return err
	} else {
		rows, _ := result.RowsAffected()
		fmt.Println("ROWS AFFECTED: ", rows)
	}

    return nil
}
```

Também é possível utilizar Prepared Statements com transações, por exemplo:
```go
func CreateUser(user entities.User) error {
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
```
</details>

Neste repositórios eu exercitei algumas ideias e maneiras de utilizar todos estes recursos, fique a vontade para explorar.

Pontos que ainda podem ser aprofundados seriam os tipos de Lock e Isolamentos das transações, pretendo abordar em breve.

