#!/usr/bin/bash

echo "[+] Creating Directory staticdeps"
mkdir -p staticdeps
cd staticdeps

echo "[+] Downloading all JavaScript Dependencies."
#curl http://code.jquery.com/jquery-2.1.4.min.js --output jquery-2.1.4.min.js
curl https://unpkg.com/react@16.0.0/umd/react.production.min.js --output react.production.min.js
curl https://unpkg.com/react-dom@16.0.0/umd/react-dom.production.min.js --output react-dom.production.min.js
curl https://unpkg.com/babel-standalone@6.26.0/babel.js --output babel.js
curl http://code.jquery.com/jquery-2.1.4.min.js --output jquery-2.1.4.min.js


echo "[+] Downloading all CSS Dependencies."
curl https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css --output bootstrap.min.css
 