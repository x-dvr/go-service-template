package user

import (
	"encoding/json"
	"strings"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		Email string `json:"email"`
	}{
		Alias: (*Alias)(&u),
		Email: strings.Repeat("*", 4) + "@mail.com",
	})
}
