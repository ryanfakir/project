package main

import (
	"os"
	"os/signal"
	"practice/server"
	"syscall"
)

func main() {
	s := server.New()
	go func() {
		s.Open()
	}()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	s.Close()
}
