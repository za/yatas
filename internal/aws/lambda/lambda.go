package lambda

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	lambdas := GetLambdas(s)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_LMD_001", CheckIfLambdaPrivate)(checkConfig, lambdas, "AWS_LMD_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_LMD_002", CheckIfLambdaInSecurityGroup)(checkConfig, lambdas, "AWS_LMD_002")
	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
