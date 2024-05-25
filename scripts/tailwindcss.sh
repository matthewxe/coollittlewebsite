#!/bin/sh
npx tailwindcss -i ../build/main.css -c ../build/tailwind.config.js -o ../web/static/main.css --minify $*
