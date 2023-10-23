package utils

func IsPersonCard(id string) bool {

	b := []byte(id)
	if len(b) == 18  {
		for _,bl := range b {
			if bl < '0' || bl > '9' {
				return false
			}
		}
		return true
	}
	return false
}
