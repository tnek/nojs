#!/bin/bash
docker build . -t notes-site
docker run -d -p ${HOST_PORT}:8080 notes-site
