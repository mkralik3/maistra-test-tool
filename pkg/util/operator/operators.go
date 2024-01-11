package operator

import (
	"fmt"
	"strings"
	"time"

	"github.com/maistra/maistra-test-tool/pkg/util/oc"
	"github.com/maistra/maistra-test-tool/pkg/util/pod"
	"github.com/maistra/maistra-test-tool/pkg/util/retry"
	"github.com/maistra/maistra-test-tool/pkg/util/shell"
	"github.com/maistra/maistra-test-tool/pkg/util/test"
)

func GetCsvName(t test.TestHelper, operatorNamespace string, partialName string) string {
	output := shell.Execute(t, fmt.Sprintf(`oc get csv -n %s -o custom-columns="NAME:.metadata.name" |grep %s ||true`, operatorNamespace, partialName))
	return strings.Trim(output, "\n")
}

func WaitForOperatorReady(t test.TestHelper, operatorNamespace string, operatorSelector string, csvName string) {
	t.Logf("Waiting for operator %s to succeed", csvName)
	// When the operator is installed, the CSV take some time to be created, need to wait until is created to validate the phase
	retry.UntilSuccessWithOptions(t, retry.Options().DelayBetweenAttempts(5*time.Second).MaxAttempts(70), func(t test.TestHelper) {
		if GetCsvName(t, operatorNamespace, csvName) == "" {
			t.Errorf("Operator %s is not yet installed", csvName)
		}
	})

	oc.WaitForPhase(t, operatorNamespace, "csv", csvName, "Succeeded")
	oc.WaitPodReadyWithOptions(t, retry.Options().MaxAttempts(70).DelayBetweenAttempts(5*time.Second), pod.MatchingSelector(operatorSelector, operatorNamespace))
}
