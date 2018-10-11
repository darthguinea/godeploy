package cfn

import (
	"os"
	"strings"

	"../log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func getSession(r string) *session.Session {
	return session.New(&aws.Config{
		Region: aws.String(r),
	})
}

func CreateChangeSet(r string,
	currentStack *cloudformation.Stack,
	name string, uri string, params []*cloudformation.Parameter) {
	log.Debug("%v", currentStack)
	getUpdatedParameters(currentStack.Parameters, params)

	// Initialise the variable:
	svc := cloudformation.New(getSession(r))
	template := cloudformation.CreateChangeSetInput{}
	count := getChangeSetCount(r, name)
	if pass, path := parseURI(uri); pass {
		template = createChangeSetFromFile(name, path, count, currentStack.Parameters)
	} else {
		template = createChangeSetFromURI(name, path, count, currentStack.Parameters)
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
		log.Print("%30v %50v %50v", *stack.StackName,
			stack.CreationTime.Format("Mon Jan 2 15:04:05 MST 2006"),
			*stack.StackStatus,
		)
	}
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
func CreateStack(r string, name string, uri string, params []*cloudformation.Parameter) {
	svc := cloudformation.New(getSession(r))

	log.Debug("Using Parameters:")
	log.Debug("%v", params)

	template := cloudformation.CreateStackInput{}
	if pass, path := parseURI(uri); pass {
		template = createStackFromFile(name, path, params)
	} else {
		template = createStackFromURI(name, path, params)
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
