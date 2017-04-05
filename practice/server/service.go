package server

import (
	"context"
	"fmt"
)

// DataService handle all data processing
type DataService struct {
	*DB
}

// CallService will do data processing below surface
func (s DataService) CallService(ctx context.Context, query string) (Result, error) {
	res, err := QueryPorts(query)
	if !validate(res) {
		return res, err
	}
	err = s.CreateEntry(ctx, res)
	if err != nil {
		fmt.Printf("Database process failed with Err: %v", err)
		res.Err = err
		return res, err
	}
	return res, err
}

func validate(res Result) bool {
	if len(res.Ports) == 0 {
		return false
	}
	return true
}

// NewService for creating DataService
func NewService() DataService {
	db, err := InitialDB()
	if err != nil {
		fmt.Printf("Database Initial failed with Err: %v", err)
		panic(err)
	}
	return DataService{DB: db}
}
