#!/bin/bash

aws cloudformation delete-stack --stack-name=inspectec2 
sleep 20
aws cloudformation create-stack --capabilities CAPABILITY_NAMED_IAM --stack-name=inspectec2 --parameters '[{"ParameterKey": "TOPICARN", "ParameterValue": "arn:aws:sns:us-west-2:107995894928:mailer"}]' --template-body file://setup-cfn.yaml



