package modules

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ArticlesInsertHandler(w http.ResponseWriter, r *http.Request) {
	if !IsHeLogin(w, r) {
		ErrHandler(w, r)
		return
	}
	user_oid := SessionTO_oid(w, r)

	title := XSSFix(r.FormValue("title"))
	content := XSSFix(r.FormValue("content"))

	if title == "" || content == "" { //빈칸일 경우 무시
		return
	}

	var anon bool
	if r.FormValue("anonymous") == "on" {
		anon = true
	}
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
	coll := db.Database("dj_board").Collection("articles")
	articles_struct := Dj_board_articles{
		Djuserid: user_oid,
		CreateAt: time.Now(),
		Title:    title,
		Content:  content,
	}
	if anon {
		articles_struct.Djuserid = primitive.NilObjectID
	}
	result, err := coll.InsertOne(context.TODO(), articles_struct)
	if err != nil {
		ErrOK(err)
	} else {
		inserted_id := result.InsertedID.(primitive.ObjectID).Hex() //result.id를 hex로 변환
		tmp_urlpath := []string{"", inserted_id}
		ArticlesDetailHandler(w, r, &tmp_urlpath)
		go CallBard(result.InsertedID.(primitive.ObjectID), title, content) //바드 비공식 파이선 api를 부름 -  비동기로
	}
}
