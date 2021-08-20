module github.com/traPtitech/neoshowcase

go 1.16

require (
	github.com/aws/aws-sdk-go v1.38.30
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustinkirkland/golang-petname v0.0.0-20191129215211-8e5a1ed0cff0
	github.com/friendsofgo/errors v0.9.2
	github.com/fsouza/go-dockerclient v1.7.2
	github.com/gavv/httpexpect/v2 v2.2.0
	github.com/go-git/go-git/v5 v5.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/mock v1.5.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/wire v0.5.0
	github.com/kat-co/vala v0.0.0-20170210184112-42e1d8b61f12
	github.com/ktr0731/evans v0.9.3
	github.com/labstack/echo/v4 v4.2.2
	github.com/leandro-lugaresi/hub v1.1.1
	github.com/moby/buildkit v0.8.3
	github.com/ncw/swift v1.0.53
	github.com/rubenv/sql-migrate v0.0.0-20210408115534-a32ed26c37ea
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/randomize v0.0.1
	github.com/volatiletech/sqlboiler/v4 v4.5.0
	github.com/volatiletech/strmangle v0.0.1
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.mongodb.org/mongo-driver v1.5.1
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	google.golang.org/grpc v1.35.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.1-0.20201117152358-0edc412565dc
	github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
)
