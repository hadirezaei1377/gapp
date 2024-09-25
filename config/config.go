package config

import (
	"gapp/repository/mysql"
	"gapp/service/authservice"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer         `koanf:"http_server"`
	Auth       authservice.Config `koanf:"auth"`
	Mysql      mysql.Config       `koanf:"mysql"`
}
