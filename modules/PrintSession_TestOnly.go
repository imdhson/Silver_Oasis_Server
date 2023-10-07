package modules

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PrintSession_TestOnly(w http.ResponseWriter, r *http.Request) {
	oid, _ := primitive.ObjectIDFromHex("64982a3c310879ece0939b32")
	temp := OidTOuser_struct(oid)
	temp2, _ := json.MarshalIndent(temp, " ", "	")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(temp2)
}
