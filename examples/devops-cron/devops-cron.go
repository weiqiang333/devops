package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/weiqiang333/devops/internal/crontab"
)


func main()  {
	crontab.CronTab()

	container := restful.NewContainer()
	api := new(restful.WebService)
	api.Path("/crontab/v1")
	api.Produces(restful.MIME_JSON)
	api.Route(api.GET("/hello").To(hello))
	container.Add(api)
	server := &http.Server{Addr: ":8081", Handler: container}
	log.Fatal(server.ListenAndServe())
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}

func init()  {
	cfgFile := flag.String("config", "", "config file (default is $HOME/.devops.yaml)")
	flag.Parse()

	initConfig(*cfgFile)
}

func initConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".devops" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".devops")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}