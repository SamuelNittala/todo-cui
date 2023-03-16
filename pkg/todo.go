package todo

import (
  "strings"
  "fmt"
	"github.com/jroimartin/gocui"
)

type Todo struct {
	ID   int
	Text string
}

type State struct {
	Todos []Todo
	Index int
}

func displayTodos(g *gocui.Gui, state *State) error {
	list_view, err := g.SetCurrentView("todoList")
	if err != nil {
		return err
	}
	list_view.Clear()
	v, err := g.View("status")
	v.Clear()
	fmt.Fprintln(v, state.Index)
	for index, todo := range state.Todos {
		if index == state.Index {
			fmt.Fprintf(list_view, ">> %d. %s\n", index + 1, todo.Text)
		} else {
			fmt.Fprintf(list_view, "%d. %s\n", index + 1, todo.Text)
		}
	}
	return nil
}

func (state *State) addTodo(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == "todoInput" {

		state.Todos = append(state.Todos, Todo{ID: len(state.Todos) + 1, Text: strings.ReplaceAll(v.Buffer(), "\n", "")})
		state.Index = len(state.Todos) - 1

		v.Clear()
		g.SetViewOnBottom("todoInput")

		err := displayTodos(g, state)

		if err != nil {
			return err
		}

	}
	return nil
}
