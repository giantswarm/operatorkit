package client

import (
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2017-10-01/dns"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-06-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/giantswarm/microerror"
)

type AzureClientSetConfig struct {
	// ClientID is the ID of the Active Directory Service Principal.
	ClientID string
	// ClientSecret is the secret of the Active Directory Service Principal.
	ClientSecret string
	// The cloud environment identifier. Takes values from
	// https://github.com/Azure/go-autorest/blob/ec5f4903f77ed9927ac95b19ab8e44ada64c1356/autorest/azure/environments.go#L13
	// TODO rename to EnvironmentName. See https://github.com/giantswarm/giantswarm/issues/4124.
	Cloud string
	// SubscriptionID is the ID of the Azure subscription.
	SubscriptionID string
	// TenantID is the ID of the Active Directory tenant.
	TenantID string
}

// clientConfig contains all essential information to create an Azure client.
type clientConfig struct {
	subscriptionID          string
	resourceManagerEndpoint string
	servicePrincipalToken   *adal.ServicePrincipalToken
}

func (c AzureClientSetConfig) Validate() error {
	if c.ClientID == "" {
		return errors.New("ClientID must not be empty")
	}
	if c.ClientSecret == "" {
		return errors.New("ClientSecret must not be empty")
	}
	if c.SubscriptionID == "" {
		return errors.New("SubscriptionID must not be empty")
	}
	if c.TenantID == "" {
		return errors.New("TenantID must not be empty")
	}

	return nil
}

// AzureClientSet is the collection of Azure API clients.
type AzureClientSet struct {
	// DeploymentsClient manages deployments of ARM templates.
	DeploymentsClient *resources.DeploymentsClient
	// GroupsClient manages ARM resource groups.
	GroupsClient *resources.GroupsClient
	// DNSRecordSetsClient manages DNS zones' records.
	DNSRecordSetsClient *dns.RecordSetsClient
	// DNSRecordSetsClient manages DNS zones.
	DNSZonesClient *dns.ZonesClient
	// InterfacesClient manages virtual network interfaces.
	InterfacesClient *network.InterfacesClient
	//StorageAccountsClient manages blobs in storage containers.
	StorageAccountsClient *storage.AccountsClient
	// VirtualNetworkClient manages virtual networks.
	VirtualNetworkClient *network.VirtualNetworksClient
	// VirtualNetworkGatewayConnectionsClient manages virtual network gateway connections.
	VirtualNetworkGatewayConnectionsClient *network.VirtualNetworkGatewayConnectionsClient
	// VirtualNetworkGatewaysClient manages virtual network gateways.
	VirtualNetworkGatewaysClient *network.VirtualNetworkGatewaysClient
	// VirtualMachineScaleSetsClient manages virtual machine scale sets.
	VirtualMachineScaleSetsClient *compute.VirtualMachineScaleSetsClient
	// VirtualMachineScaleSetVMsClient manages virtual machine scale set VMs.
	VirtualMachineScaleSetVMsClient *compute.VirtualMachineScaleSetVMsClient
	// VnetPeeringClient manages virtual network peerings.
	VnetPeeringClient *network.VirtualNetworkPeeringsClient
}

