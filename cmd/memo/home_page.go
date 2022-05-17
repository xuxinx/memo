package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xuxinx/memo"
)

const pageHome page = "Home"

type homePage struct {
	list list.Model
}

func initHomePage() pager {
	return &homePage{
		list: newList(
			"Home",
			[]list.Item{
				listItem{title: string(pagePractice), value: pagePractice},
				listItem{title: "New Question", value: pageEditQuestion},
				listItem{title: string(pageList), value: pageList},
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
			m.page = item.value.(page)
			switch m.page {
			case pagePractice:
				m.pagers[m.page] = initPracticePage(true)
			case pageEditQuestion:
				m.pagers[m.page] = initEditQuestionPage(&memo.Question{}, pageHome)
			case pageList:
				m.pagers[m.page] = initListPage(true)
			}
			return m, nil
		case "/":
			return m, nil
		case "ctrl+h":
			m.page = pageHelp
			m.pagers[m.page] = initHelpPage(
				`
ctrl+c / q: quit
     j / k: move
     enter: select item 
                `,
				pageHome,
			)
			return m, nil
		}
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return m, cmd
}

func (p *homePage) view(m *model) string {
	return p.list.View()
}
