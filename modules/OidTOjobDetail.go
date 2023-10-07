package modules

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func OidTOjobDetail(oid primitive.ObjectID) (SO_jobs_detail, error) {
	err := godotenv.Load()
	URI := os.Getenv("MONGODB_URI")
	if URI == "" {
		Critical(err)
	}
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	Critical(err)
	defer func() {
		err := db.Disconnect(context.TODO())
		Critical(err)
	}()
	coll := db.Database("gd_facilities").Collection("gd_fac_list")
	var dbres SO_jobs_detail
	filter := bson.D{{"_id", oid}}
	err = coll.FindOne(context.TODO(), filter).Decode(&dbres)
	return dbres, err
}
