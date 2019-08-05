//go:generate go run github.com/99designs/gqlgen
package main

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	todos []*Todo
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (*Todo, error) {
	todo := &Todo{
		Text: input.Text,
		// ID:     fmt.Sprintf("T%d", rand.Int()),
		UserID: input.UserID,
	}

	saveNewTodo(ctx, todo)
	// r.todos = append(r.todos, todo)
	return todo, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*User, error) {
	return saveNewUser(ctx, input.Name)
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input UpdateTodoInfo) (*Todo, error) {
	todo := &Todo{
		ID:   input.ID,
		Text: input.Text,
		Done: input.Done,
	}

	return updateTodo(ctx, todo)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input UpdateUserInfo) (*User, error) {
	user := &User{
		ID:   input.ID,
		Name: input.Name,
	}

	return updateUser(ctx, user)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]*Todo, error) {
	return loadTodos(ctx)
	// return r.todos, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*User, error) {
	return loadUsers(ctx)
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *Todo) (*User, error) {
	return loadUser(ctx, obj.UserID)
	// return &User{ID: obj.UserID, Name: "user " + obj.UserID}, nil
}
