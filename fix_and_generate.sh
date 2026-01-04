#!/bin/bash

echo "🔧 Fixing database setup and generating models..."

# Step 1: Start Docker MySQL
echo ""
echo "📦 Starting MySQL (Docker)..."
docker-compose up -d
sleep 3

# Step 2: Run migrations
echo ""
echo "🔄 Running migrations..."
mysql -u user -puserpass -h 127.0.0.1 musicapp < migration/01_create_users.up.sql 2>/dev/null || echo "  (01 may already exist)"
mysql -u user -puserpass -h 127.0.0.1 musicapp < migration/02_create_rooms.up.sql 2>/dev/null || echo "  (02 may already exist)"
mysql -u user -puserpass -h 127.0.0.1 musicapp < migration/03_create_room_members.up.sql 2>/dev/null || echo "  (03 may already exist)"

# Step 3: Generate SQLBoiler models
echo ""
echo "🔧 Generating SQLBoiler models..."
sqlboiler mysql --wipe

# Step 4: Verify build
echo ""
echo "✅ Testing build..."
if go build ./...; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed - check errors above"
    exit 1
fi

echo ""
echo "🎉 All fixed! Models generated in models/"
echo ""
echo "Next steps:"
echo "  1. Set up test database: ./setup_test_db.sh"
echo "  2. Run tests: ./run_tests.sh"
