package grpc_service

import "time"

type Result struct {
	Elapse time.Duration
	Data interface{}
}