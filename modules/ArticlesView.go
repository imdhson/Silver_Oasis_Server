package modules

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ArticlesView(w http.ResponseWriter, r *http.Request) {
	if !IsHeLogin(w, r) {
		ErrHandler(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	wwwfile, err := os.ReadFile("./www/articles.html")
	Critical(err)
	var htmlmodify Vars_on_html
	htmlmodify.Init()

	err = godotenv.Load()
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
	coll := db.Database("gd_board").Collection("articles")
	opts := options.Find().SetSort(bson.D{{"createAt", -1}})
	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	ErrOK(err)

	var will_send []Dj_board_articles
	for cursor.Next(context.TODO()) {
		var dbres Dj_board_articles
		err := cursor.Decode(&dbres)
		ErrOK(err)
		will_send = append(will_send, dbres)
	}
	ErrOK(err)

	var title_msg string

	for _, v := range will_send {
		useremail := OidTOuser_struct(v.Djuserid).Email
		if v.Djuserid == primitive.NilObjectID {
			useremail = "익명의 유저"
		} else {
			useremail, _, _ = strings.Cut(useremail, "@")
		}
		compare_time := time.Since(v.CreateAt).String()
		compare_time, _, _ = strings.Cut(compare_time, "m") // m 이후로 무시하기 위함
		if strings.Contains(compare_time, ".") {            //1분 미만이면 방금이라고 표기
			compare_time = "방금 "
		} else {
			compare_time += "분 " //1분 이상이면 숫자+분
		}
		compare_time = strings.ReplaceAll(compare_time, "h", "시간")
		compare_time = strings.ReplaceAll(compare_time, "d", "일")

		//댓글 개수 계산
		coll_for_commentCount := db.Database("gd_board").Collection("comments")
		filter_for_commentCount := bson.D{{"dj_jobs_id", v.ID}}
		commentCount, err := coll_for_commentCount.CountDocuments(context.TODO(), filter_for_commentCount)
		ErrOK(err)

		article_url := "/articles/" + v.ID.Hex()

		//그리드 작성 시작
		/*
					    <a href=%%article_url>
			    <div class="post">
			      <h2>%%v.Title</h2>
			      <p class="time"><i class="far fa-clock"></i>%%compare_time</p>
			      <p class="author"><i class="fas fa-user"></i>%%usermail</p>
			      <p class="commentCount"><i class="far fa-comments"></i>$$commentCount</p>
			    </div>
			    </a>
		*/
		title_msg += "<a href=\"" + article_url + "\">"
		title_msg += " <div class=\"post\"> 	 <h2>" + v.Title + "</h2>"
		title_msg += "<p class=\"time\"><i class=\"far fa-clock\"></i>" + compare_time + "전</p>"
		title_msg += "<p class=\"author\"><i class=\"fas fa-user\"></i>" + useremail + "</p>"
		title_msg += "<p class=\"commentCount\"><i class=\"far fa-comments\"></i>" + strconv.Itoa(int(commentCount)) + "개</p>"
		title_msg += "</div> 	</a>\n"
	}
	htmlmodify.AddVar("articles_grid", title_msg)
	html_modified := htmlmodify.VarsOnHTML(wwwfile)
	w.Write(html_modified)
}
