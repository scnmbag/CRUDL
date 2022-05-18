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

type User struct {
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

func Create(db *sql.DB, user User) (uuid.UUID, error) {
	id := uuid.New()
	currentTime := time.Now()
	_, err := db.Exec("INSERT INTO users(id, name, phone, email, created_at) VALUES ($1, $2, $3, $4, $5)", id, user.name, user.phoneNumber, user.email, currentTime)
	if err != nil {
		return uuid.Nil, err
	}
	fmt.Println("Row has been created")
	return id, nil
}

func GetById(db *sql.DB, id uuid.UUID) (User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.id, &user.name, &user.phoneNumber, &user.email, &user.createdDate)
	if err != nil {
		return user, err
	}
	fmt.Printf("%-36s %-10s %-15s %-20s %s\n", "id", "name", "phone", "email", "created_at")
	fmt.Printf("%-36s %-10s %-15s %-20s %v\n", user.id, user.name, user.phoneNumber, user.email, user.createdDate)
	return user, nil
}

func UpdateById(db *sql.DB, id uuid.UUID, fields map[string]string) error {
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
	_, err := db.Exec(queryText, id)
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
		err = rows.Scan(&user.id, &user.name, &user.phoneNumber, &user.email, &user.createdDate)
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
		fmt.Printf("%-36s %-10s %-15s %-20s %v\n", user.id, user.name, user.phoneNumber, user.email, user.createdDate)
	}
	return users, nil
}
