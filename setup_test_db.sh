#!/bin/bash

echo "Setting up test database for IMAPP pattern..."

# Database credentials
DB_USER="user"
DB_PASS="password"
DB_HOST="127.0.0.1"
DB_PORT="3306"
TEST_DB="mlm_test"

# Create test database
echo "Creating test database: $TEST_DB"
mysql -u $DB_USER -p$DB_PASS -h $DB_HOST -P $DB_PORT <<EOF
CREATE DATABASE IF NOT EXISTS $TEST_DB;
USE $TEST_DB;
EOF

# Run migrations
echo "Running migrations..."
mysql -u $DB_USER -p$DB_PASS -h $DB_HOST -P $DB_PORT $TEST_DB < migration/01_create_users.up.sql
mysql -u $DB_USER -p$DB_PASS -h $DB_HOST -P $DB_PORT $TEST_DB < migration/02_create_rooms.up.sql
mysql -u $DB_USER -p$DB_PASS -h $DB_HOST -P $DB_PORT $TEST_DB < migration/03_create_room_members.up.sql
mysql -u $DB_USER -p$DB_PASS -h $DB_HOST -P $DB_PORT $TEST_DB < migration/04_update_users_for_imapp.up.sql

echo "✅ Test database setup complete!"
echo ""
echo "To run tests:"
echo "  go test -v -count=1 ./internal/musicapp/lib/users/store/..."
