package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const pageHome page = "Home"

type homePage struct {
	list list.Model
}

func initHomePage() *homePage {
	return &homePage{
		list: newList(
			"Home",
			[]list.Item{
				listItem{title: string(pagePractice)},
				listItem{title: string(pageNewQuestion)},
			},
		),
	}
}

func (p *homePage) update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			item := p.list.SelectedItem().(listItem)
			m.page = page(item.title)
			switch m.page {
			case pagePractice:
				m.practicePage = initPracticePage(true)
			case pageNewQuestion:
				m.newQuestionPage = initNewQuestionPage()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return m, cmd
}

func (p *homePage) view() string {
	return p.list.View()
}
