package helper

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/inspector"
)

func createResourceGroup(svc inspector.Inspector, randomTag string) string {
	// 1. create resource group input
	log.Println("1. CreateResourceGroupInput started")
	rgi := &inspector.CreateResourceGroupInput{
		ResourceGroupTags: []*inspector.ResourceGroupTag{
			{
				Key:   aws.String("inspector"),
				Value: aws.String(randomTag),
			},
		},
	}

	log.Println("1.1 CreateResourceGroup started")

	// 1.1. create resource group
	rg, rgerr := svc.CreateResourceGroup(rgi)

	if rgerr != nil {
		print(rgerr)
	}
	// print("Resource group: ", *rg.ResourceGroupArn)

	return *rg.ResourceGroupArn
}

func createAssessmentTarget(rgArn string, svc inspector.Inspector, InstanceID string) string {
	// 2. Create assessment target
	log.Println("2. Assessment target creation started")
	ati := &inspector.CreateAssessmentTargetInput{
		AssessmentTargetName: aws.String(InstanceID + "_AssessmentTarget_" + time.Now().Format("2006-01-02_15.04.05")),
		ResourceGroupArn:     &rgArn,
	}

	at, aterr := svc.CreateAssessmentTarget(ati)

	if aterr != nil {
		fmt.Println(aterr)
	}

	return *at.AssessmentTargetArn
}

func getRulesPackages(svc inspector.Inspector) *inspector.ListRulesPackagesOutput {
	// 3. create rules package input
	rpi := &inspector.ListRulesPackagesInput{
		MaxResults: aws.Int64(123),
	}

	rp, erp := svc.ListRulesPackages(rpi)
	if erp != nil {
		fmt.Println(erp.Error())
	}

	return rp
}

func createAssessmentTemplate(at string, rp *inspector.ListRulesPackagesOutput, svc inspector.Inspector, InstanceID string) *inspector.CreateAssessmentTemplateOutput {
	// 4. create assessment template
	atli := &inspector.CreateAssessmentTemplateInput{
		AssessmentTargetArn:    aws.String(at),
		AssessmentTemplateName: aws.String(InstanceID + "_AssessmentTemplate_" + time.Now().Format("2006-01-02_15.04.05")),
		DurationInSeconds:      aws.Int64(3600),
		RulesPackageArns:       rp.RulesPackageArns,
		UserAttributesForFindings: []*inspector.Attribute{
			{
				Key:   aws.String("inspection-type"),
				Value: aws.String("on-launch"),
			},
		},
	}

	atl, atlerr := svc.CreateAssessmentTemplate(atli)

	if atlerr != nil {
		fmt.Println(atlerr)
	}

	return atl
}

func subscribeToEvent(svc inspector.Inspector, resourceArn string, topicArn string) (subscribeResponse string) {
	// 5. Subscribe to event - ASSESSMENT_RUN_STARTED
	steInput := &inspector.SubscribeToEventInput{
		Event:       aws.String("ASSESSMENT_RUN_STARTED"),
		ResourceArn: aws.String(resourceArn),
		TopicArn:    aws.String(topicArn),
	}

	fmt.Println("about to execute SubscribeToEvent() with:")
	fmt.Println(steInput.String())

	result, err := svc.SubscribeToEvent(steInput)
	if err != nil {
		fmt.Println(err.Error())
	}

	return result.String()
}

func startRun(atl *inspector.CreateAssessmentTemplateOutput, svc inspector.Inspector, InstanceID string) *inspector.StartAssessmentRunOutput {
	// 6. start assessment template run
	ari := &inspector.StartAssessmentRunInput{
		AssessmentRunName:     aws.String(InstanceID + "_Run_" + time.Now().Format("2006-01-02_15.04.05")),
		AssessmentTemplateArn: aws.String(*atl.AssessmentTemplateArn),
	}

	ar, arerr := svc.StartAssessmentRun(ari)
	if arerr != nil {
		fmt.Println(arerr.Error())
	}

	return ar
}

// SetTag changes tag inspector to false
func SetTag(InstanceID *string, tag string) bool {
	ec2Svc := ec2.New(session.New())

	_, errTag := ec2Svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{InstanceID},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("inspector"),
				Value: aws.String(tag),
			},
		},
	})

	if errTag != nil {
		log.Println("Could not create tags for instance", errTag)
		return false
	}

	return true
}

// GetRandomTag will return a random tag everytime it is called
func GetRandomTag() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomTag := r1.Int()
	return strconv.Itoa(randomTag)
}

// Begin will start entire execution
func Begin(InstanceID string, randomTag string) {
	sess, _ := session.NewSession()

	svc := inspector.New(sess)

	// 1. create Resource Group
	rgArn := createResourceGroup(*svc, randomTag)
	fmt.Println("Resource group: ", rgArn)

	// 2. create Assessment Target
	at := createAssessmentTarget(rgArn, *svc, InstanceID)
	fmt.Println("\nAssessment target: ", at)

	// 3. create rules package input
	rp := getRulesPackages(*svc)
	fmt.Println("Rulespackages ARNs:", rp)

	// 4. create assessment template
	atl := createAssessmentTemplate(at, rp, *svc, InstanceID)
	fmt.Println("AssessmentTemplateArn: ", *atl.AssessmentTemplateArn)

	// 5. Subscribe to event - ASSESSMENT_RUN_COMPLETED
	topicArn := os.Getenv("TOPICARN")
	arc := subscribeToEvent(*svc, *atl.AssessmentTemplateArn, topicArn)
	fmt.Println("Response to SubscribeEvent: ", arc)

	// 6. start assessment template run
	startRun(atl, *svc, InstanceID)

}
