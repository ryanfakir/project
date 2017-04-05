package server

import "fmt"

var config = &DBConfig{Username: "root", Password: "root", Host: "127.0.0.1", Port: 3306}

//DBConfig with usr and pw
type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

// DbURL return url for connection
func (c *DBConfig) DbURL() string {
	cred := c.Username + ":" + c.Password + "@"
	//username:password@tcp(host:port)/
	return fmt.Sprintf("%stcp(%s:%d)/", cred, c.Host, c.Port)
}
