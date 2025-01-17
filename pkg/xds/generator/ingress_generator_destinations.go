package generator

import (
	"reflect"

	"golang.org/x/exp/slices"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	"github.com/kumahq/kuma/pkg/dns"
	meshhttproute_api "github.com/kumahq/kuma/pkg/plugins/policies/meshhttproute/api/v1alpha1"
	"github.com/kumahq/kuma/pkg/util/pointer"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	envoy_tags "github.com/kumahq/kuma/pkg/xds/envoy/tags"
)

type destinations map[string]map[string][]envoy_tags.Tags

func (d destinations) get(mesh string, service string) []envoy_tags.Tags {
	forMesh := d[mesh]
	return append(forMesh[service], forMesh[mesh_proto.MatchAllTag]...)
}

func buildDestinations(ingressProxy *core_xds.ZoneIngressProxy) destinations {
	dest := destinations{}
	availableSvcsByMesh := map[string][]*mesh_proto.ZoneIngress_AvailableService{}
	for _, service := range ingressProxy.ZoneIngressResource.Spec.AvailableServices {
		availableSvcsByMesh[service.Mesh] = append(availableSvcsByMesh[service.Mesh], service)
	}
	for _, meshResources := range ingressProxy.MeshResourceList {
		res := xds_context.Resources{MeshLocalResources: meshResources.Resources}
		destForMesh := map[string][]envoy_tags.Tags{}
		meshHTTPRoutes := res.ListOrEmpty(meshhttproute_api.MeshHTTPRouteType).(*meshhttproute_api.MeshHTTPRouteResourceList).Items
		addTrafficRouteDestinations(res.TrafficRoutes().Items, destForMesh)
		addMeshHTTPRouteDestinations(meshHTTPRoutes, destForMesh)
		addGatewayRouteDestinations(res.GatewayRoutes().Items, destForMesh)
		addMeshGatewayDestinations(res.MeshGateways().Items, destForMesh)
		addVirtualOutboundDestinations(res.VirtualOutbounds().Items, availableSvcsByMesh[meshResources.Mesh.GetMeta().GetName()], destForMesh)
		dest[meshResources.Mesh.GetMeta().GetName()] = destForMesh
	}
	return dest
}

func addMeshGatewayDestinations(
	meshGateways []*core_mesh.MeshGatewayResource,
	destinations map[string][]envoy_tags.Tags,
) {
	for _, meshGateway := range meshGateways {
		for _, selector := range meshGateway.Selectors() {
			addMeshGatewayListenersDestinations(
				meshGateway.Spec,
				selector.GetMatch(),
				destinations,
			)
		}
	}
}

func addMeshGatewayListenersDestinations(
	meshGateway *mesh_proto.MeshGateway,
	matchTags map[string]string,
	destinations map[string][]envoy_tags.Tags,
) {
	service := matchTags[mesh_proto.ServiceTag]

	for _, listener := range meshGateway.GetConf().GetListeners() {
		if !listener.CrossMesh {
			continue
		}

		destinations[service] = append(
			destinations[service],
			mesh_proto.Merge(
				meshGateway.GetTags(),
				matchTags,
				listener.GetTags(),
			),
		)
	}
}

func addGatewayRouteDestinations(
	gatewayRoutes []*core_mesh.MeshGatewayRouteResource,
	destinations map[string][]envoy_tags.Tags,
) {
	var backends []*mesh_proto.MeshGatewayRoute_Backend

	for _, route := range gatewayRoutes {
		for _, rule := range route.Spec.GetConf().GetHttp().GetRules() {
			backends = append(backends, rule.Backends...)
		}

		for _, rule := range route.Spec.GetConf().GetTcp().GetRules() {
			backends = append(backends, rule.Backends...)
		}
	}

	for _, backend := range backends {
		addDestination(backend.Destination, destinations)
	}
}

func addTrafficRouteDestinations(
	policies []*core_mesh.TrafficRouteResource,
	destinations map[string][]envoy_tags.Tags,
) {
	for _, policy := range policies {
		for _, split := range policy.Spec.Conf.GetSplitWithDestination() {
			addDestination(split.Destination, destinations)
		}

		for _, http := range policy.Spec.Conf.Http {
			for _, split := range http.GetSplitWithDestination() {
				addDestination(split.Destination, destinations)
			}
		}
	}
}

