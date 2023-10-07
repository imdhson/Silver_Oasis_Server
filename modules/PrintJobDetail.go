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

type temp_detail struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	TypeofFacility      string             `bson:"시설종류" json:"시설종류"`
	Loc_1               string             `bson:"지역구분1" json:"지역구분1"`
	Loc_2               string             `bson:"지역구분2" json:"지역구분2"`
	NameofFacility      string             `bson:"시설명" json:"시설명"`
	OwnerofFacility     string             `bson:"시설장명" json:"시설장명"`
	Max_Cap             int32              `bson:"입소-정원" json:"입소-정원"`
	Now_Cap             int32              `bson:"현원-계" json:"현원-계"`
	Now_Cap_Male        int32              `bson:"현원-남" json:"현원-남"`
	Now_Cap_Female      int32              `bson:"현원-여" json:"현원-여"`
	Now_Employee        int32              `bson:"종사자-계" json:"종사자-계"`
	Now_Employee_Male   int32              `bson:"종사자-남" json:"종사자-남"`
	Now_Employee_Female int32              `bson:"종사자-여" json:"종사자-여"`
	Address             string             `bson:"소재지" json:"소재지"`
	Contact             string             `bson:"전화번호" json:"전화번호"`
	InitialDate         string             `bson:"설치일" json:"설치일"`
	Operator            string             `bson:"운영주체" json:"운영주체"`
	ViewCount           int                `bson:"viewCount" json:"viewCount"`
}

func PrintJobDetail(w http.ResponseWriter, r *http.Request, urlPath *[]string) {
	oid_hex := (*urlPath)[1]
	oid, err := primitive.ObjectIDFromHex(oid_hex)
	ErrOK(err)
	so_temp, err := OidTOjobDetail(oid)
	temp := temp_detail{
		ID:             so_temp.ID,
		TypeofFacility: so_temp.TypeofFacility,
		Loc_1:          so_temp.Loc_1,
		Loc_2:          so_temp.Loc_2, NameofFacility: so_temp.NameofFacility,
		OwnerofFacility:     so_temp.OwnerofFacility,
		Max_Cap:             so_temp.Max_Cap,
		Now_Cap:             so_temp.Now_Cap,
		Now_Cap_Male:        so_temp.Now_Cap_Male,
		Now_Cap_Female:      so_temp.Now_Cap_Female,
		Now_Employee:        so_temp.Now_Employee,
		Now_Employee_Male:   so_temp.Now_Cap_Male,
		Now_Employee_Female: so_temp.Now_Cap_Female,
		Address:             so_temp.Address,
		Contact:             so_temp.Contact,
		InitialDate:         so_temp.InitialDate,
		Operator:            so_temp.Operator,
		ViewCount:           so_temp.ViewCount,
	}

	// fac-list 목록에 viewCount 추가
	func(sid primitive.ObjectID) {
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
		coll_for_scrapCount := db.Database("gd_facilities").Collection("gd_fac_list")
		filter_for_scrapCount := bson.D{{"_id", sid}}
		update_for_scrapCount := bson.D{
			{"$inc", bson.D{{"viewCount", 1}}},
		}
		_, err = coll_for_scrapCount.UpdateOne(context.TODO(), filter_for_scrapCount, update_for_scrapCount)
		ErrOK(err)

	}(oid)

	if err != nil { //job을 찾지 못하였을 때
		temp := map[string]string{"시설명": "찾지 못함"}
		temp2, err := json.MarshalIndent(temp, " ", "	")
		ErrOK(err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(temp2)
	} else {
		temp2, err := json.MarshalIndent(temp, " ", "	")
		ErrOK(err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(temp2)
	}
}
