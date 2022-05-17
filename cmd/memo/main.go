package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xuxinx/memo"
)

const (
	fPath = "./memo.json"
)

var (
	serv *memo.Memo

	noStyle     = lipgloss.NewStyle()
	focusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#24b9da"))
	cursorStyle = focusStyle.Copy()
)

func main() {
	f := mustGetOrInitFile(fPath)
	defer f.Close()
	serv = mustGetMemoService(f)

	if err := tea.NewProgram(initModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

type page string

type pager interface {
	update(m *model, msg tea.Msg) (tea.Model, tea.Cmd)
	view(m *model) string
}

type model struct {
	page     page
	pagers   map[page]pager
	quitting bool
}

func initModel() *model {
	pagers := make(map[page]pager)
	pagers[pageHome] = initHomePage()
	pagers[pagePractice] = initPracticePage(false)
	pagers[pageList] = initListPage(false)
	pagers[pageEditQuestion] = initEditQuestionPage(nil, "")
	pagers[pageHelp] = initHelpPage("", "")

	return &model{
		page:     pageHome,
		pagers:   pagers,
		quitting: false,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	pager, ok := m.pagers[m.page]
	if !ok {
		panic("unknown page")
	}
	return pager.update(m, msg)
}

func (m *model) View() string {
	if m.quitting {
		return "Bye Bye.\n"
	}

	pager, ok := m.pagers[m.page]
	if !ok {
		panic("unknown page")
	}
	return pager.view(m)
}
