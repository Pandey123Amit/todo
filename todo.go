package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aquasecurity/table"
)

// Todo represents a single todo item
type Todo struct {
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// Todos is a slice of Todo
type Todos []Todo

// Add a new todo
func (todos *Todos) add(title string) {
	todo := Todo{
		Title:       title,
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   time.Now(),
	}
	*todos = append(*todos, todo)
}

// Validate index exists in the list
func (todos *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*todos) {
		return errors.New("invalid index")
	}
	return nil
}

// Delete a todo by index
func (todos *Todos) delete(index int) error {
	if err := todos.validateIndex(index); err != nil {
		return err
	}
	*todos = append((*todos)[:index], (*todos)[index+1:]...)
	return nil
}

// Toggle completion of a todo
func (todos *Todos) toggle(index int) error {
	if err := todos.validateIndex(index); err != nil {
		return err
	}
	todo := &(*todos)[index]
	if !todo.Completed {
		now := time.Now()
		todo.CompletedAt = &now
	} else {
		todo.CompletedAt = nil
	}
	todo.Completed = !todo.Completed
	return nil
}

// Edit a todo title
func (todos *Todos) edit(index int, title string) error {
	if err := todos.validateIndex(index); err != nil {
		return err
	}
	preTitel := (*todos)[index].Title
	data := EmailData{
		Name:     "Amit Pandey",
		Subject:  "Task is Edited",
		Body:     "Thank you for signing up for our todo application!",
		Prevtask: preTitel,
		Newtask:  title,
		Id:       index,
	}
	tmpl, err := template.ParseFiles("edit_template.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Fatal("Error executing template:", err)
	}
	(*todos)[index].Title = title
	err = sendEmail("checkemail@yopmail.com", data.Subject, body.String())
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
	fmt.Println("sss")
	return nil
}

// Print the todo list
func (todos *Todos) print() {
	t := table.New(os.Stdout)
	t.SetRowLines(false)
	t.SetHeaders("#", "Title", "Completed", "Created At", "Completed At", "Time Taken")

	for i, todo := range *todos {
		completed := "❌"
		completedAt := "-"
		timeTaken := ""

		if todo.Completed {
			completed = "✅"
			if todo.CompletedAt != nil {
				completedAt = todo.CompletedAt.Format(time.RFC1123)
				timeTaken = todo.CompletedAt.Sub(todo.CreatedAt).String()
			}
		} else {
			timeTaken = time.Since(todo.CreatedAt).String()
		}

		t.AddRow(
			strconv.Itoa(i),
			todo.Title,
			completed,
			todo.CreatedAt.Format(time.RFC1123),
			completedAt,
			timeTaken,
		)
	}

	t.Render()
}
