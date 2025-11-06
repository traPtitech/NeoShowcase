//go:build tools

package main

import (
	_ "github.com/aarondl/sqlboiler/v4"
	_ "github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-mysql"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/google/wire/cmd/wire"
)
