package main

import "github.com/gin-gonic/gin"
import (
	"os"
	"github.com/globalsign/mgo"
	"net/http"
	"github.com/globalsign/mgo/bson"
	"errors"
	"github.com/Sneezoo/testproject/advice"
)

type Advice struct {
	Advice string
	Keywords []string
	Funny int64
	Serious int64
}

var collection *mgo.Collection
var repo *advice.MongoRepository

func main() {
	r := gin.Default()
	mongoHost := os.Getenv("MONGO_HOST")
	session, err := mgo.Dial(mongoHost)
	defer session.Close()

	if err != nil {
		panic(err)
	}

	db := session.DB("advice")
	collection = db.C("advice")

	r.GET("/advice", func(context *gin.Context) {
		keyword := context.Query("term")
		err, advice := getAdvice(keyword)
		if err != nil {
			raiseError(http.StatusNotFound, "Couldn't find advice for keyword", err, context)
			return
		}
		context.JSON(http.StatusOK, *advice)
	})

	r.POST("/advice", func(context *gin.Context) {
		var advice Advice

		err := context.Bind(&advice)
		if err != nil {
			raiseError(http.StatusBadRequest, "Couldn't parse Request", err, context)
			return
		}

		err = collection.Insert(&advice)
		if err != nil {
			raiseError(http.StatusInternalServerError, "Couldn't save to database", err, context)
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"status": "Added new advice",
			"advice": advice,
		})
	})

	r.Run(":8080")
}

func getAdvice(keyword string) (error, *Advice) {
	var advice Advice
	pipe := collection.Pipe([]bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{"keywords": bson.M{"$regex": keyword, "$options": "i"}, },
					{"advice": bson.M{"$regex": keyword, "$options": "i"}, },
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
	pipeSample := collection.Pipe([]bson.M{{
		"$sample": bson.M{
			"size": 1,
		},
	}})

	if pipe.Iter().Next(&advice) {
		return nil, &advice
	} else if pipeSample.Iter().Next(&advice) {
		return nil, &advice
	}
	return errors.New("[Couldn't find any advice]"), nil
}

func raiseError(code int16, status string, err error, ctx *gin.Context) {
	errorString := ""
	if err != nil {
		errorString = err.Error()
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code": code,
		"status": status,
		"error": errorString,
	})
}
