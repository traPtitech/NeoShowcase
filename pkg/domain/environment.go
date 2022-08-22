package domain

// All environment variables injected by neoshowcase must have NS_ prefix
const (
	EnvPrefix             = "NS_"
	EnvMySQLPasswordKey   = EnvPrefix + "MYSQL_PASSWORD"
	EnvMySQLDatabaseKey   = EnvPrefix + "MYSQL_DATABASE"
	EnvMongoDBPasswordKey = EnvPrefix + "MONGODB_PASSWORD"
	EnvMongoDBDatabaseKey = EnvPrefix + "MYSQL_DATABASE"
)
