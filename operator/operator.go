package operator

import (
	microerror "github.com/giantswarm/microkit/error"
)

// ProcessCreate is a drop-in for an informer's AddFunc. It receives the custom
// object observed during TPR watches and anything that implements Operator.
// ProcessCreate takes care about all necessary reconciliation logic for create
// events.
//
//     func addFunc(obj interface{}) {
//         err := ProcessCreate(obj, operator)
//         if err != nil {
//             // error handling here
//         }
//     }
//
//     newResourceEventHandler := &cache.ResourceEventHandlerFuncs{
//         AddFunc:    addFunc,
//     }
//
func ProcessCreate(obj interface{}, operator Operator) error {
	currentState, err := operator.GetCurrentState(obj)
	if err != nil {
		return microerror.MaskAny(err)
	}

	desiredState, err := operator.GetDesiredState(obj)
	if err != nil {
		return microerror.MaskAny(err)
	}

	createState, err := operator.GetCreateState(obj, currentState, desiredState)
	if err != nil {
		return microerror.MaskAny(err)
	}

	err = operator.ProcessCreateState(obj, createState)
	if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

// ProcessDelete is a drop-in for an informer's DeleteFunc. It receives the
// custom object observed during TPR watches and anything that implements
// Operator. ProcessDelete takes care about all necessary reconciliation logic
// for delete events.
//
//     func deleteFunc(obj interface{}) {
//         err := ProcessDelete(obj, operator)
//         if err != nil {
//             // error handling here
//         }
//     }
//
//     newResourceEventHandler := &cache.ResourceEventHandlerFuncs{
//         DeleteFunc:    deleteFunc,
//     }
//
func ProcessDelete(obj interface{}, operator Operator) error {
	currentState, err := operator.GetCurrentState(obj)
	if err != nil {
		return microerror.MaskAny(err)
	}

	desiredState, err := operator.GetDesiredState(obj)
	if err != nil {
		return microerror.MaskAny(err)
	}

	deleteState, err := operator.GetDeleteState(obj, currentState, desiredState)
	if err != nil {
		return microerror.MaskAny(err)
	}

	err = operator.ProcessDeleteState(obj, deleteState)
	if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}
