package osclient

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"

	OBIIIConfig "github.com/ritsec/ops-bot-iii/config"
)

var (
	// identityClient is the global openstack identity client
	identityClient *gophercloud.ServiceClient
	// networkClient is the global openstack network client
	networkClient *gophercloud.ServiceClient
	// storageClient is the global openstack storage client
	storageClient *gophercloud.ServiceClient
	// computeClient is the global openstack compute client
	computeClient *gophercloud.ServiceClient

	// tokenExpiry is the expiry time of the token
	tokenExpiry time.Time
)

// Init() will call setUpClients() to make the clients during the compile
func init() {
	setUpClients()
}

func setUpClients() {
	// Checks to see if the config has the Openstack enabled option enabled
	if OBIIIConfig.Openstack.Enabled {
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

		tokenExpiry = time.Now().Add(time.Hour) // assuming 1-hour token validity
	}
}

// RefreshClientsIfNeeded() will check if the token has expired and refresh the clients
// To be used before making any openstack API calls
func RefreshClientsIfNeeded() {
	if time.Now().After(tokenExpiry) {
		setUpClients()
	}
}
