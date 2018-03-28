// +build k8srequired

package reconciliation

import (
	"testing"

	"github.com/giantswarm/operatorkit/framework/integration/client"
)

// Test_Finalizer_Integration_Reconciliation is a integration test for
// the proper replay and reconciliation of delete events with finalizers.
func Test_Finalizer_Integration_Reconciliation(t *testing.T) {
	testNamespace := "finalizer-integration-reconciliation-test"

	client.MustSetup(testNamespace)
	defer client.MustTeardown(testNamespace)

	// TODO: Implement the actual test here.

}
