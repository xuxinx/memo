package memo

import "time"

type Question struct {
	ID uint

	Tags     []string
	Question string
	Desc     string
	Answer   string

	Score            float64
	NextPracticeTime time.Time
}
