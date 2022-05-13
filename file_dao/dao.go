package file_dao

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/xuxinx/memo"
)

func NewDao(f *os.File) (dao memo.Dao, err error) {
	d := &fileDao{
		f:  f,
		qs: make(map[uint]*memo.Question),
	}

	if err := d.read(); err != nil {
		return nil, err
	}

	return d, nil
}

type fileDao struct {
	f  *os.File
	qs map[uint]*memo.Question
}

func (d *fileDao) read() (err error) {
	b, err := ioutil.ReadAll(d.f)
	if err != nil {
		return err
	}

	qss := []*memo.Question{}
	err = json.Unmarshal(b, &qss)
	if err != nil {
		return err
	}

	for i, _ := range qss {
		q := qss[i]
		d.qs[q.ID] = q
	}

	return nil
}

func (d *fileDao) save() (err error) {
	qss := make([]*memo.Question, 0, len(d.qs))
	for i, _ := range d.qs {
		qss = append(qss, d.qs[i])
	}
	sort.Slice(qss, func(i, j int) bool {
		return qss[i].ID < qss[j].ID
	})

	b, err := json.MarshalIndent(qss, "", "  ")
	if err != nil {
		return err
	}
	err = d.f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = d.f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = d.f.Write(b)
	return err
}

func (d *fileDao) nextID() uint {
	var maxID uint
	for id, _ := range d.qs {
		if id > maxID {
			maxID = id
		}
	}

	return maxID + 1
}

func (d *fileDao) New(q *memo.Question) (rq *memo.Question, err error) {
	q.ID = d.nextID()
	d.qs[q.ID] = q
	err = d.save()
	return q, err
}

func (d *fileDao) Get(id uint) (rq *memo.Question, err error) {
	return d.qs[id], nil
}

func (d *fileDao) GetTheNextReadyToPractice() (rq *memo.Question, err error) {
	now := time.Now()
	readys := make([]uint, 0, len(d.qs))
	for id, q := range d.qs {
		if q.NextPracticeTime.Sub(now) < 0 {
			readys = append(readys, id)
		}
	}
	if len(readys) == 0 {
		return nil, nil
	}
	rid := readys[rand.Intn(len(readys))]
	return d.qs[rid], nil
}

func (d *fileDao) Update(id uint, q *memo.Question) (err error) {
	d.qs[id] = q
	return d.save()
}

func (d *fileDao) Delete(id uint) (err error) {
	delete(d.qs, id)
	return d.save()
}

func (d *fileDao) GetTags() (tags []string, err error) {
	m := make(map[string]struct{})
	for _, q := range d.qs {
		for _, t := range q.Tags {
			m[t] = struct{}{}
		}
	}
	for t, _ := range m {
		tags = append(tags, t)
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i] < tags[j]
	})
	return tags, nil
}
