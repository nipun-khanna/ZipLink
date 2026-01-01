package utils

func IDToShortURL(id int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURL := ""

	for id > 0 {
		shortURL = string(charset[id%62]) + shortURL
		id = id / 62
	}
	return shortURL
}

func ShortURLToID(url string) int {
	id := 0
	for i := 0; i < len(url); i++ {
		c := int(url[i])
		if c >= int('a') && c <= int('z') {
			id = id*62 + c - int('a')
		} else if c >= int('A') && c <= int('Z') {
			id = id*62 + c - int('A') + 26
		} else {
			id = id*62 + c - int('0') + 52
		}
	}
	return id
}
