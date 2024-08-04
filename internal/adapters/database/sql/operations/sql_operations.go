package operations

import (
	"database/sql"
	"fmt"
	"go_pgsql/internal/adapters/database/sql/postgres"
	"reflect"
	"strings"
)

type SqlOperations struct{}

type ISqlOperations interface {
	WriteAndReadPreparedTx(query string, args ...interface{}) sql.Row
	WritePreparedTx(query string, args ...interface{}) (sql.Result, error)
	ReadOnPreparedeTX(query string, args ...interface{}) *sql.Row
	ReadPreparedTx(tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error)
	WriteAndReadTx(query string, args ...interface{}) sql.Row
	WriteTx(query string, args ...interface{}) (sql.Result, error)
	ReadOneTX(query string, args ...interface{}) *sql.Row
	ReadTx(tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error)
	WriteAndReadPrepared(query string, args ...interface{}) sql.Row
	WritePrepared(query string, args ...interface{}) (sql.Result, error)
	ReadOnePrepared(query string, args ...interface{}) *sql.Row
	ReadPrepared(query string, args ...interface{}) (*sql.Rows, error)
	WriteAndRead(query string, args ...interface{}) sql.Row
	Write(query string, args ...interface{}) (sql.Result, error)
	ReadOne(query string, args ...interface{}) *sql.Row
	Read(query string, args ...interface{}) (*sql.Rows, error)

	RawQuerySelectTx(tx *sql.Tx, query string, args []interface{}) (rows *sql.Rows, err error)
	RawQuerySelectPrepared(query string, args []interface{}) (rows *sql.Rows, err error)
	RawQueryPrepared(query string, args []interface{}) (rowsAffected int64, err error)
	RawQuerySelect(query string, args []interface{}) (rows *sql.Rows, err error)
	RawQuery(query string, args []interface{}) (rowsAffected int64, err error)
	Insert(entity interface{}) *sql.Row
	Select(entity interface{}, where *Where) *sql.Rows
	SelectOne(entity interface{}, where *Where) *sql.Row
	Update(entity interface{}, where *Where, set interface{})
	Delete(entity interface{}, where *Where)
	getTable(entity interface{}) (tableName string)
	getColumns(entity interface{}) (columns string)
	getValues(entity interface{}) (valuesStr string, values []interface{})
	getEntityReturn(entity interface{}) []interface{}
}

func (o *SqlOperations) WriteAndReadPreparedTx(query string, args ...interface{}) *sql.Row {
	tx, err := postgres.Client.Begin()
	if err != nil {
		fmt.Println("TRANSACTION INTANCE ERROR\n", err)
	}

	defer tx.Rollback()
	defer tx.Commit()

	stmt, err := tx.Prepare(fmt.Sprintf("%s RETURNING *", query))
	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR\n", err)
	}

	defer stmt.Close()

	return stmt.QueryRow(args...)
}

func (o *SqlOperations) WritePreparedTx(query string, args ...interface{}) (sql.Result, error) {
	tx, err := postgres.Client.Begin()
	if err != nil {
		fmt.Println("TRANSACTION INTANCE ERROR\n", err)
	}

	defer tx.Rollback()
	defer tx.Commit()

	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR\n", err)
	}

	defer stmt.Close()

	return stmt.Exec(args...)
}

func (o *SqlOperations) ReadOnePreparedTx(query string, args ...interface{}) *sql.Row {
	tx, err := postgres.Client.Begin()
	if err != nil {
		fmt.Println("TRANSACTION INTANCE ERROR\n", err)
	}

	defer tx.Rollback()
	defer tx.Commit()

	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR\n", err)
	}

	defer stmt.Close()

	rows := stmt.QueryRow(args...)

	return rows
}

func (o *SqlOperations) ReadPreparedTx(tx *sql.Tx, stmt *sql.Stmt, args ...interface{}) (*sql.Rows, error) {

	rows, err := stmt.Query(args...)
	if err != nil {
		fmt.Println("STATEMENT READ ERROR\n", err)
	}

	return rows, err
}

func (o *SqlOperations) WriteAndReadTx(query string, args ...interface{}) *sql.Row {
	tx, err := postgres.Client.Begin()
	if err != nil {
		fmt.Println("TRANSACTION INTANCE ERROR\n", err)
	}

	defer tx.Rollback()
	defer tx.Commit()

	return tx.QueryRow(fmt.Sprintf("%s RETURNING *", query), args...)
}

func (o *SqlOperations) WriteTx(query string, args ...interface{}) (sql.Result, error) {
	tx, err := postgres.Client.Begin()
	if err != nil {
		fmt.Println("TRANSACTION INTANCE ERROR\n", err)
	}

	defer tx.Rollback()
	defer tx.Commit()

	return tx.Exec(query, args...)
}

func (o *SqlOperations) ReadOneTx(query string, args ...interface{}) *sql.Row {
	tx, err := postgres.Client.Begin()
	if err != nil {
		fmt.Println("TRANSACTION INTANCE ERROR\n", err)
	}

	defer tx.Rollback()
	defer tx.Commit()

	return tx.QueryRow(query, args...)
}

func (o *SqlOperations) ReadTx(tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error) {

	rows, err := tx.Query(query, args...)
	if err != nil {
		fmt.Println("READ ERROR\n", err)
	}

	return rows, err
}

func (o *SqlOperations) WriteAndReadPrepared(query string, args ...interface{}) *sql.Row {
	stmt, err := postgres.Client.Prepare(fmt.Sprintf("%s RETURNNG *", query))
	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR\n", err)
	}

	defer stmt.Close()

	return stmt.QueryRow(args...)
}

