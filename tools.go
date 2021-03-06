//+build tools

package main

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/ktr0731/evans"
	_ "github.com/rubenv/sql-migrate"
	_ "github.com/volatiletech/sqlboiler/v4"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
