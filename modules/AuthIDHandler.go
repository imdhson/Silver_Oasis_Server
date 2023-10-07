package modules

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AuthIDHandler(w http.ResponseWriter, r *http.Request) {

	err := godotenv.Load()
	URI := os.Getenv("MONGODB_URI")
	if URI == "" {
		Critical(err)
	}
	form_email := XSSFix(r.FormValue("email"))
	if !strings.Contains(form_email, "@") || !strings.Contains(form_email, ".") { //ID가 @와 .을 포함해야함.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		redirect_msg := "<script>alert(\"ID 형식이 올바르지 않음.\")</script><meta http-equiv=\"refresh\" content=\"0;url=/login/\"></meta>"
		w.Write([]byte(redirect_msg))
		return
	}
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	Critical(err)
	defer func() {
		err := db.Disconnect(context.TODO())
		Critical(err)
	}()

	coll := db.Database("dj_users").Collection("users")
	filter := bson.D{{"email", form_email}}
	var dbres Dj_users_users
	err = coll.FindOne(context.TODO(), filter).Decode(&dbres)
	ErrOK(err)
	same_mail_not_found_on_users := func(err error) bool { //같은 email을 찾았는지 판별하는 anonymous 함수
		return err != nil
	}(err)

	if same_mail_not_found_on_users { //같은 이메일 찾지 못하였을 때 - 회원가입 모드
		//디비 registration에 항목 있는지 확인
		coll := db.Database("dj_users").Collection("registration")
		filter := bson.D{{"email", form_email}}
		var dbres_regist Dj_users_registration
		err = coll.FindOne(context.TODO(), filter).Decode(&dbres_regist)

		same_mail_not_found_on_register := func(err error) bool { //같은 email을 찾았는지 판별하는 anonymous 함수
			return err != nil
		}(err)

		if same_mail_not_found_on_register { //regist 에서 같은 이메일 못찾았을 때

			key := SmtpSender(form_email, true) //메일 보내고
			log.Println("pseudo: 회원가입으로 이동하기", form_email, key)

			//디비 registration에 추가
			now_time := time.Now()
			regist_struct := Dj_users_registration{
				Email:        form_email,
				VerifyNumber: key,
				CreateAt:     now_time,
			}
			coll_dj_regist := db.Database("dj_users").Collection("registration")
			result, err := coll_dj_regist.InsertOne(context.TODO(), regist_struct)
			ErrOK(err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			redirect_msg := "<meta http-equiv=\"refresh\" content=\"0;url=/r/" + form_email + "\"></meta>"
			w.Write([]byte(redirect_msg))

			log.Println("db regist에 저장 완료", result)

		} else { //regist에서 같은 이메일이 있을 때 Regist DB에서 삭제
			filter := bson.D{{"email", form_email}}
			coll_dj_regist := db.Database("dj_users").Collection("registration")
			result, err := coll_dj_regist.DeleteMany(context.TODO(), filter)
			ErrOK(err)
			log.Println("regist에서 겹치는 이메일 삭제", result.DeletedCount)
			AuthIDHandler(w, r) //삭제하고 다시 호출해서 다시수행
		}
	} else { //같은 이메일 user에서 찾았을 때 - 로그인모드
		log.Println("pseudo:", dbres.Email, "을 E-Mail로 로그인하기")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		redirect_msg := "<meta http-equiv=\"refresh\" content=\"0;url=/login/id/" + form_email + "\"></meta>"
		w.Write([]byte(redirect_msg))
	}
}
