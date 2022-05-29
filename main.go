package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
)

type User struct {
	Id          uuid.UUID
	Name        string
	PhoneNumber string
	Email       string
	CreatedDate time.Time
}

var dataBaseURL string = "postgres://login:pass@localhost:5432/database-name"

func main() {
	db, err := sql.Open("pgx", dataBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func Create(db *sql.DB, user User) (uuid.UUID, error) {
	id := uuid.New()
	currentTime := time.Now()
	_, err := db.Exec("INSERT INTO users(id, name, phone, email, created_at) VALUES ($1, $2, $3, $4, $5)", id, user.Name, user.PhoneNumber, user.Email, currentTime)
	if err != nil {
		return uuid.Nil, err
	}
	fmt.Println("Row has been created")
	return id, nil
}

func GetById(db *sql.DB, id uuid.UUID) (User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.Name, &user.PhoneNumber, &user.Email, &user.CreatedDate)
	if err != nil {
		return user, err
	}
	fmt.Printf("%-36s %-10s %-15s %-20s %s\n", "id", "name", "phone", "email", "created_at")
	fmt.Printf("%-36s %-10s %-15s %-20s %v\n", user.Id, user.Name, user.PhoneNumber, user.Email, user.CreatedDate)
	return user, nil
}

func UpdateById(db *sql.DB, id uuid.UUID, fields map[string]interface{}, template User) error {
	queryText := "UPDATE users SET"
	fieldsList, err := fieldsFilter(fields, template)
	if err != nil {
		return err
	}
	queryText, queryFields := getQueryDataForUpdate(queryText, fieldsList)
	//---debug
	fmt.Println(queryText)
	fmt.Println(queryFields)
	//---
	queryFields = append(queryFields, id)
	_, err = db.Exec(queryText, queryFields...)
	if err != nil {
		return err
	}
	fmt.Println("Row has been updated")
	return nil
}

func DeleteById(db *sql.DB, id uuid.UUID) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	fmt.Println("Row has been deleted")
	return nil
}

func List(db *sql.DB, limit, offest uint64) ([]User, error) {
	var (
		user         User
		limitString  string
		offsetString string
		users        []User
	)
	if limit != 0 {
		limitString = " LIMIT " + strconv.Itoa(int(limit))
	}
	if offest != 0 {
		offsetString = " OFFSET " + strconv.Itoa(int(offest))
	}
	rows, err := db.Query("SELECT * FROM users" + limitString + offsetString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//----- для отладки, убрать в скором времени
	rowsNames, err := rows.Columns()
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("%-36s %-10s %-15s %-20s %s\n", rowsNames[0], rowsNames[1], rowsNames[2], rowsNames[3], rowsNames[4])
	//-----
	users = make([]User, 0, limit)
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.PhoneNumber, &user.Email, &user.CreatedDate)
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
		fmt.Printf("%-36s %-10s %-15s %-20s %v\n", user.Id, user.Name, user.PhoneNumber, user.Email, user.CreatedDate)
	}
	return users, nil
}

func fieldsFilter(fields map[string]interface{}, template User) (filteredFields map[string]interface{}, err error) {
	v := reflect.ValueOf(template)
	typeOfTemplate := v.Type()
	if v.Kind() != reflect.Struct {
		return nil, errors.New("Wrong template parameter. It must be struct or have underlying type as struct")
	}
	filteredFields = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		templateFieldName := typeOfTemplate.Field(i).Name
		templateFieldType := v.Field(i).Type()
		if value, ok := fields[templateFieldName]; ok && value != nil && reflect.ValueOf(value).Type() == templateFieldType {
			filteredFields[templateFieldName] = value
		}
	}
	if len(filteredFields) == 0 {
		return filteredFields, errors.New("Not enough fields to update.")
	}
	return
}

func getQueryDataForUpdate(queryText string, filteredFields map[string]interface{}) (outputQueryText string, fieldsValues []interface{}) {
	fieldsValues = make([]interface{}, 0, len(filteredFields))
	comma := ","
	i := 0
	for k, v := range filteredFields {
		if i == len(filteredFields)-1 {
			comma = " "
		}
		queryText += " " + k + " = " + "$" + strconv.Itoa(i+1) + comma
		fieldsValues = append(fieldsValues, v)
		i++
	}
	queryText += "WHERE id = " + "$" + strconv.Itoa(i+1)
	return queryText, fieldsValues
}
