package entities

type User struct{}

type Project struct {
	Users []User
}

type Organization struct {
	Users []User
}

type CorpGroup struct {
	Users []User
}

type Community struct {
	Users []User
}
