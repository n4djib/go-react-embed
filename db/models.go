package db

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Pokemon struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
