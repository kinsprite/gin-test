package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type userItem struct {
	gorm.Model
	Name string
}

type todoItem struct {
	gorm.Model
	Text   string
	Done   bool
	UserID string
}

func (item *todoItem) fromNewTodo(todo *Todo) {
	item.Text = todo.Text
	item.Done = todo.Done
	item.UserID = todo.UserID
}

func (item *todoItem) toTodo() *Todo {
	var todo Todo
	todo.ID = strconv.FormatInt(int64(item.ID), 10)
	todo.Text = item.Text
	todo.Done = item.Done
	todo.UserID = item.UserID
	return &todo
}

func initDB() {
	driverName := os.Getenv("SQL_DRIVER_NAME")
	dataSourceName := os.Getenv("SQL_DATA_SOURCE_NAME")

	var err error
	db, err = gorm.Open(driverName, dataSourceName)

	if err != nil {
		log.Println("ERROR    Gorm open DB failed")
	} else {
		log.Println("INFO    Gorm open DB OK")
	}

	// Migrate the schema
	db.AutoMigrate(&userItem{}, &todoItem{})

	// if isEmptyBooks() {
	// 	log.Println("INFO    Gorm DB books is Empty")
	// 	createUserBooks()
	// }
}

func closeDB() {
	if db != nil {
		db.Close()
	}
}

func saveNewTodo(ctx context.Context, todo *Todo) {
	item := todoItem{}
	item.fromNewTodo(todo)
	db.Save(&item)
	todo.ID = strconv.FormatInt(int64(item.ID), 10)
}

func loadTodos(ctx context.Context) ([]*Todo, error) {
	todoItems := []todoItem{}
	errors := db.Find(&todoItems).GetErrors()

	todos := make([]*Todo, len(todoItems))

	for i, item := range todoItems {
		todos[i] = item.toTodo()
	}

	var err error

	if len(errors) > 0 {
		err = errors[0]
	}

	return todos, err
}

func saveNewUser(ctx context.Context, name string) (*User, error) {
	item := userItem{
		Name: name,
	}

	errors := db.Save(&item).GetErrors()

	var err error

	if len(errors) > 0 {
		err = errors[0]
		return nil, err
	}

	user := User{
		ID:   strconv.FormatInt(int64(item.ID), 10),
		Name: name,
	}

	return &user, err
}

func loadUser(ctx context.Context, userID string) (*User, error) {
	id, err := strconv.ParseUint(userID, 10, 0)

	if err != nil {
		return nil, err
	}

	item := userItem{}
	errors := db.Where("id = ?", id).Find(&item).GetErrors()

	if len(errors) > 0 {
		err = errors[0]
		return nil, err
	}

	user := User{
		ID:   strconv.FormatInt(int64(item.ID), 10),
		Name: item.Name,
	}

	return &user, nil
}
