#!/bin/bash

set -a
if [ -f ../.env ]; then
  export $(grep -v '^#' ../.env | xargs -0)
  echo ".env file loaded"
else
  echo ".env file not found. Ensure it exists in the correct directory."
  exit 1
fi
set +a

migrate -database "$DATABASE_URL" -path ../migrations down -all
if [ $? -ne 0 ]; then
  echo "Migration down failed!"
  exit 1
fi

migrate -database "$DATABASE_URL" -path ../migrations up
if [ $? -ne 0 ]; then
  echo "Migration up failed!"
  exit 1
fi

go run ./seed.go
if [ $? -ne 0 ]; then
  echo "Seeder failed!"
  exit 1
fi

echo "Database migration down, migration up, and seeding completed successfully!"