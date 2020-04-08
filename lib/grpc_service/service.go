package grpc_service

import "time"

type Result struct {
	ID        string
	ServiceName                 string
	Payload              []byte
	Err                  string
	Elapse time.Duration
}