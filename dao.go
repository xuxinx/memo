package memo

type Dao interface {
	New(q *Question) (rq *Question, err error)
	Get(id uint) (rq *Question, err error)
	GetAll() (rqs []*Question, err error)
	GetTheNextReadyToPractice() (rq *Question, err error)
	Update(id uint, q *Question) (err error)
	Delete(id uint) (err error)
}
