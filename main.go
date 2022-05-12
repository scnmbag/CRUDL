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
	db := OpenDB(dataBaseURL)
	defer db.Close()
}

func OpenDB(url string) *sql.DB {
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Create(data dbUserEntity, dataBase *sql.DB) {
	uuid := uuid.New()
	//time := time.Now() использую now() в postgreSQL
	_, err := dataBase.Exec("INSERT INTO users(id, name, phone, email, created_at) VALUES ($1, $2, $3, $4, now())", uuid, data.name, data.phoneNumber, data.email)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Entry have been created")
}

//какой тип использовать для параметра id? uuid.UUID or [16]byte?
func GetById(id [16]byte, dataBase *sql.DB) {
	var user dbUserEntity
	row := dataBase.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.id, &user.name, &user.phoneNumber, &user.email, &user.createdDate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-36s %-10s %-15s %-20s %s\n", "id", "name", "phone", "email", "created_at")
	fmt.Printf("%-36s %-10s %-15s %-20s %v\n", user.id, user.name, user.phoneNumber, user.email, user.createdDate)
}

func UpdateById(id [16]byte, fields map[string]string, dataBase *sql.DB) {
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
		log.Fatal(err)
	}
	fmt.Println("Row have been updated")
}

func DeleteById(id [16]byte, dataBase *sql.DB) {
	_, err := dataBase.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Row have been deleted")
}

func List(dataBase *sql.DB, limit int) {
	var limitString string
	var user dbUserEntity
	if limit > 0 {
		limitString = " LIMIT " + strconv.Itoa(limit)
	}
	rows, err := dataBase.Query("SELECT * FROM users" + limitString)
	if err != nil {
		log.Fatal(err)
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
}

func List1(a ...any) {
	var (
		dataBase      *sql.DB
		limitQuantity int
		limitString   string
		user          dbUserEntity
	)
	for _, v := range a {
		switch v := v.(type) {
		case *sql.DB:
			dataBase = v
		case int:
			limitQuantity = v
		}
	}
	if dataBase == nil {
		log.Fatal("This function must take a *sql.DB to work")
	}
	if limitQuantity > 0 {
		limitString = " LIMIT " + strconv.Itoa(limitQuantity)
	}
	rows, err := dataBase.Query("SELECT * FROM users" + limitString)
	if err != nil {
		log.Fatal(err)
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
}
