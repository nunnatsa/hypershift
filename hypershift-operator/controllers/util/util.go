package util

import (
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"

	"k8s.io/apimachinery/pkg/types"
)

const HypershiftRouteLabel = "hypershift.openshift.io/hosted-control-plane"

// ParseNamespacedName expects a string with the format "namespace/name"
// and returns the proper types.NamespacedName.
// This is useful when watching a CR annotated with the format above to requeue the CR
// described in the annotation.
func ParseNamespacedName(name string) types.NamespacedName {
	parts := strings.SplitN(name, string(types.Separator), 2)
	if len(parts) > 1 {
		return types.NamespacedName{Namespace: parts[0], Name: parts[1]}
	}
	return types.NamespacedName{Name: parts[0]}
}

// NoopReconcile is just a default mutation function that does nothing.
var NoopReconcile controllerutil.MutateFn = func() error { return nil }
