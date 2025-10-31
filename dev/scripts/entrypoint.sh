#!/bin/sh
set -e  # Exit on error

cd /app

if [ "$ENVIRONMENT" = "development" ]; then
    echo "Starting in development mode..."

    echo "Starting dev server..."
    air &

    echo "Starting templ watch..."
    make templ-watch &

	echo "Starting tailwindcss watch..."
	make tailwind-watch &

    wait
else
    echo "Starting in production mode..."
    make build
    ./bin/app
fi
