#!/bin/bash

rsync -ru --delete ../* root@staging.httprelay.io:~/httprelay-stag
ssh root@staging.httprelay.io "
cd ~/httprelay-stag/deploy
COMPOSE_PROJECT_NAME=httprelay-stag SUB_DOMAIN=staging DOMAIN=httprelay.io docker-compose down &
COMPOSE_PROJECT_NAME=httprelay-stag SUB_DOMAIN=staging DOMAIN=httprelay.io docker-compose up --build
"