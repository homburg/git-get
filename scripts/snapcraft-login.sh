#!/usr/bin/env bash

echo "${SNAPCRAFT_LOGIN}" | base64 -d > /tmp/snap.login
snapcraft login --with /tmp/snap.login
