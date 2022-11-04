#!/bin/bash

docker build . -f capturer.Dockerfile -t capturer
docker tag capturer localhost:5001/capturer
docker push localhost:5001/capturer

docker build . -f demo-app/s3.Dockerfile -t s3test
docker tag s3test localhost:5001/s3test
docker push localhost:5001/s3test