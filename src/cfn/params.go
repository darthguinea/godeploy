package cfn

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func getParameters(p []string) []*cloudformation.Parameter {
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
