package modules

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dj_users_users struct {
	ID        primitive.ObjectID      `bson:"_id,omitempty"`
	Email     string                  `bson:"email"`
	Password  [64]byte                `bson:"password"`
	LastLogin time.Time               `bson:"lastLogin"`
	ScrapList []primitive.ObjectID    `bson:"scrapList"`
	Settings  Dj_users_users_settings `bson:"settings"`
}

type Dj_users_users_settings struct {
	Loc   string `bson:"loc"`
	Type1 string `bson:"type1"`
	Type2 string `bson:"type2"`
	Type3 string `bson:"type3"`
}

type Dj_users_registration struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email"`
	VerifyNumber string             `bson:"verifyNumber"`
	CreateAt     time.Time          `bson:"createAt"`
}

type Dj_user_session struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Djuserid primitive.ObjectID `bson:"dj_user_id"`
	Session  int                `bson:"dj_session"`
	CreateAt time.Time          `bson:"createAt"`
}

type GD_jobs_detail struct {
	AI_List_num         int
	AI_List_score       int
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
}

type GD_jobs_refined struct {
	AI_List_num    int
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	NameofFacility string             `bson:"시설명" json:"시설명"`
	TypeofFacility string             `bson:"시설종류" json:"시설종류"`
	Address        string             `bson:"소재지" json:"소재지"`
	Operator       string             `bson:"운영주체" json:"운영주체"`
}

type GD_jobs_detail_s []GD_jobs_detail

func (a GD_jobs_detail_s) will_send_append(i GD_jobs_detail, score int) {
	for _, v := range a {
		if v.ID == i.ID {
			v.AI_List_score += score
			return
		}
	}
	i.AI_List_score += score
	a = append(a, i)
	return
}

func (a GD_jobs_detail_s) Len() int {
	return len(a)
}

func (a GD_jobs_detail_s) Swap(i int, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a GD_jobs_detail_s) Less(i int, j int) bool {
	return a[i].AI_List_score < a[j].AI_List_score
}

type Dj_board_comments struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Djjobsid primitive.ObjectID `bson:"dj_jobs_id"`
	Djuserid primitive.ObjectID `bson:"dj_user_id"`
	CreateAt time.Time          `bson:"createAt"`
	Content  string             `bson:"content"`
	GenbyAI  bool               `bson:"genbyAI"`
}
type Dj_board_articles struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Djuserid primitive.ObjectID `bson:"dj_user_id"`
	CreateAt time.Time          `bson:"createAt"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}
type GD_service_type struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Type           string             `bson:"시설종류"`
	VisitService   int                `bson:"방문서비스"`
	ObjectSupport  int                `bson:"물품지원"`
	WelfareService int                `bson:"복지서비스"`
	LifeSupport    int                `bson:"생활지원"`
}
