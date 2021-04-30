package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	id       int
	username string
	surname  string
}

func main() {
	db, _ := sql.Open("sqlite3", "godb.db")
	db.Exec("create table if not exists usersTable (id integer primary key autoincrement,username text, surname text)")
loop:
	for {
		var answer int

		user := User{}
		fmt.Println("MENU:\n1. Add new user\n2. Get all users\n3. Delete user\n4. Exit")

		fmt.Print("Choose command, please: ")
		_, err := fmt.Fscan(os.Stdin, &answer)
		checkError(err)

		switch answer {
		case 1:

			fmt.Print("Enter first name, please: ")
			fmt.Fscan(os.Stdin, &user.username)

			fmt.Print("Enter last name, please: ")
			fmt.Fscan(os.Stdin, &user.surname)

			addUser(db, user.username, user.surname)

		case 2:
			fmt.Println("List of users from db: ")
			getUsers(db)

		case 3:
			fmt.Print("Enter the id of the user you want to delete, please: ")
			fmt.Fscan(os.Stdin, &user.id)

			deleteUser(db, user.id)

		case 4:
			break loop
		}

	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func addUser(db *sql.DB, username string, surname string) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into usersTable (username,surname) values (?,?)")
	_, err := stmt.Exec(username, surname)
	checkError(err)
	tx.Commit()
}

func getUsers(db *sql.DB) {
	userList := make([]*User, 0)
	rows, err := db.Query("select * from usersTable")
	checkError(err)

	for rows.Next() {
		u := new(User)
		err := rows.Scan(&u.id, &u.username, &u.surname)
		checkError(err)
		userList = append(userList, u)
	}

	for _, u := range userList {
		fmt.Printf("%d, %s, %s\n", u.id, u.username, u.surname)
	}

}

func deleteUser(db *sql.DB, id int) {
	sid := strconv.Itoa(id)
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("delete from usersTable where id=?")
	_, err := stmt.Exec(sid)
	checkError(err)
	tx.Commit()
}
