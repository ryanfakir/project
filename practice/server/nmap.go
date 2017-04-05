package server

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// QueryPorts for hostname or ip lookup
func QueryPorts(query string) (Result, error) {
	var err error
	var ids []int
	cmdName := "nmap"
	cmdArgs := []string{"-p", "1-1000", query}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			lines := strings.Split(line, "/")
			i, err := strconv.Atoi(lines[0])
			if err != nil {
				continue
			}
			ids = append(ids, i)
		}
	}()
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
	start := time.Now()
	return Result{Hostname: query, Ports: ids, QueryTime: start}, nil
}
