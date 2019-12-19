package adobe

import (
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_api_v2_route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	ingressroutev1 "github.com/projectcontour/contour/apis/contour/v1beta1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ignore properties added/removed by our customization
var ignoreProperties = []cmp.Option{
	cmpopts.IgnoreFields(v2.Cluster{}, "CircuitBreakers", "DrainConnectionsOnHostRemoval"),
	cmpopts.IgnoreFields(envoy_api_v2_core.HealthCheck_HttpHealthCheck{}, "ExpectedStatuses"),
	cmpopts.IgnoreFields(envoy_api_v2_route.RouteAction{}, "RetryPolicy", "Timeout", "HashPolicy"),
	cmpopts.IgnoreFields(envoy_api_v2_route.VirtualHost{}, "RetryPolicy"),
}

// list of tests to ignore (assuming names are unique across suites)
var ignoreTests = map[string]struct{}{
	"ingressroute root delegates to another ingressroute root":      {}, // root to root delegation
	"root ingress delegating to another root w/ different hostname": {}, // root to root delegation
	"self-edge produces a cycle":                                    {}, // root to root delegation
	"multiple tls ingress with secrets should be sorted":            {}, // we group the filter chains together
}

func IgnoreFields() []cmp.Option {
	return ignoreProperties
}

func ShouldSkipTest(name string) bool {
	if _, ignore := ignoreTests[name]; ignore {
		return true
	}
	return false
}

// Object resources
func AdobefyObject(data interface{}) {
	switch obj := data.(type) {
	case *v1beta1.Ingress:
		addClassAnnotation(&obj.ObjectMeta)
	case *ingressroutev1.IngressRoute:
		addClassAnnotation(&obj.ObjectMeta)
	}
}

func addClassAnnotation(om *metav1.ObjectMeta) {
	if metav1.HasAnnotation(*om, "kubernetes.io/ingress.class") {
		return
	}
	metav1.SetMetaDataAnnotation(om, "kubernetes.io/ingress.class", "contour")
}
