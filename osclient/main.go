package osclient

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"

	OBIIIConfig "github.com/ritsec/ops-bot-iii/config"
)

var (
	// networkClient is the global openstack identity client
	identityClient *gophercloud.ServiceClient
	networkClient  *gophercloud.ServiceClient
	storageClient  *gophercloud.ServiceClient
	computeClient  *gophercloud.ServiceClient
)

func init() {
	ctx := context.Background()
	ao, eo, tlsConfig, err := clouds.Parse(clouds.WithCloudName("openstack"), clouds.WithLocations(OBIIIConfig.Openstack.CloudsPath))
	if err != nil {
		log.Fatalf("Failed to parse the clouds.yaml: %v", err)
	}

	providerClient, err := config.NewProviderClient(ctx, ao, config.WithTLSConfig(tlsConfig))
	if err != nil {
		log.Fatalf("Failed to make providerClient with NewProviderClient: %v", err)
	}

	_identityClient, err := openstack.NewIdentityV3(providerClient, eo)
	if err != nil {
		log.Fatalf("Failed to make _identityClient with NewIdentityV3: %v", err)
	}

	_networkClient, err := openstack.NewNetworkV2(providerClient, eo)
	if err != nil {
		log.Fatalf("Failed to make _networkClient with NewNetworkV2: %v", err)
	}

	_storageClient, err := openstack.NewBlockStorageV3(providerClient, eo)
	if err != nil {
		log.Fatalf("Failed to make _storageClient with NewBlockStorageV3: %v", err)
	}

	_computeClient, err := openstack.NewComputeV2(providerClient, eo)
	if err != nil {
		log.Fatalf("Failed to make _computeClient with NewComputeV2: %v", err)
	}

	identityClient = _identityClient
	networkClient = _networkClient
	storageClient = _storageClient
	computeClient = _computeClient
}
