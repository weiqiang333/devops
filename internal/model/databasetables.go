package model

import "time"


type TableService struct {
	Id int	`json:"id" form:"id"`
	Server string	`json:"server" form:"server"`
	Service string	`json:"service" form:"service"`
	Status string	`json:"status" form:"status"`
}


type TableGoogleAuth struct {
	Id int	`json:"id" form:"id"`
	Name string	`json:"name" form:"name"`
	Secret string	`json:"secret" form:"secret"`
	UpdatedAt time.Time `json:"updated_at"`
}


type TableLdapPwdExpired struct {
	Id int	`json:"id" form:"id"`
	Name string	`json:"name" form:"name"`
	PwdLastSet time.Time	`json:"pwd_last_set"`
	PwdExpired time.Time	`json:"pwd_expired"`
}
