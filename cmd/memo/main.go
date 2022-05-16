package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xuxinx/memo"
)

const (
	fPath     = "./memo.json"
	enterChar = "¬"
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

type model struct {
	page     page
	quitting bool

	homePage        *homePage
	newQuestionPage *newQuestionPage
	practicePage    *practicePage
}

func initModel() *model {
	return &model{
		page:     pageHome,
		quitting: false,

		homePage:        initHomePage(),
		practicePage:    initPracticePage(false),
		newQuestionPage: initNewQuestionPage(),
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.page {
	case pageHome:
		return m.homePage.update(m, msg)
	case pagePractice:
		return m.practicePage.update(m, msg)
	case pageNewQuestion:
		return m.newQuestionPage.update(m, msg)
	}
	panic("unknown page Update")
}

func (m *model) View() string {
	if m.quitting {
		return "Bye Bye.\n"
	}
	switch m.page {
	case pageHome:
		return m.homePage.view()
	case pagePractice:
		return m.practicePage.view()
	case pageNewQuestion:
		return m.newQuestionPage.view()
	}

	panic("unknown page View")
}
