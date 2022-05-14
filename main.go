package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
)

type dbUserEntity struct {
	id          uuid.UUID
	name        string
	phoneNumber string
	email       string
	createdDate time.Time
}

var dataBaseURL string = "postgres://login:pass@localhost:5432/database-name"

func main() {
	db, err := sql.Open("pgx", dataBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func Create(data dbUserEntity, dataBase *sql.DB) error {
	uuid := uuid.New()
	currentTime := time.Now()
	_, err := dataBase.Exec("INSERT INTO users(id, name, phone, email, created_at) VALUES ($1, $2, $3, $4, $5)", uuid, data.name, data.phoneNumber, data.email, currentTime)
	if err != nil {
		return err
	}
	fmt.Println("Row has been created")
	return nil
}

func GetById(id uuid.UUID, dataBase *sql.DB) error {
	var user dbUserEntity
	row := dataBase.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.id, &user.name, &user.phoneNumber, &user.email, &user.createdDate)
	if err != nil {
		return err
	}
	fmt.Printf("%-36s %-10s %-15s %-20s %s\n", "id", "name", "phone", "email", "created_at")
	fmt.Printf("%-36s %-10s %-15s %-20s %v\n", user.id, user.name, user.phoneNumber, user.email, user.createdDate)
	return nil
}

func UpdateById(id uuid.UUID, fields map[string]string, dataBase *sql.DB) error {
	var (
		updateString string
		queryText    string
		fieldsNames  []string
		comma        = "', "
	)
	for key := range fields {
		fieldsNames = append(fieldsNames, key)
	}
	for i, v := range fieldsNames {
		if i == len(fields)-1 {
			comma = "'"
		}
		updateString += v + " = " + "'" + fields[v] + comma
	}
	queryText = "UPDATE users SET " + updateString + " WHERE id = $1"
	_, err := dataBase.Exec(queryText, id)
	if err != nil {
		return err
	}
	fmt.Println("Row has been updated")
	return nil
}

func DeleteById(id uuid.UUID, dataBase *sql.DB) error {
	_, err := dataBase.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	fmt.Println("Row has been deleted")
	return nil
}

func List(dataBase *sql.DB, limit ...any) error {
	var limitString string
	var user dbUserEntity
	for _, v := range limit {
		limitNumber, ok := v.(int)
		if ok && limitNumber > 0 {
			limitString = " LIMIT " + strconv.Itoa(limitNumber)
		}
	}
	rows, err := dataBase.Query("SELECT * FROM users" + limitString)
	if err != nil {
		return err
	}
	defer rows.Close()
	rowsNames, err := rows.Columns()
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("%-36s %-10s %-15s %-20s %s\n", rowsNames[0], rowsNames[1], rowsNames[2], rowsNames[3], rowsNames[4])
	for rows.Next() {
		rows.Scan(&user.id, &user.name, &user.phoneNumber, &user.email, &user.createdDate)
		fmt.Printf("%-36s %-10s %-15s %-20s %v\n", user.id, user.name, user.phoneNumber, user.email, user.createdDate)
	}
	return nil
}