// NewAzureClientSet returns the Azure API clients.
func NewAzureClientSet(config AzureClientSetConfig) (*AzureClientSet, error) {
	if err := config.Validate(); err != nil {
		return nil, microerror.Maskf(invalidConfigError, "config.%s", err)
	}

	// Returns environment object contains all API endpoints for specific Azure cloud.
	// For empty config.Cloud returns Azure public cloud.
	env, err := parseAzureEnvironment(config.Cloud)
	if err != nil {
		return nil, microerror.Maskf(err, "parsing Azure environment")
	}

	servicePrincipalToken, err := newServicePrincipalToken(config, env)
	if err != nil {
		return nil, microerror.Maskf(err, "creating service principal token")
	}

	c := &clientConfig{
		subscriptionID:          config.SubscriptionID,
		resourceManagerEndpoint: env.ResourceManagerEndpoint,
		servicePrincipalToken:   servicePrincipalToken,
	}

	clientSet := &AzureClientSet{
		DeploymentsClient:                      newDeploymentsClient(c),
		GroupsClient:                           newGroupsClient(c),
		DNSRecordSetsClient:                    newDNSRecordSetsClient(c),
		DNSZonesClient:                         newDNSZonesClient(c),
		InterfacesClient:                       newInterfacesClient(c),
		StorageAccountsClient:                  newStorageAccountsClient(c),
		VirtualNetworkClient:                   newVirtualNetworkClient(c),
		VirtualNetworkGatewayConnectionsClient: newVirtualNetworkGatewayConnectionsClient(c),
		VirtualNetworkGatewaysClient:           newVirtualNetworkGatewaysClient(c),
		VirtualMachineScaleSetVMsClient:        newVirtualMachineScaleSetVMsClient(c),
		VirtualMachineScaleSetsClient:          newVirtualMachineScaleSetsClient(c),
		VnetPeeringClient:                      newVnetPeeringClient(c),
	}

	return clientSet, nil
}

// ResponseWasNotFound returns true if the response code from the Azure API
// was a 404.
func ResponseWasNotFound(resp autorest.Response) bool {
	if resp.Response != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}

	return false
}

func newDeploymentsClient(config *clientConfig) *resources.DeploymentsClient {
	c := resources.NewDeploymentsClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newGroupsClient(config *clientConfig) *resources.GroupsClient {
	c := resources.NewGroupsClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newDNSRecordSetsClient(config *clientConfig) *dns.RecordSetsClient {
	c := dns.NewRecordSetsClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newDNSZonesClient(config *clientConfig) *dns.ZonesClient {
	c := dns.NewZonesClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newInterfacesClient(config *clientConfig) *network.InterfacesClient {
	c := network.NewInterfacesClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newVirtualNetworkClient(config *clientConfig) *network.VirtualNetworksClient {
	c := network.NewVirtualNetworksClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newVirtualNetworkGatewayConnectionsClient(config *clientConfig) *network.VirtualNetworkGatewayConnectionsClient {
	c := network.NewVirtualNetworkGatewayConnectionsClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newVirtualNetworkGatewaysClient(config *clientConfig) *network.VirtualNetworkGatewaysClient {
	c := network.NewVirtualNetworkGatewaysClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newVirtualMachineScaleSetsClient(config *clientConfig) *compute.VirtualMachineScaleSetsClient {
	c := compute.NewVirtualMachineScaleSetsClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newVirtualMachineScaleSetVMsClient(config *clientConfig) *compute.VirtualMachineScaleSetVMsClient {
	c := compute.NewVirtualMachineScaleSetVMsClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newVnetPeeringClient(config *clientConfig) *network.VirtualNetworkPeeringsClient {
	c := network.NewVirtualNetworkPeeringsClientWithBaseURI(config.resourceManagerEndpoint, config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)

	return &c
}

func newServicePrincipalToken(config AzureClientSetConfig, env azure.Environment) (*adal.ServicePrincipalToken, error) {
	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, config.TenantID)
	if err != nil {
		return nil, microerror.Maskf(err, "creating OAuth config")
	}

	token, err := adal.NewServicePrincipalToken(*oauthConfig, config.ClientID, config.ClientSecret, env.ServiceManagementEndpoint)
	if err != nil {
		return nil, microerror.Maskf(err, "getting token")
	}

	return token, nil
}

func newStorageAccountsClient(config *clientConfig) *storage.AccountsClient {
	c := storage.NewAccountsClient(config.subscriptionID)
	c.Authorizer = autorest.NewBearerAuthorizer(config.servicePrincipalToken)
	return &c
}

// parseAzureEnvironment returns azure environment by name.
func parseAzureEnvironment(cloudName string) (azure.Environment, error) {
	if cloudName == "" {
		return azure.PublicCloud, nil
	}

	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		return azure.Environment{}, microerror.Mask(err)
	}

	return env, nil
}
