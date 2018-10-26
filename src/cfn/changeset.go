package cfn

import (
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
