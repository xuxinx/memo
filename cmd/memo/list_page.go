package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xuxinx/memo"
)

const pageList page = "List"

type listPage struct {
	list list.Model
}

func initListPage(load bool) pager {
	items := []list.Item{}
	if load {
		qs, err := serv.GetAllQuestions()
		if err != nil {
			panic(err)
		}
		items = make([]list.Item, 0, len(qs))
		for i := range qs {
			q := qs[i]
			items = append(items, listItem{
				id:          q.ID,
				title:       q.Question,
				filterValue: strings.Join([]string{q.Question, q.Desc, q.Answer}, " "),
				value:       q,
			})
		}
	}
	return &listPage{
		list: newList(
			"",
			items,
		),
	}
}

func (p *listPage) update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			if p.list.FilterState() == list.Filtering {
				break
			}
			m.page = pageHome
			return m, nil
		case "enter":
			if p.list.FilterState() == list.Filtering {
				break
			}
			q := p.list.SelectedItem().(listItem).value.(*memo.Question)
			m.page = pageEditQuestion
			m.pagers[pageEditQuestion] = initEditQuestionPage(q, pageList)
			return m, nil
		case "x":
			if p.list.FilterState() == list.Filtering {
				break
			}
			q := p.list.SelectedItem().(listItem).value.(*memo.Question)
			if err := serv.DeleteQuestion(q.ID); err != nil {
				panic(err)
			}
			p.list.RemoveItem(p.list.Index())
			return m, nil
		case "ctrl+h":
			m.page = pageHelp
			m.pagers[m.page] = initHelpPage(
				`
ctrl+c / q: back to previous page
     j / k: move
     enter: select item 
         x: delete item
         /: search item
       esc: quit search
                `,
				pageList,
			)
			return m, nil
		}
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return m, cmd
}

func (p *listPage) view(m *model) string {
	return p.list.View()
}