func (o *SqlOperations) WritePrepared(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := postgres.Client.Prepare(query)
	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR\n", err)
	}

	defer stmt.Close()

	return stmt.Exec(args...)
}

func (o *SqlOperations) ReadOnePrepared(query string, args ...interface{}) *sql.Row {
	stmt, err := postgres.Client.Prepare(query)
	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR\n", err)
	}

	defer stmt.Close()

	rows := stmt.QueryRow(args...)

	return rows
}

func (o *SqlOperations) ReadPrepared(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := postgres.Client.Prepare(query)
	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR\n", err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		fmt.Println("STATEMENT READ ERROR\n", err)
	}

	return rows, err
}

func (o *SqlOperations) WriteAndRead(query string, args ...interface{}) *sql.Row {
	return postgres.Client.QueryRow(fmt.Sprintf("%s RETURNING *", query), args...)
}

func (o *SqlOperations) Write(query string, args ...interface{}) (sql.Result, error) {
	return postgres.Client.Exec(query, args...)
}

func (o *SqlOperations) ReadOne(query string, args ...interface{}) *sql.Row {
	return postgres.Client.QueryRow(query, args...)
}

func (o *SqlOperations) Read(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := postgres.Client.Query(query, args...)
	if err != nil {
		fmt.Println("READ ERROR\n", err)
	}

	return rows, err
}

func (o *SqlOperations) RawQuerySelectTx(tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error) {

	rows, err := tx.Query(query, args...)
	if err != nil {
		fmt.Println("QUERY ERR", err)
	}

	return rows, err
}

func (o *SqlOperations) RawQueryPrepared(query string, args ...interface{}) (int64, error) {

	stmt, err := postgres.Client.Prepare(query)
	if err != nil {
		fmt.Println("PREPARE ERR", err)
	}

	result, err := stmt.Exec(args...)

	stmt.Close()

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (o *SqlOperations) RawQuerySelectPrepared(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := postgres.Client.Prepare(query)
	if err != nil {
		fmt.Println("PREPARE ERR", err)
	}

	rows, err := stmt.Query(args...)

	stmt.Close()

	return rows, err
}

func (o *SqlOperations) RawQuerySelect(query string, args ...interface{}) (*sql.Rows, error) {
	return postgres.Client.Query(query, args...)
}

func (o *SqlOperations) RawQuery(query string, args ...interface{}) (rowsAffected int64, err error) {
	result, err := postgres.Client.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (o *SqlOperations) Insert(entity interface{}) *sql.Row {

	table := o.getTable(entity)
	columns := o.getColumns(entity)
	valuesStr, args := o.getValues(entity)

	query := fmt.Sprintf("INSERT INTO teste.%s (%s) VALUES (%s) RETURNING *;", table, columns, valuesStr)

	row := postgres.Client.QueryRow(query, args...)

	return row
}

func (o *SqlOperations) Select(entity interface{}, where *Where) *sql.Rows {
	table := o.getTable(entity)
	var rows *sql.Rows
	var err error

	if where != nil {
		whereStr, args := where.Handle()

		query := fmt.Sprintf("SELECT * FROM teste.%s %s;", table, whereStr)

		rows, err = postgres.Client.Query(query, args...)
	} else {
		query := fmt.Sprintf("SELECT * FROM teste.%s;", table)
		rows, err = postgres.Client.Query(query)
	}

	if err != nil {
		fmt.Println("SELECT ERROR\n", err)
	}

	return rows
}

func (o *SqlOperations) SelectOne(entity interface{}, where *Where) *sql.Row {
	table := o.getTable(entity)
	whereStr, args := where.Equal.handle()

	query := fmt.Sprintf("SELECT * FROM teste.%s %s;", table, whereStr)
	row := postgres.Client.QueryRow(query, args)

	return row
}

func (o *SqlOperations) Update(entity interface{}, set Set, where Where) (sql.Result, error) {
	table := o.getTable(entity)
	setStr, setArgs := set.handle()
	whereStr, whereArgs := where.handleUpdateWhere(len(setArgs))
	args := append(setArgs, whereArgs...)

	query := fmt.Sprintf("UPDATE teste.%s %s %s;", table, setStr, whereStr)

	return postgres.Client.Exec(query, args...)
}

func (o *SqlOperations) Delete(entity interface{}, id int) (sql.Result, error) {
	table := o.getTable(entity)

	query := fmt.Sprintf("DELETE FROM teste.%s WHERE id = $1;", table)

	return postgres.Client.Exec(query, id)
}

func (o *SqlOperations) getTable(entity interface{}) (tableName string) {
	typeEntity := reflect.TypeOf(entity)
	split := strings.Split(typeEntity.Name(), ".")
	tableName = fmt.Sprintf("%ss", split[len(split)-1])
	return
}

func (o *SqlOperations) getColumns(entity interface{}) (columns string) {
	typeEntity := reflect.TypeOf(entity)

	for i := 0; i < typeEntity.NumField(); i++ {
		column := typeEntity.Field(i)

		if column.Name == "Id" {
			continue
		}

		if columns != "" {
			columns = fmt.Sprintf("%s, %s", columns, column.Name)
		} else {
			columns = column.Name
		}
	}

	return
}

func (o *SqlOperations) getValues(entity interface{}) (valuesStr string, values []interface{}) {
	valueEntity := reflect.ValueOf(entity)
	index := 1

	for i := 0; i < valueEntity.NumField(); i++ {
		value := valueEntity.Field(i).Interface()

		if value == "" || value == 0 {
			continue
		}

		if valuesStr != "" {
			valuesStr = fmt.Sprintf("%s, $%v", valuesStr, index)
			index++
		} else {
			valuesStr = "$1"
			index++
		}

		values = append(values, value)

	}
	return
}
