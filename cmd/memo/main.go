package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/xuxinx/memo"
	"github.com/xuxinx/memo/file_dao"
)

const fPath = "./memo.json"

func main() {
	f := mustGetOrInitFile(fPath)
	defer f.Close()
	serv := mustGetMemoService(f)

Home:
	{
		clearScreen()
		p := promptui.Select{
			Label:    "Home",
			Items:    []string{"Practice", "New"},
			HideHelp: true,
		}
		_, r, err := p.Run()
		if err != nil {
			if err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				fmt.Println("Bye Bye.")
				return
			}
			panic(err)
		}
		switch r {
		case "New":
			goto New
		case "Practice":
			goto Practice
		default:
			panic("unknown case")
		}
	}
New:
	{
		var inputTags []string
		{
			tags, err := serv.GetTags()
			if err != nil {
				panic(err)
			}
		SelectTag:
			tagsItems := append([]string{"New Tags"}, tags...)
			sp := promptui.Select{
				Label:    "Tags",
				Items:    tagsItems,
				HideHelp: true,
				Searcher: func(input string, index int) bool {
					v := tagsItems[index]
					return strings.Contains(v, input)
				},
			}
			_, r, err := sp.Run()
			if err != nil {
				goto Home
			}
			switch r {
			case "New Tags":
				p := promptui.Prompt{
					Label: "Tags(seperate by comma)",
					Validate: func(s string) error {
						if s == "" {
							return errors.New("required")
						}
						return nil
					},
				}
				s, err := p.Run()
				if err != nil {
					goto Home
				}
				inputTags = strings.Split(s, ",")
			default:
				inputTags = append(inputTags, r)
			}

			sp = promptui.Select{
				Label:    fmt.Sprintf("Tags finished? %v", inputTags),
				Items:    []string{"Yes", "No"},
				HideHelp: true,
			}
			_, r, err = sp.Run()
			if err != nil {
				goto Home
			}
			if r != "Yes" {
				goto SelectTag
			}
		}
		p := promptui.Prompt{
			Label: "Question",
			Validate: func(s string) error {
				if s == "" {
					return errors.New("required")
				}
				return nil
			},
		}
		question, err := p.Run()
		if err != nil {
			goto Home
		}
		p = promptui.Prompt{
			Label: "Answer",
			Validate: func(s string) error {
				if s == "" {
					return errors.New("required")
				}
				return nil
			},
		}
		answer, err := p.Run()
		if err != nil {
			goto Home
		}
		p = promptui.Prompt{
			Label: "Desc",
		}
		desc, err := p.Run()
		if err != nil {
			goto Home
		}
		_, err = serv.NewQuestion(inputTags, question, desc, answer)
		if err != nil {
			panic(err)
		}
		goto Home
	}
Practice:
	{
		clearScreen()
		q, err := serv.GetTheNextReadyToPracticeQuestion()
		if err != nil {
			panic(err)
		}
		if q == nil {
			fmt.Println("no more questions need to practice.")
			p := promptui.Prompt{
				Label: "Back to Home",
			}
			p.Run()
			goto Home
		}
		fmt.Println(q.Tags, q.Question)
		if q.Desc != "" {
			fmt.Println(q.Desc)
		}
		p := promptui.Prompt{
			Label: "Show Answer",
		}
		_, err = p.Run()
		if err != nil {
			goto Home
		}
		fmt.Println(q.Answer)
		sp := promptui.Select{
			Label:    "Do you remember?",
			Items:    []string{"Remember", "Not Sure", "Forget"},
			HideHelp: true,
		}
		_, r, err := sp.Run()
		if err != nil {
			goto Home
		}
		switch r {
		case "Remember":
			err := serv.MarkQuestion(q.ID, memo.Mark_Remember)
			if err != nil {
				panic(err)
			}
		case "Not Sure":
			err := serv.MarkQuestion(q.ID, memo.Mark_NotSure)
			if err != nil {
				panic(err)
			}
		case "Forget":
			err := serv.MarkQuestion(q.ID, memo.Mark_Forget)
			if err != nil {
				panic(err)
			}
		default:
			panic("unknown case")
		}
		goto Practice
	}
}

func mustGetOrInitFile(p string) *os.File {
	if _, err := os.Stat(p); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
		f, err := os.Create(p)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString("[]")
		if err != nil {
			panic(err)
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(p, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	return f
}

func mustGetMemoService(f *os.File) *memo.Memo {
	dao, err := file_dao.NewDao(f)
	if err != nil {
		panic(err)
	}

	s, err := memo.NewMemo(dao)
	if err != nil {
		panic(err)
	}
	return s
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
