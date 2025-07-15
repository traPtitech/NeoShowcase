package mocks

//go:generate go tool moq -stub -pkg $GOPACKAGE -out registry_mock.go ../../domain/builder RegistryClient
//go:generate go tool moq -stub -pkg $GOPACKAGE -out controller_mock.go ../../domain ControllerServiceClient
//go:generate go tool moq -stub -pkg $GOPACKAGE -out dbmanager_mock.go ../../domain MariaDBManager MongoDBManager
//go:generate go tool moq -stub -pkg $GOPACKAGE -out git_mock.go ../../domain GitService GitRepository
