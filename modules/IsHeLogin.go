package modules

import (
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsHeLogin(w http.ResponseWriter, r *http.Request) bool {
	oid := SessionTO_oid(w, r)
	if oid != primitive.NilObjectID {
		users_struct := OidTOuser_struct(oid)
		if strings.Contains(users_struct.Email, "@") {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
