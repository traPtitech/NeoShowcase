package domain

// All environment variables injected by neoshowcase must have NS_ prefix
const (
	EnvPrefix             = "NS_"
	EnvMySQLUserKey       = EnvPrefix + "MYSQL_USER"
	EnvMySQLPasswordKey   = EnvPrefix + "MYSQL_PASSWORD"
	EnvMySQLDatabaseKey   = EnvPrefix + "MYSQL_DATABASE"
	EnvMySQLHostnameKey   = EnvPrefix + "MYSQL_HOSTNAME"
	EnvMongoDBUserKey     = EnvPrefix + "MONGODB_USER"
	EnvMongoDBPasswordKey = EnvPrefix + "MONGODB_PASSWORD"
	EnvMongoDBDatabaseKey = EnvPrefix + "MONGODB_DATABASE"
	EnvMongoDBHostnameKey = EnvPrefix + "MONGODB_HOSTNAME"
)
