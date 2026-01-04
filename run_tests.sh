#!/bin/bash

echo "🧪 Running Store Layer Tests..."
echo ""

# Check if test database exists
echo "📊 Checking test database..."
if mysql -u user -ppassword -e "USE mlm_test;" 2>/dev/null; then
    echo "✅ Test database exists"
else
    echo "❌ Test database 'mlm_test' not found"
    echo "💡 Run: ./setup_test_db.sh"
    exit 1
fi

echo ""
echo "🔬 Running tests..."
echo ""

# Run tests with verbose output
go test -v -count=1 ./internal/musicapp/lib/users/store/...

# Check exit code
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All tests passed!"
    echo ""
    echo "📊 Run with coverage:"
    echo "   go test -v -cover ./internal/musicapp/lib/users/store/..."
    echo ""
    echo "🏃 Run with race detection:"
    echo "   go test -v -race ./internal/musicapp/lib/users/store/..."
else
    echo ""
    echo "❌ Some tests failed"
    echo "💡 Check the output above for details"
    exit 1
fi
