package model

import (
	"log"
	"os"
	"fmt"
)


func init(){
	file, err := os.OpenFile("logs/devops.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	log.SetFlags(log.Ldate|log.Lmicroseconds|log.LUTC)
	log.SetOutput(file)
}
