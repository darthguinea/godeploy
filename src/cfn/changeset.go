package cfn

import (
	"os"

	"../log"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// func getChangeSetCount() {
// 	svc := cloudformation.New(getSession(r))

// 	cloudformation.DescribeChangeSetInput

// 	svc.DescribeChangeSetRequest()
// }

func createChangeSetFromFile(name string, uri string, params []*cloudformation.Parameter) cloudformation.CreateChangeSetInput {
	file, err := os.Open(uri)
	if err != nil {
		log.Error("Could not open file %v", uri)
	}

	data := make([]byte, 1048576)
	x, _ := file.Read(data)
	templateBody := string(data[:x])

	return cloudformation.CreateChangeSetInput{
		StackName:     &name,
		ChangeSetName: &name,
		TemplateBody:  &templateBody,
		Parameters:    params,
	}
}

func createChangeSetFromURI(name string, uri string, params []*cloudformation.Parameter) cloudformation.CreateChangeSetInput {
	return cloudformation.CreateChangeSetInput{
		StackName:     &name,
		ChangeSetName: &name,
		TemplateURL:   &uri,
		Parameters:    params,
	}
}
