package model

type Admin struct {
	UserName 		string 	`json:"user_name"`
	Password 		string 	`json:"password"`
	AuthorityLevel 	int 	`json:"authority_level"`
}