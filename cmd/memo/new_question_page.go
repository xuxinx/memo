package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const pageNewQuestion page = "New Question"

const (
	newQuestionFidxTagsInput = iota
	newQuestionFidxQuestionInput
	newQuestionFidxDescInput
	newQuestionFidxAnswerInput
	maxNewQuestionFidxForInputs = iota - 1
	newQuestionFidxSubmitBtn
	maxNewQuestionFidx = iota - 2
)

type newQuestionPage struct {
	fidx   int
	inputs []textinput.Model
}

func initNewQuestionPage() *newQuestionPage {
	tagsIn := textinput.New()
	tagsIn.CursorStyle = cursorStyle
	tagsIn.Placeholder = "Tags"
	tagsIn.Focus()
	tagsIn.SetCursorMode(textinput.CursorStatic)

	qIn := textinput.New()
	qIn.CursorStyle = cursorStyle
	qIn.Placeholder = "Question"
	qIn.SetCursorMode(textinput.CursorStatic)

	descIn := textinput.New()
	descIn.CursorStyle = cursorStyle
	descIn.Placeholder = "Desc"
	descIn.SetCursorMode(textinput.CursorStatic)

	answerIn := textinput.New()
	answerIn.CursorStyle = cursorStyle
	answerIn.Placeholder = "Answer"
	answerIn.SetCursorMode(textinput.CursorStatic)

	return &newQuestionPage{
		inputs: []textinput.Model{
			tagsIn, qIn, descIn, answerIn,
		},
	}
}

func (p *newQuestionPage) update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			if key == "q" && p.fidx <= maxNewQuestionFidxForInputs && p.inputs[p.fidx].CursorMode() == textinput.CursorBlink {
				break
			}
			m.page = pageHome
			return m, nil
		case "i":
			if p.fidx <= maxNewQuestionFidxForInputs && p.inputs[p.fidx].CursorMode() == textinput.CursorStatic {
				cmd := p.inputs[p.fidx].SetCursorMode(textinput.CursorBlink)
				p.inputs[p.fidx].CursorMode()
				return m, cmd
			}
		case "ctrl+[", "esc":
			if p.fidx <= maxNewQuestionFidxForInputs && p.inputs[p.fidx].CursorMode() == textinput.CursorBlink {
				cmd := p.inputs[p.fidx].SetCursorMode(textinput.CursorStatic)
				return m, cmd
			}
			return m, nil
		case "j":
			if p.fidx > maxNewQuestionFidxForInputs || p.inputs[p.fidx].CursorMode() == textinput.CursorStatic {
				p.fidx++
				if p.fidx > maxNewQuestionFidx {
					p.fidx = maxNewQuestionFidx
				}
				if p.fidx <= maxNewQuestionFidxForInputs+1 {
					p.inputs[p.fidx-1].Blur()
				}
				if p.fidx <= maxNewQuestionFidxForInputs {
					p.inputs[p.fidx].Focus()
				}

				return m, nil
			}
		case "k":
			if p.fidx > maxNewQuestionFidxForInputs || p.inputs[p.fidx].CursorMode() == textinput.CursorStatic {
				p.fidx--
				if p.fidx < 0 {
					p.fidx = 0
				}
				p.inputs[p.fidx].Focus()
				if p.fidx+1 <= maxNewQuestionFidxForInputs {
					p.inputs[p.fidx+1].Blur()
				}
				return m, nil
			}
		case "enter":
			if p.fidx <= maxNewQuestionFidxForInputs {
				v := p.inputs[p.fidx].Value()
				p.inputs[p.fidx].Reset()
				p.inputs[p.fidx].SetValue(v + enterChar)
			}
			if p.fidx == newQuestionFidxSubmitBtn {
				tags := []string{}
				if v := p.inputs[0].Value(); v != "" {
					tags = strings.Split(v, ",")
				}
				_, err := serv.NewQuestion(
					tags,
					p.inputs[1].Value(),
					p.inputs[2].Value(),
					strings.Replace(p.inputs[3].Value(), enterChar, "\n", -1),
				)
				if err != nil {
					panic(err)
				}
				m.page = pageHome
			}
			return m, nil
		default:
			if p.fidx <= maxNewQuestionFidxForInputs && p.inputs[p.fidx].CursorMode() == textinput.CursorStatic {
				return m, nil
			}
		}
	}

	cmds := make([]tea.Cmd, len(p.inputs))
	for i := range p.inputs {
		p.inputs[i], cmds[i] = p.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (p *newQuestionPage) view() string {
	var b strings.Builder
	for i := range p.inputs {
		if i == p.fidx {
			b.WriteString(focusStyle.Render(p.inputs[i].View()))
		} else {
			b.WriteString(p.inputs[i].View())
		}
		b.WriteString("\n")
	}
	btn := "[ Submit ]"
	if p.fidx == newQuestionFidxSubmitBtn {
		btn = focusStyle.Render(btn)
	}
	b.WriteString("\n")
	b.WriteString(btn)
	b.WriteString("\n")
	return b.String()
}
