package cmd

import (
	"go-blog/config"
	"go-blog/server"
)

func cmdStart() error {
	var err error
	defer func(){
		log.Error("Error when start server: ", err)
	}()

	cfg, err := config.OpenFileConfig()
	if err != nil {
		return err
	}
	ser := server.NewGinServer(cfg)
	go ser.Start()
	return ser.Run()
}
