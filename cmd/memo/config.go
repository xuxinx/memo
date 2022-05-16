package main

import (
	"errors"
	"os"

	"github.com/xuxinx/memo"
	"github.com/xuxinx/memo/file_dao"
)

func mustGetOrInitFile(p string) *os.File {
	if _, err := os.Stat(p); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
		f, err := os.Create(p)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString("{}")
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
