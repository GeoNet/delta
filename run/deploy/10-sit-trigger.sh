#!/bin/bash

if [[ "$TRAVIS_BRANCH" == "main" ]] && [[ "$TRAVIS_PULL_REQUEST" == "false" ]]
then
    echo 'Triggering SIT Build'

    body='{
    "request": {
    "branch":"main"
    }}'

    STATUSCODE=$(curl --silent --output /dev/stderr --write-out "%{http_code}" -X POST \
       -H "Content-Type: application/json" \
       -H "Accept: application/json" \
       -H "Travis-API-Version: 3" \
       -H "Authorization: token $geonetci_api_token" \
       -d "$body" \
       https://api.travis-ci.com/repo/GeoNet%2Fsit/requests
    )

    if [[ $STATUSCODE -ne 200 ]] && [[ $STATUSCODE -ne 202 ]]
    then
            exit 1
    fi

fi
