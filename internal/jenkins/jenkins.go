package jenkins

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bndr/gojenkins"
)


//ConnJenkins 连接 jenkins
func ConnJenkins(jenkinsBase, username, token string) (*gojenkins.Jenkins, error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	jenkinsClient := gojenkins.CreateJenkins(client, jenkinsBase, username, token)
	_, err := jenkinsClient.Init()

	if err != nil {
		return jenkinsClient, fmt.Errorf("ConnJenkins Something Went Wrong %v", err)
	}
	return jenkinsClient, nil
}
