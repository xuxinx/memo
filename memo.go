package memo

import (
	"errors"
	"fmt"
	"time"
)

type Memo struct {
	dao Dao
}

type Mark int

const (
	Mark_Remember = iota
	Mark_NotSure
	Mark_Forget
)

func NewMemo(dao Dao) (*Memo, error) {
	if dao == nil {
		return nil, errors.New("dao is nil")
	}
	return &Memo{
		dao: dao,
	}, nil
}

func (me *Memo) NewQuestion(
	tags []string,
	question string,
	desc string,
	answer string,
) (rq *Question, err error) {
	return me.dao.New(&Question{
		Tags:             tags,
		Question:         question,
		Desc:             desc,
		Answer:           answer,
		Score:            1,
		NextPracticeTime: time.Now(),
	})
}

func (me *Memo) GetQuestion(id uint) (rq *Question, err error) {
	return me.dao.Get(id)
}

func (me *Memo) GetTheNextReadyToPracticeQuestion() (rq *Question, err error) {
	return me.dao.GetTheNextReadyToPractice()
}

func (me *Memo) UpdateQuestion(id uint, q *Question) (err error) {
	return me.dao.Update(id, q)
}

func (me *Memo) DeleteQuestion(id uint) (err error) {
	return me.dao.Delete(id)
}

func (me *Memo) MarkQuestion(id uint, mark Mark) (err error) {
	q, err := me.dao.Get(id)
	if err != nil {
		return err
	}
	if q == nil || q.NextPracticeTime.Sub(time.Now()) > 0 {
		return nil
	}
	score := q.Score
	switch mark {
	case Mark_Remember:
		score *= 1.3
	case Mark_NotSure:
		if score < 1.3*1.3 {
			score = 1
		} else {
			score = 1.3 * 1.3
		}
	case Mark_Forget:
		score = 1
	default:
		return fmt.Errorf("unknown mark: %v", mark)
	}
	q.Score = score

	if mark == Mark_Remember {
		now := time.Now()
		yy := now.Year()
		mm := now.Month()
		dd := now.Day()
		increseDays := int(score) * int(score)
		if increseDays > 180 {
			increseDays = 180
		}
		q.NextPracticeTime = time.Date(yy, mm, dd+increseDays, 2, 0, 0, 0, time.Local)
	}

	return me.UpdateQuestion(id, q)
}
