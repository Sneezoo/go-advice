package main

import "github.com/gin-gonic/gin"
import (
	"os"
	"github.com/globalsign/mgo"
	"net/http"
	"github.com/Sneezoo/advicery/advice"
)

var collection *mgo.Collection
var repository *advice.MongoRepository

func main() {
	var err error
	r := gin.Default()
	mongoHost := os.Getenv("MONGO_HOST")
	err, repository = (&advice.MongoRepository{}).Init(mongoHost)
	defer repository.Session.Close()
	if err != nil {
		panic(err)
	}

	r.GET("/advice", func(context *gin.Context) {
		keyword := context.Query("term")

		if err, advice := repository.Search(keyword); err == nil {
			context.JSON(http.StatusOK, advice)
			return
		}
		if err, advice := repository.Random(); err == nil {
			context.JSON(http.StatusOK, advice)
			return
		}
		raiseError(http.StatusNotFound, "Couldn't find advice for keyword", err, context)
	})

	r.POST("/advice", func(context *gin.Context) {
		var advice advice.Advice

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
