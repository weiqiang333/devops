package model


type TableService struct {
	Id int	`json:"id" form:"id"`
	Server string	`json:"server" form:"server"`
	Service string	`json:"service" form:"service"`
	Status string	`json:"status" form:"status"`
}