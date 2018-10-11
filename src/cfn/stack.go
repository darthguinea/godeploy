package cfn

import (
	"os"

	"../log"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func createStackFromFile(name string, uri string, params []*cloudformation.Parameter, capabilities []*string) cloudformation.CreateStackInput {
	file, err := os.Open(uri)
	if err != nil {
		log.Error("Could not open file %v", uri)
	}

	data := make([]byte, 1048576)
	x, _ := file.Read(data)
	templateBody := string(data[:x])

	return cloudformation.CreateStackInput{
		StackName:    &name,
		TemplateBody: &templateBody,
		Parameters:   params,
		Capabilities: capabilities,
	}
}

func createStackFromURI(name string, uri string, params []*cloudformation.Parameter, capabilities []*string) cloudformation.CreateStackInput {
	return cloudformation.CreateStackInput{
		StackName:    &name,
		TemplateURL:  &uri,
		Parameters:   params,
		Capabilities: capabilities,
	}
}
