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

func (item *userItem) toUser() *User {
	user := User{
		ID:   idToStr(item.ID),
		Name: item.Name,
	}

	return &user
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
	todo.ID = idToStr(item.ID)
	todo.Text = item.Text
	todo.Done = item.Done
	todo.UserID = item.UserID
	return &todo
}

func strToID(strID string) (uint64, error) {
	return strconv.ParseUint(strID, 10, 0)
}

func idToStr(ID uint) string {
	return strconv.FormatInt(int64(ID), 10)
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
	todo.ID = idToStr(item.ID)
}

func updateTodo(ctx context.Context, todo *Todo) (*Todo, error) {
	item := todoItem{}
	id, err := strToID(todo.ID)

	if err != nil {
		return nil, err
	}

	item.ID = uint(id)
	item.Text = todo.Text
	item.Done = todo.Done

	errors := db.Model(&item).Updates(item).GetErrors()

	if len(errors) > 0 {
		err = errors[0]
		return nil, err
	}

	errors = db.First(&item, item.ID).GetErrors()

	if len(errors) > 0 {
		err = errors[0]
		return nil, err
	}

	return item.toTodo(), err
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

	user := item.toUser()
	return user, err
}

func updateUser(ctx context.Context, user *User) (*User, error) {
	item := userItem{}

	id, err := strToID(user.ID)

	if err != nil {
		return nil, err
	}

	item.ID = uint(id)
	item.Name = user.Name

	errors := db.Model(&item).Updates(item).GetErrors()

	if len(errors) > 0 {
		err = errors[0]
		return nil, err
	}

	errors = db.First(&item, item.ID).GetErrors()

	if len(errors) > 0 {
		err = errors[0]
		return nil, err
	}

	return item.toUser(), err
}

func loadUsers(ctx context.Context) ([]*User, error) {
	userItems := []userItem{}
	errors := db.Find(&userItems).GetErrors()

	users := make([]*User, len(userItems))

	for i, item := range userItems {
		users[i] = item.toUser()
	}

	var err error

	if len(errors) > 0 {
		err = errors[0]
	}

	return users, err
}

func loadUser(ctx context.Context, userID string) (*User, error) {
	id, err := strToID(userID)

	if err != nil {
		return nil, err
	}

	item := userItem{}
	errors := db.Where("id = ?", id).Find(&item).GetErrors()

	if len(errors) > 0 {
		err = errors[0]
		return nil, err
	}

	user := item.toUser()
	return user, nil
}
