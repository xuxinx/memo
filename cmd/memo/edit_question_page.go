package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xuxinx/memo"
)

const pageEditQuestion page = "Edit Question"

const (
	editQuestionFidxTagsInput = iota
	editQuestionFidxQuestionInput
	editQuestionFidxDescInput
	editQuestionFidxAnswerInput
	maxEditQuestionFidxForInputs = iota - 1
	editQuestionFidxSubmitBtn
	maxEditQuestionFidx = iota - 2
)

type editQuestionPage struct {
	fidx   int
	inputs []textinput.Model

	q        *memo.Question
	backPage page
}

func initEditQuestionPage(
	q *memo.Question,
	backPage page,
) pager {
	if q == nil {
		return nil
	}

	tagsIn := textinput.New()
	tagsIn.CursorStyle = cursorStyle
	tagsIn.Placeholder = "Tags"
	tagsIn.SetCursorMode(textinput.CursorStatic)
	tagsIn.SetValue(strings.Join(q.Tags, ","))

	qIn := textinput.New()
	qIn.CursorStyle = cursorStyle
	qIn.Placeholder = "Question"
	qIn.SetCursorMode(textinput.CursorStatic)
	qIn.SetValue(q.Question)

	descIn := textinput.New()
	descIn.CursorStyle = cursorStyle
	descIn.Placeholder = "Desc"
	descIn.SetCursorMode(textinput.CursorStatic)
	descIn.SetValue(q.Desc)

	answerIn := textinput.New()
	answerIn.CursorStyle = cursorStyle
	answerIn.Placeholder = "Answer"
	answerIn.SetCursorMode(textinput.CursorStatic)
	answerIn.SetValue(q.Answer)

	tagsIn.Focus()
	qIn.Blur()
	descIn.Blur()
	answerIn.Blur()

	return &editQuestionPage{
		inputs: []textinput.Model{
			tagsIn, qIn, descIn, answerIn,
		},
		q:        q,
		backPage: backPage,
	}
}

func (p *editQuestionPage) update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			if key == "q" && p.fidx <= maxEditQuestionFidxForInputs && p.inputs[p.fidx].CursorMode() == textinput.CursorBlink {
				break
			}
			m.page = p.backPage
			return m, nil
		case "i":
			if p.fidx <= maxEditQuestionFidxForInputs && p.inputs[p.fidx].CursorMode() == textinput.CursorStatic {
				cmd := p.inputs[p.fidx].SetCursorMode(textinput.CursorBlink)
				p.inputs[p.fidx].CursorMode()
				return m, cmd
			}
		case "ctrl+[", "esc":
			if p.fidx <= maxEditQuestionFidxForInputs && p.inputs[p.fidx].CursorMode() == textinput.CursorBlink {
				cmd := p.inputs[p.fidx].SetCursorMode(textinput.CursorStatic)
				return m, cmd
			}
			return m, nil
		case "j":
			if p.fidx > maxEditQuestionFidxForInputs || p.inputs[p.fidx].CursorMode() == textinput.CursorStatic {
				p.fidx++
				if p.fidx > maxEditQuestionFidx {
					p.fidx = maxEditQuestionFidx
				}
				if p.fidx <= maxEditQuestionFidxForInputs+1 {
					p.inputs[p.fidx-1].Blur()
				}
				if p.fidx <= maxEditQuestionFidxForInputs {
					p.inputs[p.fidx].Focus()
				}

				return m, nil
			}
		case "k":
			if p.fidx > maxEditQuestionFidxForInputs || p.inputs[p.fidx].CursorMode() == textinput.CursorStatic {
				p.fidx--
				if p.fidx < 0 {
					p.fidx = 0
				}
				p.inputs[p.fidx].Focus()
				if p.fidx+1 <= maxEditQuestionFidxForInputs {
					p.inputs[p.fidx+1].Blur()
				}
				return m, nil
			}
		case "enter":
			if p.fidx <= maxEditQuestionFidxForInputs && p.inputs[p.fidx].CursorMode() == textinput.CursorBlink {
				v := p.inputs[p.fidx].Value()
				p.inputs[p.fidx].Reset()
				p.inputs[p.fidx].SetValue(v + enterChar)
			}
			if p.fidx == editQuestionFidxSubmitBtn {
				tags := []string{}
				if v := p.inputs[0].Value(); v != "" {
					tags = strings.Split(v, ",")
				}
				p.q.Tags = tags
				p.q.Question = p.inputs[1].Value()
				p.q.Desc = p.inputs[2].Value()
				p.q.Answer = strings.Replace(p.inputs[3].Value(), enterChar, "\n", -1)
				if p.q.ID == 0 {
					_, err := serv.NewQuestion(tags, p.q.Question, p.q.Desc, p.q.Answer)
					if err != nil {
						panic(err)
					}
				} else {
					err := serv.UpdateQuestion(
						p.q.ID,
						p.q,
					)
					if err != nil {
						panic(err)
					}
				}
				m.page = p.backPage
			}
			return m, nil
		case "ctrl+h":
			m.page = pageHelp
			m.pagers[m.page] = initHelpPage(
				`
  ctrl+c / q: back to previous page
       j / k: move
       enter: select item 
           i: enter input mode
ctrl+[ / esc: exit input mode
                `,
				pageEditQuestion,
			)
			return m, nil
		default:
			if p.fidx <= maxEditQuestionFidxForInputs && p.inputs[p.fidx].CursorMode() == textinput.CursorStatic {
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

func (p *editQuestionPage) view(m *model) string {
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
	if p.fidx == editQuestionFidxSubmitBtn {
		btn = focusStyle.Render(btn)
	}
	b.WriteString("\n")
	b.WriteString(btn)
	b.WriteString("\n")
	return b.String()
}
