#!/bin/bash

if [[ "$TRAVIS_BRANCH" == "master" ]] && [[ "$TRAVIS_PULL_REQUEST" == "false" ]]
then
    echo 'Triggering SIT Build'

    body='{
    "request": {
    "branch":"master"
    }}'

    curl -s -X POST \
       -H "Content-Type: application/json" \
       -H "Accept: application/json" \
       -H "Travis-API-Version: 3" \
       -H "Authorization: token $TRAVIS_API_TOKEN" \
       -d "$body" \
       https://api.travis-ci.com/repo/GeoNet%2Fsit/requests
fi