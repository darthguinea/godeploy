package cfn

import (
	"strings"

	"github.com/darthguinea/golib/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// getUpdatedParameters - will update the current parameters to
// include the new parameters and will set the UsePreviousValue to
// true if the user has requested that param is updated
func getUpdatedParameters(currentParams *[]*cloudformation.Parameter,
	newParams []*cloudformation.Parameter) {
	for _, currentParam := range *currentParams {
		currentParam.ParameterValue = aws.String("")
		currentParam.UsePreviousValue = aws.Bool(true)
	}
outer:
	for _, newParam := range newParams {
		if isNewParameter(newParam.ParameterKey, *currentParams) {
			newParam.UsePreviousValue = aws.Bool(false)
		}
		for _, currentParam := range *currentParams {
			if strings.Compare(*currentParam.ParameterKey, *newParam.ParameterKey) == 0 {
				currentParam.UsePreviousValue = aws.Bool(false)
				currentParam.ParameterValue = newParam.ParameterValue
				continue outer
			}
		}
		*currentParams = append(*currentParams, newParam)
	}
	log.Debug("Parameters: %v", &currentParams)
	log.Info("Using:")
	for _, param := range *currentParams {
		if *param.UsePreviousValue {
			log.Info("\t%v=UsePreviousValue", *param.ParameterKey)
		} else {
			log.Info("\t%v=%v", *param.ParameterKey, *param.ParameterValue)
		}
	}
	return
}

// isNewParameter - This function will check to see if the parameter exists in the
// current aws stack, returns true or false
func isNewParameter(newParamKey *string, currentParams []*cloudformation.Parameter) bool {
	for _, currentParam := range currentParams {
		if strings.Compare(*currentParam.ParameterKey, *newParamKey) == 0 {
			return false
		}
	}
	return true
}
