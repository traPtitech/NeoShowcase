module github.com/traPtitech/neoshowcase

go 1.16

require (
	github.com/aws/aws-sdk-go v1.37.33
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustinkirkland/golang-petname v0.0.0-20191129215211-8e5a1ed0cff0
	github.com/friendsofgo/errors v0.9.2
	github.com/fsouza/go-dockerclient v1.7.1
	github.com/gavv/httpexpect/v2 v2.2.0
	github.com/go-git/go-git/v5 v5.2.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kat-co/vala v0.0.0-20170210184112-42e1d8b61f12
	github.com/labstack/echo/v4 v4.2.0
	github.com/leandro-lugaresi/hub v1.1.1
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/moby/buildkit v0.8.1
	github.com/ncw/swift v1.0.53
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/rubenv/sql-migrate v0.0.0-20210215143335-f84234893558
	github.com/sirupsen/logrus v1.8.0
	github.com/spf13/afero v1.4.0 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/randomize v0.0.1
	github.com/volatiletech/sqlboiler/v4 v4.4.0
	github.com/volatiletech/strmangle v0.0.1
	go.mongodb.org/mongo-driver v1.4.6
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	google.golang.org/genproto v0.0.0-20201008135153-289734e2e40c // indirect
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/ini.v1 v1.61.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.1-0.20201215193253-e922d5553d12
	github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
)
