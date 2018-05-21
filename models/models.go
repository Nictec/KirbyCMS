package models 

import (
	"gopkg.in/mgo.v2/bson"
)

type Video struct {
	Id bson.ObjectId `bson:"_id", json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	CatId string `json:"cat_id"`
	Privacy string `json:"privacy"`
	YtId string `json:"yt_id"`
	File string `json:"file"`
	Status string `json:"status"`
}