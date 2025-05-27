package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	// If you're not using this, remove it
)

type cmdFlags struct {
	Add    string
	Del    int
	Edit   string
	Toggle int
	List   bool
}

func NewCmdFlags() *cmdFlags {
	cf := cmdFlags{}

	flag.StringVar(&cf.Add, "add", "", "Add a new todo title")
	flag.StringVar(&cf.Edit, "Edit", "", "Edit a todo by index")
	flag.IntVar(&cf.Del, "del", -1, "Specify a todo by index to delete")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Specify a todo by index to toggle")
	flag.BoolVar(&cf.List, "List", false, "List all todos") // ✅ comma added

	flag.Parse()

	return &cf
}

func (cf *cmdFlags) Execute(todos *Todos) { // ✅ syntax fixed
	switch {
	case cf.List:
		todos.print()
	case cf.Add != "":
		todos.add(cf.Add)
		todos.print()
	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Error: Invalid format. Use <index>:<new title>")
			os.Exit(1)
		}
		index, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error: Invalid index")
			os.Exit(1)
		}
		todos.edit(index, parts[1])
		todos.print()
	case cf.Toggle != -1:
		todos.toggle(cf.Toggle)
		todos.print()
	case cf.Del != -1:
		todos.delete(cf.Del)
		todos.print()
	default:
		fmt.Println("Invalid command")
	}
}
