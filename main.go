package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/samuelnittala/todo-app/pkg"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Println("Error initializing CUI:", err)
		return
	}
	defer g.Close()

	state := todo.State{
		Todos: []todo.Todo{{1, "Buy groceries", false}, {2, "Do laundry", false}, {3, "Clean the house", false}},
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

func layout(g *gocui.Gui, state *todo.State) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("todoList", 0, 0, maxX/2, maxY-3); err != nil {
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

	if v, err := g.SetView("status", 0, maxY-6, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Press <a> to add a new todo. <m> to mark as done, <d> to delete Press <q> to quit.")
		fmt.Fprintln(v, "\n")
		fmt.Fprintln(v, "Press <F1> to show all todos. <F2> to show done todos, <F3> to show remainig todos.")
	}

	return nil
}

func keybindings(g *gocui.Gui, state *todo.State) error {
	if err := g.SetKeybinding("todoList", 'a', gocui.ModNone, showTodoInput); err != nil {
		return err
	}

	if err := g.SetKeybinding("todoList", 'd', gocui.ModNone, state.DeleteTodo); err != nil {
		return err
	}

	if err := g.SetKeybinding("todoList", 'm', gocui.ModNone, state.MarkToDoAsDone); err != nil {
		return err
	}

	if err := g.SetKeybinding("todoList", gocui.KeyF1, gocui.ModNone, state.ShowAllTasks); err != nil {
		return err
	}

	if err := g.SetKeybinding("todoList", gocui.KeyF2, gocui.ModNone, state.ShowDoneTasks); err != nil {
		return err
	}

	if err := g.SetKeybinding("todoList", gocui.KeyF3, gocui.ModNone, state.ShowRemainingTasks); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, state.AddTodo); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, state.NextTodo); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, state.PrevTodo); err != nil {
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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
