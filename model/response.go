package model

import "time"

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

type HttpResponse struct {
	Data       interface{} `json:"data"`
	MetaData   MetaData    `json:"metadata"`
	Pagination *Pagination `json:"pagination"`
}

type MetaData struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type AuthResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type UserResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
