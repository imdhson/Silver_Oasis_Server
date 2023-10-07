package modules

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SessionTO_oid(w http.ResponseWriter, r *http.Request) primitive.ObjectID {
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
	coll_dj_registration := db.Database("gd_users").Collection("sessions")
	sessionkey, err := r.Cookie("dj_session") //key to value로 쿠키를 가져옴
	if err != nil {
		ErrOK(err)
		return primitive.NilObjectID
	}
	sessionkey_int, err := strconv.Atoi(sessionkey.Value)
	ErrOK(err)
	var dbres Dj_user_session
	filter := bson.D{{"dj_session", sessionkey_int}}
	err = coll_dj_registration.FindOne(context.TODO(), filter).Decode(&dbres)
	if err != nil {
		ErrOK(err)
		return primitive.NilObjectID
	}
	return dbres.Djuserid
}
