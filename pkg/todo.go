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

func DisplayTodos(g *gocui.Gui, state *State) error {
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

func (state *State) AddTodo(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == "todoInput" {

		state.Todos = append(state.Todos, Todo{ID: len(state.Todos) + 1, Text: strings.ReplaceAll(v.Buffer(), "\n", "")})
		state.Index = len(state.Todos) - 1

		v.Clear()
		g.SetViewOnBottom("todoInput")

		err := DisplayTodos(g, state)

		if err != nil {
			return err
		}

	}
	return nil
}


func (state *State) PrevTodo(g *gocui.Gui, v *gocui.View) error {
	if state.Index > 0 {
		state.Index--
		err := DisplayTodos(g, state)
		if err != nil {
			return err
		}
	}
	return nil
}

func (state *State) NextTodo(g *gocui.Gui, v *gocui.View) error {
	if state.Index < len(state.Todos)-1 {
		state.Index++
		err := DisplayTodos(g, state)
		if err != nil {
			return err
		}
	}
	return nil
}

func (state *State) DeleteTodo(g *gocui.Gui, v *gocui.View) error {
	list_view, err := g.SetCurrentView("todoList")
	_vw, err := g.View("status")
	_vw.Clear()
	_, err = list_view.Line(state.Index)
	if err != nil {
		return err
	}

	state.Todos = removeElementFromAnArray(state.Todos, state.Index)
	state.Index = 0
	fmt.Fprintln(_vw, state.Todos)

	err = DisplayTodos(g, state)
	if err != nil {
		return err
	}
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
