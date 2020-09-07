package auth

import (
	"time"
)

func CreateAuth(userid string, td *TokenDetails) error {
	//unix to UTC
    at := time.Unix(td.AtExpires, 0) 
    rt := time.Unix(td.RtExpires, 0)
    now := time.Now()
	//store tokens
    errAccess := Redis.Set(td.AccessUuid, userid, at.Sub(now)).Err()
    if errAccess != nil {
        return errAccess
    }
    errRefresh := Redis.Set(td.RefreshUuid, userid, rt.Sub(now)).Err()
    if errRefresh != nil {
        return errRefresh
    }
    return nil
}