package modules

import (
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type session_strcut_tmp struct {
	ID        primitive.ObjectID      `bson:"_id,omitempty"`
	Email     string                  `bson:"email"`
	LastLogin time.Time               `bson:"lastLogin"`
	ScrapList []primitive.ObjectID    `bson:"scrapList"`
	Settings  Dj_users_users_settings `bson:"settings"`
}

func PrintSession(w http.ResponseWriter, r *http.Request) {
	if IsHeLogin(w, r) {
		oid := SessionTO_oid(w, r)
		temp := OidTOuser_struct(oid)
		temp2 := session_strcut_tmp{
			ID:        temp.ID,
			Email:     temp.Email,
			LastLogin: temp.LastLogin,
			ScrapList: temp.ScrapList,
			Settings:  temp.Settings,
		}
		temp3, _ := json.MarshalIndent(temp2, " ", "	")

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(temp3)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		msg := map[string]string{"error": "Not LOGIN"}
		msg_json, _ := json.Marshal(msg)
		w.Write(msg_json)
	}

}
