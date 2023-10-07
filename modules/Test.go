package modules

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SampleAIList(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	Critical(err)
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
	coll := db.Database("dj_jobs").Collection("job_list")

	//쿼리1
	filter := bson.D{
		{"$in", bson.A{
			bson.D{{"사업장명", "대한적십자사"}},
			bson.D{{"사업장명", "코오롱글로벌"}},
		},
		},
	}
	//filter := bson.D{{"사업장명", bson.D{{"$regex", "서울아산병원"}}}}
	var dbres_1 Dj_jobs_detail
	err = coll.FindOne(context.TODO(), filter).Decode(&dbres_1)
	ErrOK(err)
	dbres_1.AI_List_num = 1
	// 쿼리2
	filter = bson.D{{"사업장명", "용인시청"}}
	var dbres_2 Dj_jobs_detail
	err = coll.FindOne(context.TODO(), filter).Decode(&dbres_2)
	ErrOK(err)
	dbres_2.AI_List_num = 2
	//쿼리3
	filter = bson.D{{"사업장명", "강원랜드(주)"}}
	var dbres_3 Dj_jobs_detail
	err = coll.FindOne(context.TODO(), filter).Decode(&dbres_3)
	ErrOK(err)
	dbres_3.AI_List_num = 3
	//병합
	var will_send []Dj_jobs_detail
	will_send = append(will_send, dbres_1, dbres_2, dbres_3)
	will_send_json, _ := json.MarshalIndent(will_send, " ", "	")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(will_send_json)
}

func Test2(w http.ResponseWriter, r *http.Request) { //applist를 보낼 때 모델이 될 예정
	err := godotenv.Load()
	Critical(err)
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

	coll := db.Database("dj_jobs").Collection("job_list")
	//쿼리1
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"사업장 주소", bson.D{{"$regex", "서울특별시"}}}},
			bson.D{{"사업장 주소", bson.D{{"$regex", "중구"}}}},
			bson.D{{"필수부위", bson.D{{"$regex", ""}}}},
		},
		},
	}
	cursor, err := coll.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())
	var will_send []Dj_jobs_detail
	App_List_num := 0
	for cursor.Next(context.TODO()) {
		var dbres Dj_jobs_detail
		cursor.Decode(&dbres)
		dbres.AI_List_num = App_List_num
		will_send = append(will_send, dbres)
		App_List_num++
	}
	ErrOK(err)

	will_send_json, _ := json.MarshalIndent(will_send, " ", "	")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(will_send_json)
}

/* 앤드 연산 예제
filter := bson.D{
   {"$and",
      bson.A{
         bson.D{{"rating", bson.D{{"$gt", 7}}}},
         bson.D{{"rating", bson.D{{"$lte", 10}}}},
      },
   },
}

*/

func Test3(w http.ResponseWriter, r *http.Request) { //object id로 joblist 를 불러오는 예제.
	oid, _ := primitive.ObjectIDFromHex("648e92b1f2d0f84208c426f1")
	will_send, _ := OidTOjobDetail(oid)
	will_send_json, _ := json.MarshalIndent(will_send, " ", "	")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(will_send_json)
}
