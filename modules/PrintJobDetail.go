package modules

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type job_detail struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Duration       string             `bson:"구인신청일자" json:"구인신청일자"`
	CompanyName    string             `bson:"사업장명" json:"사업장명"`
	RecuritType    string             `bson:"모집직종" json:"모집직종"`
	RecuritShape   string             `bson:"고용형태" json:"고용형태"`
	WageType       string             `bson:"임금형태" json:"임금형태"`
	Wage           int                `bson:"임금" json:"임금"`
	ComeType       string             `bson:"입사형태" json:"입사형태"`
	RequireHistory string             `bson:"요구경력" json:"요구경력"`
	RequireStudy   string             `bson:"요구학력" json:"요구학력"`
	RelateMajor    string             `bson:"전공계열" json:"전공계열"`
	RequireLicense string             `bson:"요구자격증" json:"요구자격증"`
	Address        string             `bson:"사업장 주소" json:"사업장 주소"`
	CompanyType    string             `bson:"기업형태" json:"기업형태"`
	ResponsiveIns  string             `bson:"담당기관" json:"담당기관"`
	CreateAt       string             `bson:"등록일" json:"등록일"`
	Contact        string             `bson:"연락처" json:"연락처"`
	BodySpec       string             `bson:"필수부위" json:"필수부위"`
}

func PrintJobDetail(w http.ResponseWriter, r *http.Request, urlPath *[]string) {
	oid_hex := (*urlPath)[1]
	oid, err := primitive.ObjectIDFromHex(oid_hex)
	ErrOK(err)
	dj_temp, err := OidTOjobDetail(oid)
	temp := job_detail{
		ID:             dj_temp.ID,
		Duration:       dj_temp.Duration,
		CompanyName:    dj_temp.CompanyName,
		RecuritType:    dj_temp.RecuritType,
		RecuritShape:   dj_temp.RecuritShape,
		WageType:       dj_temp.WageType,
		Wage:           dj_temp.Wage,
		ComeType:       dj_temp.ComeType,
		RequireHistory: dj_temp.RequireHistory,
		RequireStudy:   dj_temp.RequireStudy,
		RelateMajor:    dj_temp.RelateMajor,
		RequireLicense: dj_temp.RequireLicense,
		Address:        dj_temp.Address,
		CompanyType:    dj_temp.CompanyType,
		ResponsiveIns:  dj_temp.ResponsiveIns,
		CreateAt:       dj_temp.CreateAt,
		Contact:        dj_temp.Contact,
		BodySpec:       dj_temp.BodySpec,
	}

	if err != nil { //job을 찾지 못하였을 때
		temp := map[string]string{"사업장명": "찾지 못함"}
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
