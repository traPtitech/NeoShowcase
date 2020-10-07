//go:generate sqlboiler mysql
//go:generate pkger -o cmd/ns-migrate
//go:generate protoc -I ./api/proto --go_out=plugins=grpc:. --go_opt=module=github.com/traPtitech/neoshowcase ./api/proto/neoshowcase/entities/entity.proto
//go:generate protoc -I ./api/proto --go_out=plugins=grpc:. --go_opt=module=github.com/traPtitech/neoshowcase ./api/proto/neoshowcase/api/sites.proto
package main
