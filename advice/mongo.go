package advice

import (
	"os"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"errors"
)

type MongoRepository struct {
	Session *mgo.Session
	Collection *mgo.Collection
}

func (repo *MongoRepository) Init(host string) (error, *MongoRepository){
	if host == "" {
		host = os.Getenv("MONGO_HOST")
	}
	s, err := mgo.Dial(host)
	if err != nil {
		return err, nil
	}
	repo.Session = s
	repo.Collection = s.DB("advice").C("advice")
	return err, repo
}

func (repo *MongoRepository) Create(advice *Advice) (error, *Advice) {
	err := repo.Collection.Insert(advice)
	if err != nil {
		return err, nil
	}
	return nil, advice
}

func (repo *MongoRepository) Search(term string) (error, *Advice) {
	var advice Advice
	pipe := repo.Collection.Pipe([]bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{"keywords": bson.M{"$regex": term, "$options": "i"}, },
					{"advice": bson.M{"$regex": term, "$options": "i"}, },
				},
			},
		},
		{
			"$sort": bson.M{
				"serious": -1,
			},
		},
		{
			"$limit": 1,
		},
	})

	if pipe.Iter().Next(&advice) {
		return nil, &advice
	}
	return errors.New("[Couldn't find any advice]"), nil
}

func (repo *MongoRepository) Random() (error, *Advice) {
	var advice Advice
	pipe := repo.Collection.Pipe([]bson.M{{
		"$sample": bson.M{
			"size": 1,
		},
	}})

	if pipe.Iter().Next(&advice) {
		return nil, &advice
	}
	return errors.New("[Couldn't find any advice]"), nil
}

func (repo *MongoRepository) Delete(advice *Advice) (error, *Advice) {
	return nil, nil
}