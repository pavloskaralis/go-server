package auth

func FetchAuth(authD *AccessDetails) (string, error) {
	userID, err := Redis.Get(authD.AccessUuid).Result()
	if err != nil {
	   return "", err
	}
	return userID, nil
  }