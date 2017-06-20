package operator

import (
	microerror "github.com/giantswarm/microkit/error"
)

// ProcessCreate is a drop-in for an informer's AddFunc. It receives the custom
// object observed during TPR watches and anything that implements Resource.
// ProcessCreate takes care about all necessary reconciliation logic for create
// events.
//
//     func addFunc(obj interface{}) {
//         err := ProcessCreate(obj, resources...)
//         if err != nil {
//             // error handling here
//         }
//     }
//
//     newResourceEventHandler := &cache.ResourceEventHandlerFuncs{
//         AddFunc:    addFunc,
//     }
//
func ProcessCreate(obj interface{}, resources ...Resource) error {
	if len(resources) == 0 {
		return microerror.MaskAnyf(executionFailedError, "resources must not be empty")
	}

	for _, r := range resources {
		currentState, err := r.GetCurrentState(obj)
		if err != nil {
			return microerror.MaskAny(err)
		}

		desiredState, err := r.GetDesiredState(obj)
		if err != nil {
			return microerror.MaskAny(err)
		}

		createState, err := r.GetCreateState(obj, currentState, desiredState)
		if err != nil {
			return microerror.MaskAny(err)
		}

		err = r.ProcessCreateState(obj, createState)
		if err != nil {
			return microerror.MaskAny(err)
		}
	}

	return nil
}

// ProcessDelete is a drop-in for an informer's DeleteFunc. It receives the
// custom object observed during TPR watches and anything that implements
// Resource. ProcessDelete takes care about all necessary reconciliation logic
// for delete events.
//
//     func deleteFunc(obj interface{}) {
//         err := ProcessDelete(obj, resources...)
//         if err != nil {
//             // error handling here
//         }
//     }
//
//     newResourceEventHandler := &cache.ResourceEventHandlerFuncs{
//         DeleteFunc:    deleteFunc,
//     }
//
func ProcessDelete(obj interface{}, resources ...Resource) error {
	if len(resources) == 0 {
		return microerror.MaskAnyf(executionFailedError, "resources must not be empty")
	}

	for _, r := range resources {
		currentState, err := r.GetCurrentState(obj)
		if err != nil {
			return microerror.MaskAny(err)
		}

		desiredState, err := r.GetDesiredState(obj)
		if err != nil {
			return microerror.MaskAny(err)
		}

		deleteState, err := r.GetDeleteState(obj, currentState, desiredState)
		if err != nil {
			return microerror.MaskAny(err)
		}

		err = r.ProcessDeleteState(obj, deleteState)
		if err != nil {
			return microerror.MaskAny(err)
		}
	}

	return nil
}
