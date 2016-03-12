#!/bin/bash

INDEX_DIST='./public'
JS_DIST='./public/js'

git submodule foreach git pull origin master
git commit -a -m 'update submodules'
git push origin master

cd 'web'

npm install

npm run build

cd ..

test -d $INDEX_DIST || mkdir -p $INDEX_DIST && cp './web/dist/index.html' $INDEX_DIST
test -d $JS_DIST || mkdir -p $JS_DIST && find './web/dist' -name '*.js' -exec cp {} $JS_DIST \;
