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
		qm: make(map[uint]*memo.Question),
	}

	if err := d.read(); err != nil {
		return nil, err
	}

	return d, nil
}

type fileDao struct {
	f      *os.File
	qm     map[uint]*memo.Question
	maxQID uint
}

type data struct {
	Questions []*memo.Question
	// TODO Tags
}

func (d *fileDao) read() (err error) {
	b, err := ioutil.ReadAll(d.f)
	if err != nil {
		return err
	}

	da := data{}
	err = json.Unmarshal(b, &da)
	if err != nil {
		return err
	}

	for i := range da.Questions {
		q := da.Questions[i]
		d.qm[q.ID] = q
		if q.ID > d.maxQID {
			d.maxQID = q.ID
		}
	}

	return nil
}

func (d *fileDao) save() (err error) {
	qs := make([]*memo.Question, 0, len(d.qm))
	for i := range d.qm {
		qs = append(qs, d.qm[i])
	}
	sort.Slice(qs, func(i, j int) bool {
		return qs[i].ID < qs[j].ID
	})

	b, err := json.MarshalIndent(&data{
		Questions: qs,
	}, "", "  ")
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
	d.maxQID++
	return d.maxQID
}

func (d *fileDao) New(q *memo.Question) (rq *memo.Question, err error) {
	q.ID = d.nextID()
	d.qm[q.ID] = q
	err = d.save()
	return q, err
}

func (d *fileDao) Get(id uint) (rq *memo.Question, err error) {
	return d.qm[id], nil
}

func (d *fileDao) GetAll() (rqs []*memo.Question, err error) {
	rqs = make([]*memo.Question, 0, len(d.qm))
	for k := range d.qm {
		rqs = append(rqs, d.qm[k])
	}
	sort.Slice(rqs, func(i, j int) bool {
		return rqs[i].ID > rqs[j].ID
	})
	return rqs, nil
}

func (d *fileDao) GetTheNextReadyToPractice() (rq *memo.Question, err error) {
	now := time.Now()
	readys := make([]uint, 0, len(d.qm))
	for id, q := range d.qm {
		if q.NextPracticeTime.Sub(now) < 0 {
			readys = append(readys, id)
		}
	}
	if len(readys) == 0 {
		return nil, nil
	}
	rid := readys[rand.Intn(len(readys))]
	return d.qm[rid], nil
}

func (d *fileDao) Update(id uint, q *memo.Question) (err error) {
	d.qm[id] = q
	return d.save()
}

func (d *fileDao) Delete(id uint) (err error) {
	delete(d.qm, id)
	return d.save()
}
