package todo

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
)

type Todo struct {
	ID   int
	Text string
	Done bool
}

type State struct {
	Todos []Todo
	Index int
}

func DisplayTodos(g *gocui.Gui, state *State, todos []Todo) error {
	list_view, err := g.SetCurrentView("todoList")
	if err != nil {
		return err
	}
	list_view.Clear()
	for index, todo := range todos {
		if todo.Done {
			if index == state.Index {
				fmt.Fprintf(list_view, ">> %d. %s [x] \n", index+1, todo.Text)
			} else {
				fmt.Fprintf(list_view, "[x] %d. %s\n", index+1, todo.Text)
			}
		} else if index == state.Index {
			if todo.Done {
				fmt.Fprintf(list_view, ">> %d. %s [x] \n", index+1, todo.Text)
			} else {
				fmt.Fprintf(list_view, ">> %d. %s\n", index+1, todo.Text)
			}
		} else {
			fmt.Fprintf(list_view, "%d. %s\n", index+1, todo.Text)
		}
	}
	return nil
}

func (state *State) ShowDoneTasks(g *gocui.Gui, v *gocui.View) error {
	doneTodos := []Todo{}
	for _, todo := range state.Todos {
		if todo.Done {
			doneTodos = append(doneTodos, todo)
		}
	}
	DisplayTodos(g, state, doneTodos)
	return nil
}

func (state *State) ShowAllTasks(g *gocui.Gui, v *gocui.View) error {
	DisplayTodos(g, state, state.Todos)
	return nil
}

func (state *State) ShowRemainingTasks(g *gocui.Gui, v *gocui.View) error {
	remTodos := []Todo{}
	for _, todo := range state.Todos {
		if !todo.Done {
			remTodos = append(remTodos, todo)
		}
	}
	DisplayTodos(g, state, remTodos)
	return nil
}

func (state *State) AddTodo(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == "todoInput" {

		state.Todos = append(state.Todos, Todo{ID: len(state.Todos) + 1, Text: strings.ReplaceAll(v.Buffer(), "\n", "")})
		state.Index = len(state.Todos) - 1

		v.Clear()
		g.SetViewOnBottom("todoInput")

		err := DisplayTodos(g, state, state.Todos)

		if err != nil {
			return err
		}

	}
	return nil
}

func (state *State) PrevTodo(g *gocui.Gui, v *gocui.View) error {
	if state.Index > 0 {
		state.Index--
		err := DisplayTodos(g, state, state.Todos)
		if err != nil {
			return err
		}
	}
	return nil
}

func (state *State) NextTodo(g *gocui.Gui, v *gocui.View) error {
	if state.Index < len(state.Todos)-1 {
		state.Index++
		err := DisplayTodos(g, state, state.Todos)
		if err != nil {
			return err
		}
	}
	return nil
}

func (state *State) DeleteTodo(g *gocui.Gui, v *gocui.View) error {
	list_view, err := g.SetCurrentView("todoList")
	_, err = list_view.Line(state.Index)
	if err != nil {
		return err
	}

	state.Todos = removeElementFromAnArray(state.Todos, state.Index)
	state.Index = 0
	err = DisplayTodos(g, state, state.Todos)
	if err != nil {
		return err
	}
	return nil
}

func (state *State) MarkToDoAsDone(g *gocui.Gui, v *gocui.View) error {
	state.Todos[state.Index].Done = true
	DisplayTodos(g, state, state.Todos)
	return nil
}

func removeElementFromAnArray[T Todo | int](arr []T, index int) []T {
	if index < 0 || index >= len(arr) {
		return arr
	}
	result := make([]T, len(arr)-1)

	copy(result, arr[:index])
	copy(result[index:], arr[index+1:])

	return result
}
