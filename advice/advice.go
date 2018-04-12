package advice

import (
	"github.com/globalsign/mgo"
)

var session *mgo.Session

type Repository interface {
	Init(host string)
	Create(advice *Advice) (error, *Advice)
	Search(term string) (error, *Advice)
	Random() (error, *Advice)
	Delete(advice *Advice) (error, *Advice)
}

type Advice struct {
	Advice string
	Keywords []string
	Funny int64
	Serious int64
}