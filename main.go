package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/prabhatsharma/aws-inspector-assessment/helper"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Detail is event detail structure
type Detail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}

// HandleRequest lambda handler
func HandleRequest(ctx context.Context, cEvent events.CloudWatchEvent) (string, error) {
	log.Println("-----------Execution begins------------")

	var detail Detail
	json.Unmarshal(cEvent.Detail, &detail)
	log.Println("InstanceID: ", detail.InstanceID)

	helper.SetTag(&detail.InstanceID, "true")
	log.Println("SetTag:true complete")

	time.Sleep(60 * time.Second) // sleep for 60 seconds allowing instance to start
	log.Println("60 seconds sleep to allow ec2 to initialize complete")

	helper.Begin(detail.InstanceID)
	time.Sleep(60 * time.Second) // sleep for 60 seconds allowing scanning to begin
	helper.SetTag(&detail.InstanceID, "false")

	log.Println("-----------Execution ends------------")
	return "Execution Completed", nil
}

func main() {
	lambda.Start(HandleRequest)

}
