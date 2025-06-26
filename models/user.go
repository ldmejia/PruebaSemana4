package models

type ApiResponse struct {
	Results []User `json:"results"`
}

type User struct {
	Gender string `json:"gender"`
	Name Name `json:"name"`
	Location Location `json:"location"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Name struct {
	First string `json:"first"`
	Last string `json:"last"`
}

type Location struct{
	Street struct {
		Number int `json:"number"`
		Name string `json:"name"`
	}
	City string `json:"city"`
	State string `json:"state"`
	Country string `json:"country"`
}

type Job struct{
	ID int 
}