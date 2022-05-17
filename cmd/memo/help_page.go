package main

import tea "github.com/charmbracelet/bubbletea"

const pageHelp page = "Help"

type helpPage struct {
	helps    string
	backPage page
}

func initHelpPage(helps string, backPage page) pager {
	return &helpPage{
		helps:    helps,
		backPage: backPage,
	}
}

func (p *helpPage) update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			m.page = p.backPage
			return m, nil
		}
	}

	return m, nil
}

func (p *helpPage) view(m *model) string {
	return "Helps:\n" + p.helps + "\n"
}
