package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xuxinx/memo"
)

const pagePractice page = "Practice"

const (
	practiceFidxRememberBtn = iota
	practiceFidxNotSureBtn
	practiceFidxForgetBtn
	maxPracticeFidx = iota - 1
)

type practicePage struct {
	fidx         int
	answerShowed bool

	q *memo.Question
}

func initPracticePage(loadQuestion bool) pager {
	p := &practicePage{}
	if loadQuestion {
		q, err := serv.GetTheNextReadyToPracticeQuestion()
		if err != nil {
			panic(err)
		}
		p.q = q
	}
	return p
}

func (p *practicePage) update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			m.page = pageHome
			return m, nil
		case "j":
			if p.answerShowed {
				p.fidx++
				if p.fidx > maxPracticeFidx {
					p.fidx = maxPracticeFidx
				}
			}
			return m, nil
		case "k":
			if p.answerShowed {
				p.fidx--
				if p.fidx < 0 {
					p.fidx = 0
				}
			}
			return m, nil
		case "enter":
			if p.q == nil {
				m.page = pageHome
				return m, nil
			}
			if !p.answerShowed {
				p.fidx = 0
				p.answerShowed = true
				return m, nil
			}

			q := p.q
			switch p.fidx {
			case practiceFidxRememberBtn:
				err := serv.MarkQuestion(q.ID, memo.Mark_Remember)
				if err != nil {
					panic(err)
				}
			case practiceFidxNotSureBtn:
				err := serv.MarkQuestion(q.ID, memo.Mark_NotSure)
				if err != nil {
					panic(err)
				}
			case practiceFidxForgetBtn:
				err := serv.MarkQuestion(q.ID, memo.Mark_Forget)
				if err != nil {
					panic(err)
				}
			default:
				panic("unknown action")
			}
			q, err := serv.GetTheNextReadyToPracticeQuestion()
			if err != nil {
				panic(err)
			}
			p.q = q
			p.fidx = 0
			p.answerShowed = false
			return m, nil
		case "ctrl+h":
			m.page = pageHelp
			m.pagers[m.page] = initHelpPage(
				`
ctrl+c / q: back to previous page
     j / k: move
     enter: select item 
                `,
				pagePractice,
			)
			return m, nil
		}
	}

	return m, nil
}

func (p *practicePage) view(m *model) string {
	var b strings.Builder
	q := p.q
	if q == nil {
		b.WriteString("no more questions to practice. ")
		b.WriteString("\n")
		return b.String()
	}
	b.WriteString(fmt.Sprint(q.Tags))
	b.WriteString(" ")
	b.WriteString(q.Question)
	b.WriteString("\n")
	if q.Desc != "" {
		b.WriteString(q.Desc)
		b.WriteString("\n")
	}
	if p.answerShowed {
		b.WriteString(q.Answer)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	if !p.answerShowed {
		b.WriteString(focusStyle.Render("[ Show Answer ]"))
	} else {
		rBtn := "[ Remember ]"
		if p.fidx == practiceFidxRememberBtn {
			rBtn = focusStyle.Render(rBtn)
		}
		nBtn := "[ Not Sure ]"
		if p.fidx == practiceFidxNotSureBtn {
			nBtn = focusStyle.Render(nBtn)
		}
		fBtn := "[ Forget ]"
		if p.fidx == practiceFidxForgetBtn {
			fBtn = focusStyle.Render(fBtn)
		}
		b.WriteString(rBtn)
		b.WriteString(" ")
		b.WriteString(nBtn)
		b.WriteString(" ")
		b.WriteString(fBtn)
	}
	b.WriteString("\n")
	return b.String()
}
