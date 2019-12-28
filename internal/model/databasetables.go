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


type TableRdsRsyncWorkorder struct {
	Id int	`json:"id" form:"id"`
	Database string	`json:"database" form:"database"`
	Username string	`json:"username" form:"username"`
	CreatedAt time.Time	`json:"created_at" form:"created_at"`
	PassAt time.Time	`json:"pass_at" form:"pass_at"`
	OrderStatus string	`json:"order_status" form:"order_status"`
}

type TableRdsRsyncOrder struct {
	Id int	`json:"id" form:"id"`
	Database string	`json:"database" form:"database"`
	Priority int	`json:"priority" form:"priority"`
	AuthorizedUser string	`json:"authorized_user" form:"authorized_user"`
}

type TableRdsRsyncOrderLogs struct {
	Id int	`json:"id" form:"id"`
	Workorderid int	`json:"workorderid" form:"workorderid"`
	Orderid int	`json:"orderid" form:"orderid"`
	Status	bool	`json:"status" form:"status"`
	CreatedAt	time.Time	`json:"created_at" form:"created_at"`
}
