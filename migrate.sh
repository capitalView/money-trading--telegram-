#!/bin/bash

# get .env in file
export $(grep -v '^#' .env | xargs)

# check if the number of arguments is correct
if [ "$#" -lt 1 ] || [ "$#" -gt 2 ]; then
  echo "Usage: $0 <direction> [steps]"
  echo "Example: $0 down 1"
  exit 1
fi

MIGRATE_PATH="./db/migrations"

# check if the migration direction is correct
MIGRATION_DIRECTION="$1"


if [ "$MIGRATION_DIRECTION" = "create" ]; then
  echo "Create migration..."
  migrate create -ext sql -dir db/migrations -seq "$2"    # Создаем новую миграцию
else
  echo "Running migration up..."
  migrate -path "$MIGRATE_PATH" -database "postgres://$DB_USER:$DB_PASSWORD@localhost:$DB_PORT/$DB_NAME?sslmode=disable" "$MIGRATION_DIRECTION"
fi


# check status of the migration
if [ $? -eq 0 ]; then
  echo "Migration $MIGRATION_DIRECTION completed successfully."
else
  echo "Migration $MIGRATION_DIRECTION failed."
fi