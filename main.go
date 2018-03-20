package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/inspector"
)

func createResourceGroup(svc inspector.Inspector) string {
	// 1. create resource group input
	rgi := &inspector.CreateResourceGroupInput{
		ResourceGroupTags: []*inspector.ResourceGroupTag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("inspector-test"),
			},
		},
	}

	// 2. create resource group
	rg, rgerr := svc.CreateResourceGroup(rgi)

	if rgerr != nil {
		print(rgerr)
	}
	// print("Resource group: ", *rg.ResourceGroupArn)

	return *rg.ResourceGroupArn
}

func createAssessmentTarget(rgArn string, svc inspector.Inspector) string {
	ati := &inspector.CreateAssessmentTargetInput{
		AssessmentTargetName: aws.String("ExampleAssessmentTarget"),
		ResourceGroupArn:     &rgArn,
	}

	at, aterr := svc.CreateAssessmentTarget(ati)

	if aterr != nil {
		fmt.Println(aterr)
	}

	return *at.AssessmentTargetArn

	// fmt.Println("\nAssessment target: ", *at.AssessmentTargetArn)
}

func getRulesPackages(svc inspector.Inspector) *inspector.ListRulesPackagesOutput {
	rpi := &inspector.ListRulesPackagesInput{
		MaxResults: aws.Int64(123),
	}

	rp, erp := svc.ListRulesPackages(rpi)
	if erp != nil {
		fmt.Println(erp.Error())
	}

	return rp
}

func createAssessmentTemplate(at string, rp *inspector.ListRulesPackagesOutput, svc inspector.Inspector) *inspector.CreateAssessmentTemplateOutput {
	atli := &inspector.CreateAssessmentTemplateInput{
		AssessmentTargetArn:    aws.String(at),
		AssessmentTemplateName: aws.String("ExampleAssessmentTemplate"),
		DurationInSeconds:      aws.Int64(3600),
		RulesPackageArns:       rp.RulesPackageArns,
		UserAttributesForFindings: []*inspector.Attribute{
			{
				Key:   aws.String("Example"),
				Value: aws.String("example"),
			},
		},
	}

	atl, atlerr := svc.CreateAssessmentTemplate(atli)

	if atlerr != nil {
		fmt.Println(atlerr)
	}

	return atl
}

func startRun(atl *inspector.CreateAssessmentTemplateOutput, svc inspector.Inspector) *inspector.StartAssessmentRunOutput {
	ari := &inspector.StartAssessmentRunInput{
		AssessmentRunName:     aws.String("examplerun"),
		AssessmentTemplateArn: aws.String(*atl.AssessmentTemplateArn),
	}

	ar, arerr := svc.StartAssessmentRun(ari)
	if arerr != nil {
		fmt.Println(arerr.Error())
	}

	return ar
}

func main() {
	sess, _ := session.NewSession()
	svc := inspector.New(sess, &aws.Config{
		Region: aws.String("us-west-2"),
	})

	// 1. create Resource Group
	rgArn := createResourceGroup(*svc)
	fmt.Println("Resource group: ", rgArn)

	// 2. create Assessment Target
	at := createAssessmentTarget(rgArn, *svc)
	fmt.Println("\nAssessment target: ", at)

	// 3. create rules package input
	rp := getRulesPackages(*svc)
	fmt.Println("Rulespackages ARNs:", rp)

	// 4. create assessment template
	atl := createAssessmentTemplate(at, rp, *svc)
	fmt.Println("AssessmentTemplateArn: ", *atl.AssessmentTemplateArn)

	// 5. start assessment template run
	ar := startRun(atl, *svc)

	fmt.Println(ar)

}
