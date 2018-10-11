package cfn

import (
	"fmt"
	"os"
	"strings"

	"../log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// getUpdatedParameters - will update the current parameters to
// include the new parameters and will set the UsePreviousValue to
// true if the user has requested that param is updated
func getUpdatedParameters(currentParams []*cloudformation.Parameter,
	newParams []*cloudformation.Parameter) {
	for _, currentParam := range currentParams {
		currentParam.ParameterValue = aws.String("")
		currentParam.UsePreviousValue = aws.Bool(true)
	}
outer:
	for _, newParam := range newParams {
		for _, currentParam := range currentParams {
			if strings.Compare(*currentParam.ParameterKey, *newParam.ParameterKey) == 0 {
				currentParam.UsePreviousValue = aws.Bool(false)
				currentParam.ParameterValue = newParam.ParameterValue
				continue outer
			}
		}
		currentParams = append(currentParams, newParam)
	}
	log.Debug("Parameters: %v", currentParams)
	return
}

func getChangeSetCount(r string, name string) int {
	svc := cloudformation.New(getSession(r))
	sets, err := svc.ListChangeSets(&cloudformation.ListChangeSetsInput{
		StackName: &name,
	})
	if err != nil {
		log.Error("%v", err)
	}
	return len(sets.Summaries) + 1
}

func createChangeSetFromFile(name string, uri string, count int,
	params []*cloudformation.Parameter) cloudformation.CreateChangeSetInput {
	file, err := os.Open(uri)
	if err != nil {
		log.Error("Could not open file %v", uri)
	}

	data := make([]byte, 1048576)
	x, _ := file.Read(data)
	templateBody := string(data[:x])

	changeSetName := fmt.Sprintf("%s-%d", name, count)
	return cloudformation.CreateChangeSetInput{
		StackName:     &name,
		ChangeSetName: &changeSetName,
		TemplateBody:  &templateBody,
		Parameters:    params,
	}
}

func createChangeSetFromURI(name string, uri string, count int,
	params []*cloudformation.Parameter) cloudformation.CreateChangeSetInput {
	changeSetName := fmt.Sprintf("%s-%d", name, count)
	return cloudformation.CreateChangeSetInput{
		StackName:     &name,
		ChangeSetName: &changeSetName,
		TemplateURL:   &uri,
		Parameters:    params,
	}
}
