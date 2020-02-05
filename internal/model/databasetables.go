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

type TableRdsRsyncWorkorderLogs struct {
	Id int	`json:"id" form:"id"`
	Workorderid	int	`json:"workorderid" form:"workorderid"`
	Username	string	`json:"username" form:"username"`
	CreatedAt	time.Time	`json:"created_at" form:"created_at"`
	GetSnapshotAt	time.Time	`json:"get_snapshot_at" form:"get_snapshot_at"`
	DeleteAt	time.Time	`json:"delete_at" form:"delete_at"`
	RestoreAt	time.Time	`json:"restore_at" form:"restore_at"`
	ModifyConfigAt	time.Time	`json:"modify_config_at" form:"modify_config_at"`
	ExecuteSqlAt	time.Time	`json:"execute_sql_at" form:"execute_sql_at"`
	Status	string	`json:"status" form:"status"`
}

//TableAwsInstanceTypes tables aws_instance_types
type TableAwsInstanceTypes struct {
	Id	int	`json:"id"`
	InstanceType	string	`json:"instance_type"`
	Vcpus	int64	`json:"vcpus"`
	Memory	int64	`json:"memory"`
	UpdateAt	time.Time	`json:"update_at"`
}

//TableAwsVolumes tables aws_volumes
type TableAwsVolumes struct {
	Id	int	`json:"id"`
	DiskId	string	`json:"disk_id"`
	CreateAt	time.Time	`json:"create_at"`
	Size	int64	`json:"size"`	
	Iops	int64	`json:"iops"`
	State	string	`json:"state"`
	Type	string	`json:"type"`
	TagName	string	`json:"tag_name"`
	TagApp	string	`json:"tag_app"`
	TagEnv	string	`json:"tag_env"`
	TagPillar	string	`json:"tag_pillar"`
	UpdateAt	time.Time	`json:"update_at"`
}
