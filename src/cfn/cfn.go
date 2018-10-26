package cfn

import (
	"fmt"
	"os"
	"strings"

	"../log"
	"../utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// Creates an AWS Session if you pass in the region
func getSession(r string) *session.Session {
	return session.New(&aws.Config{
		Region: aws.String(r),
	})
}

// UpdateStack - This will update the cfn stack
func UpdateStack(r string, currentStack *cloudformation.Stack,
	uri string, name string, params []*cloudformation.Parameter,
	capabilities []*string) {
	svc := cloudformation.New(getSession(r))
	log.Info("UpdateStack")

	getUpdatedParameters(currentStack.Parameters, params)
	template := cloudformation.UpdateStackInput{
		StackName:    &name,
		Parameters:   currentStack.Parameters,
		Capabilities: capabilities,
	}

	if pass, path := parseURI(uri); pass {
		templateBody := utils.LoadTemplate(path)
		template.TemplateBody = &templateBody
	} else {
		template.TemplateURL = &path
	}

	stack, err := svc.UpdateStack(&template)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("%v", stack)
}

// CreateChangeSet - This will create a change set for a given stack
func CreateChangeSet(r string, currentStack *cloudformation.Stack, name string, uri string, params []*cloudformation.Parameter, capabilities []*string) {
	log.Debug("%v", currentStack)

	// Initialise the variable:
	svc := cloudformation.New(getSession(r))

	count := len(DescribeChangeSets(r, name).Summaries) + 1
	changeSetName := fmt.Sprintf("%s-%d", name, count)

	getUpdatedParameters(currentStack.Parameters, params)
	template := cloudformation.CreateChangeSetInput{
		StackName:     &name,
		ChangeSetName: &changeSetName,
		Parameters:    currentStack.Parameters,
		Capabilities:  capabilities,
	}

	if pass, path := parseURI(uri); pass {
		templateBody := utils.LoadTemplate(path)
		template.TemplateBody = &templateBody
	} else {
		template.TemplateURL = &path
	}
	changeSet, err := svc.CreateChangeSet(&template)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("%v", changeSet)
}

// DescribeStacks describes the cloudformation stacks in the region
func DescribeStacks(r string) {
	svc := cloudformation.New(getSession(r))

	statii := []*string{
		aws.String("CREATE_COMPLETE"),
		aws.String("CREATE_IN_PROGRESS"),
		aws.String("UPDATE_COMPLETE"),
		aws.String("UPDATE_IN_PROGRESS"),
	}

	stackSummaries, err := svc.ListStacks(&cloudformation.ListStacksInput{
		StackStatusFilter: statii,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, stack := range stackSummaries.StackSummaries {
		log.Print("%-52v %-40v %-20v", *stack.StackName,
			stack.CreationTime.Format("Mon Jan 2 15:04:05 MST 2006"),
			*stack.StackStatus,
		)
		// List change sets:
		changeSets := DescribeChangeSets(r, *stack.StackName)
		for _, change := range changeSets.Summaries {
			log.Print("\tchange set -> %-30v %-40v %-20v", *change.ChangeSetName,
				change.CreationTime.Format("Mon Jan 2 15:04:05 MST 2006"),
				*change.ExecutionStatus,
			)
		}
	}
}

// Describes the ChangeSets for a given Stack in a region
func DescribeChangeSets(r string, name string) *cloudformation.ListChangeSetsOutput {
	svc := cloudformation.New(getSession(r))
	sets, err := svc.ListChangeSets(&cloudformation.ListChangeSetsInput{
		StackName: &name,
	})
	if err != nil {
		log.Error("%v", err)
	}
	return sets
}

// StackExists - Find stack with name, return stack information
func StackExists(r string, name string) (bool, *cloudformation.Stack) {
	svc := cloudformation.New(getSession(r))
	log.Debug("Getting stack information for %v", name)

	stackDetails, _ := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(name),
	})

	if len(stackDetails.Stacks) == 1 {
		return true, stackDetails.Stacks[0]
	}

	return false, nil
}

// CreateStack API call to create aws cloudformation stack
func CreateStack(r string, name string, uri string, params []*cloudformation.Parameter, capabilities []*string) {
	svc := cloudformation.New(getSession(r))
	log.Info("CreateStack")

	log.Debug("Using Parameters:")
	log.Debug("%v", params)

	template := cloudformation.CreateStackInput{
		StackName:    &name,
		Parameters:   params,
		Capabilities: capabilities,
	}
	if pass, path := parseURI(uri); pass {
		templateBody := utils.LoadTemplate(path)
		template.TemplateBody = &templateBody
	} else {
		template.TemplateURL = &path
	}

	stack, err := svc.CreateStack(&template)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("%v", stack)
}

func parseURI(uri string) (bool, string) {
	if strings.HasPrefix(uri, "s3://") {
		return false, uri
	} else if strings.HasPrefix(uri, "file://") {
		uri = uri[7:]
	}

	if _, err := os.Stat(uri); os.IsNotExist(err) {
		log.Error("Could not locate file")
		os.Exit(1)
	}
	return true, uri
}

// GetParameters - Convers an array of strings to array of parameters
func GetParameters(p []string) []*cloudformation.Parameter {
	var parameters []*cloudformation.Parameter
	for _, val := range p {
		strKeyPair := strings.Split(val, "=")
		parameters = append(parameters, &cloudformation.Parameter{
			ParameterKey:   &strKeyPair[0],
			ParameterValue: &strKeyPair[1],
		})
	}
	return parameters
}

// GetCapabilities - Will return an array of string pointers for capabilities
func GetCapabilities(cap string) []*string {
	x := []*string{}
	if strings.Compare(cap, "") == 0 {
		return nil
	}
	for _, c := range strings.Split(cap, ",") {
		x = append(x, &c)
	}
	return x
}
