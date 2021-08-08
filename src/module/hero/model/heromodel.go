package heromodel

type Hero struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ReturnFirstThreeCharacters will return the first charactesrs of the input string
func ReturnFirstThreeCharacters(name string) string {
	return name[0:3]
}
