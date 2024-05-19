package db

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" binding:"required"`
}

type Pokemon struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Data struct {
	Data interface{} `json:"data"`
}

type DataList struct {
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}
