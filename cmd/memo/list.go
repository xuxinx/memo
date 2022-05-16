package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newList(title string, items []list.Item) list.Model {
	h := len(items) + 5
	if h > 20 {
		h = 20
	}
	l := list.New(items, listItemDelegate{}, 50, h)
	l.Title = title
	if l.Title == "" {
		l.SetShowTitle(false)
	}
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	return l
}

type listItem struct {
	title   string
	context map[string]interface{}
}

func (i listItem) FilterValue() string {
	return i.title
}

type listItemDelegate struct{}

func (d listItemDelegate) Height() int {
	return 1
}
func (d listItemDelegate) Spacing() int {
	return 0
}
func (d listItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
func (d listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(listItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("  %d. %s", index+1, i.title)
	if index == m.Index() {
		str = focusStyle.Render(fmt.Sprintf("> %d. %s", index+1, i.title))
	}

	fmt.Fprintf(w, str)
}
