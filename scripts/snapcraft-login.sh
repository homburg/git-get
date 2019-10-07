#!/usr/bin/env bash

echo "${SNAPCRAFT_LOGIN}" | base64 -d > snap.login
snapcraft login --with snap.login
