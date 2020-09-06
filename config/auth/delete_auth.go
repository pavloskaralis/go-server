
package auth


func DeleteAuth(givenUuid string) (int64,error) {
	deleted, err := Redis.Del(givenUuid).Result()
	if err != nil {
	   return 0, err
	}
	return deleted, nil
}
