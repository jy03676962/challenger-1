#!/bin/bash
cd web

webpack

cp -r 'dist' '../server'

rm -rf '../server/public'

mv '../server/dist' '../server/public'
