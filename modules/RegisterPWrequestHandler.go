package modules

import (
	"context"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterPWrequestHandler(w http.ResponseWriter, r *http.Request, urlPath *[]string) {
	url_email := (*urlPath)[1]
	url_key := (*urlPath)[2]
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
	coll := db.Database("dj_users").Collection("registration")
	filter := bson.D{{"email", url_email}}
	var dbres Dj_users_registration
	err = coll.FindOne(context.TODO(), filter).Decode(&dbres)

	ErrOK(err)
	if dbres.VerifyNumber == url_key { //키가 맞으면
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		wwwfile, err := os.ReadFile("./www/register.html")
		Critical(err)
		var i Vars_on_html
		i.Init()
		i.AddVar("url_email", url_email)
		i.AddVar("url_key", url_key)
		i.VarsOnHTML(wwwfile)
		w.Write(i.VarsOnHTML(wwwfile))
	} else if url_key == "" { //키가 비었으면 - 기본 빈 폼으로 반환시켜줌
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		wwwfile, err := os.ReadFile("./www/register.html")
		Critical(err)
		var i Vars_on_html
		i.Init()
		i.AddVar("url_email", url_email)
		i.AddVar("url_key", "")
		i.VarsOnHTML(wwwfile)
		w.Write(i.VarsOnHTML(wwwfile))
	} else { //키가 다르면
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		redirect_msg := "<script>alert(\"키가 다릅니다. 다시시도해주세요.\")</script><meta http-equiv=\"refresh\" content=\"0;url=/login\">"
		w.Write([]byte(redirect_msg))
	}

}
