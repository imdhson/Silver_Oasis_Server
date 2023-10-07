package modules

import (
	"context"
	"crypto/sha512"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	URI := os.Getenv("MONGODB_URI")
	if URI == "" {
		Critical(err)
	}
	form_key := XSSFix(r.FormValue("verifyNumber"))
	form_pw1 := XSSFix(r.FormValue("password1"))
	form_pw2 := XSSFix(r.FormValue("password2"))
	if form_pw1 != form_pw2 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		msg := "<script>alert(\"비밀번호가 불일치합니다. 다시 입력해주세요.\")</script><meta http-equiv=\"refresh\" content=\"0;url=/login/\"></meta>"
		w.Write([]byte(msg))
		return
	}

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	Critical(err)
	defer func() {
		err := db.Disconnect(context.TODO())
		Critical(err)
	}()
	coll_dj_registration := db.Database("gd_users").Collection("registration")
	var dbres Dj_users_registration
	filter_for_key_email_search := bson.D{{"verifyNumber", string(form_key)}}
	err = coll_dj_registration.FindOne(context.TODO(), filter_for_key_email_search).Decode(&dbres) //key로 email 찾기
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		msg := "<script>alert(\"키를 찾기 못함.\")</script><meta http-equiv=\"refresh\" content=\"0;url=/login/\"></meta>"
		w.Write([]byte(msg))
		return
	}
	encryptedPW := sha512.Sum512([]byte(form_pw1)) //비밀 번호 해시 단방향 암호화
	users_struct := Dj_users_users{
		Email: dbres.Email, Password: encryptedPW,
		LastLogin: time.Now(), ScrapList: []primitive.ObjectID{}, // primitive.NewObjectID 로 나중에 Push 가능
		Settings: Dj_users_users_settings{},
	}
	coll_dj_users := db.Database("gd_users").Collection("users")
	result, err := coll_dj_users.InsertOne(context.TODO(), users_struct)
	log.Println(result)

	//사용된 key 삭제
	filter := bson.D{{"verifyNumber", form_key}}
	coll_dj_regist := db.Database("gd_users").Collection("registration")
	result_del, err := coll_dj_regist.DeleteMany(context.TODO(), filter)
	ErrOK(err)
	log.Println("regist 사용된 키 삭제:", result_del.DeletedCount)

	Immed_Login_AfterRegister(w, r, users_struct.Email, users_struct.Password)
}
