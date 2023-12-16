package client

//
// import (
//	"errors"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
// )
//
// var (
//	ErrCacheMiss = errors.New("cache miss")
// )
//
// func asClientError(err error) error {
//	gstatus, ok := status.FromError(err)
//	if !ok {
//		return err
//	}
//
//	switch gstatus.Code() {
//	case codes.NotFound:
//		return ErrCacheMiss
//	default:
//		return err
//	}
// }
