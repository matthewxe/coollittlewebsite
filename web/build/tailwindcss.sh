#!/usr/bin/env sh
# Run in the web/ folder
# pwd
# npx tailwindcss -i ../build/main.css -c ../build/tailwind.config.js -o ../web/static/whataboutme/main.css --minify $*
tailwindcss -i build/main.css -c build/tailwind.config.js -o static/whataboutme/main.css --minify $*
