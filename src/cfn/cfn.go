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

func CreateChangeSet(r string, name string, uri string, params []string) {
	// svc := cloudformation.New(getSession(r))

	cfnParams := getParameters(params)
	template := cloudformation.CreateChangeSetInput{}
	if pass, path := parseURI(uri); pass {
		template = createChangeSetFromFile(name, path, cfnParams)
	} else {
		template = createChangeSetFromURI(name, path, cfnParams)
	}
	log.Info("%v", template)

	// svc.CreateChangeSet()
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
func CreateStack(r string, name string, uri string, params []string) {
	svc := cloudformation.New(getSession(r))

	cfnParams := getParameters(params)
	log.Debug("Using Parameters:")
	log.Debug("%v", cfnParams)

	template := cloudformation.CreateStackInput{}
	if pass, path := parseURI(uri); pass {
		template = createStackFromFile(name, path, cfnParams)
	} else {
		template = createStackFromURI(name, path, cfnParams)
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
