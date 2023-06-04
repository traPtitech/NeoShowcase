package domain

// All environment variables injected by neoshowcase must have NS_ prefix
const (
	EnvPrefix = "NS_"

	EnvMariaDBHostnameKey = EnvPrefix + "MARIADB_HOSTNAME"
	EnvMariaDBPortKey     = EnvPrefix + "MARIADB_PORT"
	EnvMariaDBUserKey     = EnvPrefix + "MARIADB_USER"
	EnvMariaDBPasswordKey = EnvPrefix + "MARIADB_PASSWORD"
	EnvMariaDBDatabaseKey = EnvPrefix + "MARIADB_DATABASE"

	EnvMongoDBHostnameKey = EnvPrefix + "MONGODB_HOSTNAME"
	EnvMongoDBPortKey     = EnvPrefix + "MONGODB_PORT"
	EnvMongoDBUserKey     = EnvPrefix + "MONGODB_USER"
	EnvMongoDBPasswordKey = EnvPrefix + "MONGODB_PASSWORD"
	EnvMongoDBDatabaseKey = EnvPrefix + "MONGODB_DATABASE"
)
