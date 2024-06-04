package api

import "github.com/go-playground/validator/v10"

// Initialize validator
var validate = validator.New()

// type Data struct {
// 	Message string `json:"message"`
// 	Data interface{} `json:"data"`
// }

// type ResultList struct {
// 	Count  int         `json:"count"`
// 	Limit  int         `json:"limit"`
// 	Offset int         `json:"offset"`
// 	Result interface{} `json:"result"`
// }
