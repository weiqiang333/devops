package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Instance struct {
	Id	int
	Cloud	string
	Address	string
	PublicAddress	string
	Name	string
	Cpu	int
	Mem	int
	InstanceType	string
	OperatingSystem	string
	DiskCapacity	JsonToMap
	State	string
}

type JsonToMap map[string]interface{}

func (m JsonToMap) Value() (driver.Value, error) {
	j, err := json.Marshal(m)
	return j, err
}

func (m *JsonToMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*m, ok = i.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}