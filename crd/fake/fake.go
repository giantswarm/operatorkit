// Package fake provides primitives to create fake CRDs and corresponding CROs
// like the following example.
//
// CRD
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: tests.example.com
//     spec:
//       group: example.com
//       version: v1
//       scope: Cluster
//       names:
//         plural: tests
//         singular: test
//         kind: Test
//
// CRO
//
//     apiVersion: "example.com/v1"
//     kind: Test
//     metadata:
//       name: al9qy
//     spec:
//       id: al9qy
//
package fake
