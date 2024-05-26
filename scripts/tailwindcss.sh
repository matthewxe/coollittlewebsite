#!/usr/bin/env sh
# Run in the scripts folder
pwd
npx tailwindcss -i ../build/main.css -c ../build/tailwind.config.js -o ../web/static/whataboutme/main.css --minify $*