func addMeshHTTPRouteDestinations(
	policies []*meshhttproute_api.MeshHTTPRouteResource,
	destinations map[string][]envoy_tags.Tags,
) {
	addTrafficFlowByDefaultDestinationIfMeshHTTPRoutesExist(policies, destinations)

	// Note that we're not merging these resources, but that's OK because the
	// set of destinations after merging is a subset of the set we get here by
	// iterating through them.
	for _, policy := range policies {
		for _, to := range policy.Spec.To {
			if toTags, ok := tagsFromTargetRef(to.TargetRef); ok {
				addMeshHTTPRouteToDestinations(to.Rules, toTags, destinations)
			}
		}
	}
}

func addMeshHTTPRouteToDestinations(
	rules []meshhttproute_api.Rule,
	toTags envoy_tags.Tags,
	destinations map[string][]envoy_tags.Tags,
) {
	for _, rule := range rules {
		if rule.Default.BackendRefs == nil {
			addDestination(toTags, destinations)
			continue
		}

		for _, backendRef := range pointer.Deref(rule.Default.BackendRefs) {
			if tags, ok := tagsFromTargetRef(backendRef.TargetRef); ok {
				addDestination(tags, destinations)
			}
		}
	}
}

func addDestination(tags map[string]string, destinations map[string][]envoy_tags.Tags) {
	service := tags[mesh_proto.ServiceTag]
	destinations[service] = append(destinations[service], tags)
}

// addTrafficFlowByDefaultDestinationIfMeshHTTPRoutesExist makes sure that when
// at least one MeshHTTPRoute policy exists there will be a "match all"
// destination pointing to all services (kuma.io/service:* -> kuma.io/service:*)
// This logic is necessary because of conflicting behaviors of TrafficRoute and
// MeshHTTPRoute policies. TrafficRoute expects that by default traffic doesn't
// flow, and there is necessary TrafficRoute with appropriate configuration
// to make communication between services possible. MeshHTTPRoute on the other
// hand expects the traffic to flow by default. As a result, when there is
// at least one MeshHTTPRoute policy present, traffic between services will flow
// by default, when there is none, it will flow, when appropriate TrafficRoute
// policy will exist.
func addTrafficFlowByDefaultDestinationIfMeshHTTPRoutesExist(
	policies []*meshhttproute_api.MeshHTTPRouteResource,
	destinations map[string][]envoy_tags.Tags,
) {
	// If there are no MeshHTTPRoutes, we are not modifying destinations
	if len(policies) == 0 {
		return
	}

	// We need to add a destination to route any service to any instance of
	// that service
	matchAllTags := envoy_tags.Tags{mesh_proto.ServiceTag: mesh_proto.MatchAllTag}
	matchAllDestinations := destinations[mesh_proto.MatchAllTag]
	foundAllServicesDestination := slices.ContainsFunc(
		matchAllDestinations,
		func(tagsElem envoy_tags.Tags) bool {
			return reflect.DeepEqual(tagsElem, matchAllTags)
		},
	)

	if !foundAllServicesDestination {
		matchAllDestinations = append(matchAllDestinations, matchAllTags)
	}

	destinations[mesh_proto.MatchAllTag] = matchAllDestinations
}

func addVirtualOutboundDestinations(
	virtualOutbounds []*core_mesh.VirtualOutboundResource,
	availableServices []*mesh_proto.ZoneIngress_AvailableService,
	destinations map[string][]envoy_tags.Tags,
) {
	// If there are no VirtualOutbounds, we are not modifying destinations
	if len(virtualOutbounds) == 0 {
		return
	}

	for _, availableService := range availableServices {
		for _, matched := range dns.Match(virtualOutbounds, availableService.Tags) {
			service := availableService.Tags[mesh_proto.ServiceTag]
			tags := envoy_tags.Tags{}
			for _, param := range matched.Spec.GetConf().GetParameters() {
				tags[param.TagKey] = availableService.Tags[param.TagKey]
			}
			destinations[service] = append(destinations[service], tags)
		}
	}
}
