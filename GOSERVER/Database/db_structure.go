package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReqInfo struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	RequestUrl    string             `bson:"requestUrl"`
	RequestType   string             `bson:"requestType"`
	Timestamp     time.Time          `bson:"timestamp"`
	RequestStatus int                `bson:"requestStatus"`
	RequestDataID primitive.ObjectID `bson:"requestDataID"`
	ResponseTime  int                `bson:"responseTime"`
}

type ReqData struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Params   any                `bson:"params"`
	Response any                `bson:"response"`
}
