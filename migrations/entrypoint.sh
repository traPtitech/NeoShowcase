#!/usr/bin/env sh

set -eu

SQL_FILE=$1

APP_VERSION=${APP_VERSION:-unknown}
APP_REVISION=${APP_REVISION:-unknown}

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-root}
DB_PASS=${DB_PASS:-password}
DB_NAME=${DB_NAME:-neoshowcase}

echo "ns-migrate $APP_VERSION@$APP_REVISION"

mysqldef --host="$DB_HOST" --port="$DB_PORT" --user="$DB_USER" --password="$DB_PASS" "$DB_NAME" < "$SQL_FILE"
