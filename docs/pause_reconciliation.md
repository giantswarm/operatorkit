# Pause Reconciliation

Based on [the upstream CAPI support for pausing annotations] `operatorkit`
supports the upstream pausing annotation as well as its own default.
Additionally you can define your own pausing annotation with [the controller
configuration]. Once the pausing annotation is added to a runtime object,
`operatorkit` will not reconcile this object any further and will not execute
any configured handler for this runtime object anymore until the pausing
annotation is removed again.



[the upstream CAPI support for pausing annotations]: https://cluster-api.sigs.k8s.io/developer/providers/v1alpha2-to-v1alpha3.html#support-the-clusterx-k8siopaused-annotation-and-clusterspecpaused-field
[the controller configuration]: https://pkg.go.dev/github.com/giantswarm/operatorkit@v1.2.0/controller?tab=doc#Config
