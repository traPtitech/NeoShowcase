package domain

// All environment variables injected by neoshowcase must have NS_ prefix
const (
	EnvPrefix = "NS_"

	EnvMySQLHostnameKey = EnvPrefix + "MYSQL_HOSTNAME"
	EnvMySQLPortKey     = EnvPrefix + "MYSQL_PORT"
	EnvMySQLUserKey     = EnvPrefix + "MYSQL_USER"
	EnvMySQLPasswordKey = EnvPrefix + "MYSQL_PASSWORD"
	EnvMySQLDatabaseKey = EnvPrefix + "MYSQL_DATABASE"

	EnvMongoDBHostnameKey = EnvPrefix + "MONGODB_HOSTNAME"
	EnvMongoDBPortKey     = EnvPrefix + "MONGODB_PORT"
	EnvMongoDBUserKey     = EnvPrefix + "MONGODB_USER"
	EnvMongoDBPasswordKey = EnvPrefix + "MONGODB_PASSWORD"
	EnvMongoDBDatabaseKey = EnvPrefix + "MONGODB_DATABASE"
)
