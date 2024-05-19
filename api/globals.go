package api

import "github.com/go-playground/validator/v10"

// Initialize validator
var validate = validator.New()

type Data struct {
	Data interface{} `json:"data"`
}

type DataList struct {
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}

type Error struct {
	Error string `json:"error"`
}

type Status struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

