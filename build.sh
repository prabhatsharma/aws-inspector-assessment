#!/bin/bash

echo '--- starting build ---'
go build main.go
echo '--- build complete ---'

echo '--- building archive ---'
zip main.zip main
echo '--- archiving build complete ---'

echo '--- uploading archive to s3 ---'
aws s3 cp main.zip s3://prabhat00-public/
echo '--- uploading archive to s3 complete ---'