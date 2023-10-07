package modules

import (
	"context"
	"crypto/sha512"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AuthPWHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	URI := os.Getenv("MONGODB_URI")
	if URI == "" {
		Critical(err)
	}
	form_email := XSSFix(r.FormValue("email"))
	form_pw := XSSFix(r.FormValue("password"))
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	Critical(err)
	defer func() {
		err := db.Disconnect(context.TODO())
		Critical(err)
	}()
	coll := db.Database("dj_users").Collection("users")
	encryptedPW := sha512.Sum512([]byte(form_pw)) //비밀 번호 해시 단방향 암호화
	filter := bson.D{{"email", form_email}, {"password", encryptedPW}}
	var dbres Dj_users_users
	err = coll.FindOne(context.TODO(), filter).Decode(&dbres)
	if err != nil { //로그인이 실패함
		log.Println("ID, 비밀번호 매칭 실패")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		redirect_msg := "<script>alert(\"로그인 실패\")</script><meta http-equiv=\"refresh\" content=\"0;url=/login/id/" + form_email + "\"></meta>" //다시 원래 pwrequst
		w.Write([]byte(redirect_msg))
	} else {
		//로그인이 성공함
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmp := []byte("<meta http-equiv=\"refresh\" content=\"0;url=/exit/\"></meta>")
		sessionkey := rand.Int()

		//http cookie에 세션키 저장
		cookieid := http.Cookie{
			Name:     "dj_session",
			Value:    strconv.Itoa(int(sessionkey)),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookieid)
		//db에 세션 클리어
		filter := bson.D{{"dj_user_id", dbres.ID}}
		coll_dj_session := db.Database("dj_users").Collection("sessions")
		result, err := coll_dj_session.DeleteMany(context.TODO(), filter)
		ErrOK(err)
		log.Println("session이 겹치는 이메일 삭제", result.DeletedCount)

		//db에 세션키 저장
		session_struct := Dj_user_session{
			Djuserid: dbres.ID,
			Session:  int(sessionkey),
			CreateAt: time.Now(),
		}
		result_1, err_1 := coll_dj_session.InsertOne(context.TODO(), session_struct)
		ErrOK(err_1)
		log.Println(result_1.InsertedID)
		w.Write([]byte(tmp))

		//users db에 last login 업데이트
		coll_dj_users := db.Database("dj_users").Collection("users")
		filter_users := bson.D{{"_id", dbres.ID}}
		update_users := bson.D{{"$set", bson.D{{"lastLogin", time.Now()}}}}
		_, err_users := coll_dj_users.UpdateOne(context.TODO(), filter_users, update_users)
		ErrOK(err_users)
	}
}
