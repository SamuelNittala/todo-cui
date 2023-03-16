package main

import (
	"fmt"
	"strings"

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

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Println("Error initializing CUI:", err)
		return
	}
	defer g.Close()

	state := State{
		Todos: []Todo{{1, "Buy groceries"}, {2, "Do laundry"}, {3, "Clean the house"}},
		Index: 0,
	}

	g.SetManagerFunc(func(g *gocui.Gui) error {
		return layout(g, &state)
	})

	if err := keybindings(g, &state); err != nil {
		fmt.Println("Error setting keybindings:", err)
		return
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Println("Error running CUI:", err)
	}
}


func layout(g *gocui.Gui, state *State) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("todoList", 0, 0, maxX-1, maxY-3); err != nil {
		v.Title = "Todo-List"
		if err != gocui.ErrUnknownView {
			return err
		}

		for index, todo := range state.Todos {
			if index == state.Index {
				fmt.Fprintf(v, ">> %d. %s\n", index+1, todo.Text)
			} else {
				fmt.Fprintf(v, "%d. %s\n", index+1, todo.Text)
			}
		}

		if _, err := g.SetCurrentView("todoList"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("todoInput", 0, 0, maxX/2, maxY/4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Enter new todo"
		g.SetViewOnBottom("todoInput")
	}

	if v, err := g.SetView("status", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Press <a> to add a new todo. Press <q> to quit.")
	}

	return nil
}

func keybindings(g *gocui.Gui, state *State) error {
	if err := g.SetKeybinding("todoList", 'a', gocui.ModNone, showTodoInput); err != nil {
		return err
	}

	if err := g.SetKeybinding("todoList", 'd', gocui.ModNone, state.deleteTodo); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, state.addTodo); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, state.nextTodo); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, state.prevTodo); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func showTodoInput(g *gocui.Gui, v *gocui.View) error {
	input_view, err := g.SetCurrentView("todoInput")
	input_view.Editable = true
	input_view.SetCursor(0, 0)

	g.SetViewOnTop("todoInput")
	if err != nil {
		return err
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

func (state *State) nextTodo(g *gocui.Gui, v *gocui.View) error {
	if state.Index < len(state.Todos)-1 {
		state.Index++
		err := displayTodos(g, state)
		if err != nil {
			return err
		}
	}
	return nil
}

func (state *State) prevTodo(g *gocui.Gui, v *gocui.View) error {
	if state.Index > 0 {
		state.Index--
		err := displayTodos(g, state)
		if err != nil {
			return err
		}
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

func (state *State) deleteTodo(g *gocui.Gui, v *gocui.View) error {
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

	err = displayTodos(g, state)
	if err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
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
