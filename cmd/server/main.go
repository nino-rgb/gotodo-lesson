package main

import (
	"database/sql"
	"fmt"
	"go-todo/repositories"
	"go-todo/services"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASS")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user,
		password,
		host,
		port,
		dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB接続成功")

	repo := repositories.NewTodoRepository(db)
	service := services.NewTodoService(repo)
	todos, err := service.GetTodos()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(todos)
}
