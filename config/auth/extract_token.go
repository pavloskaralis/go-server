package auth

import (
	"net/http"
	"strings"
)

//extract token from header
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
	   return strArr[1]
	}
	return ""
}