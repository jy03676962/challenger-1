#!/bin/bash
webpack

cp -r 'dist' '../server'

rm -rf '../server/public'

mv '../server/dist' '../server/public'

cd ../server


go run main.go
