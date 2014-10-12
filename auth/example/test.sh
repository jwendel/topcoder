#!/bin/bash

#password generated using:
#   echo -n ilovetopcoder | openssl dgst -binary -sha256 | openssl base64

echo "### topcoder.com pass"
curl -i --data "username=takumi&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost:8080/api/2/domains/topcoder.com/proxyauth
echo 

echo;echo;echo "### appirio.com pass"
curl -i --data "username=jun&password={SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=" http://localhost:8080/api/2/domains/appirio.com/proxyauth
echo

echo;echo;echo "### appirio.com password fail"
curl -i --data "username=jun&password={SHA256}p/GZlDycaRbeggsp0wQuvZQ7yV4IzVstwEKKFhGNyGo=" http://localhost:8080/api/2/domains/appirio.com/proxyauth
echo

echo;echo;echo "### appirio.com username fail"
curl -i --data "username=kyrra&password={SHA256}p/GZlDycaRbeggsp0wQuvZQ7yV4IzVstwEKKFhGNyGo=" http://localhost:8080/api/2/domains/appirio.com/proxyauth
echo

echo;echo;echo "### domain fail, expect 404"
curl -i --data "username=takumi&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost:8080/api/2/domains/example.com/proxyauth
echo
