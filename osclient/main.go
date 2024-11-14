package osclient

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"

	OBIIIConfig "github.com/ritsec/ops-bot-iii/config"
)

var (
	// networkClient is the global openstack identity client
	identityClient *gophercloud.ServiceClient
)

func init() {
	if OBIIIConfig.OpenstackEnabled {
		ctx := context.Background()
		ao, eo, _, err := clouds.Parse(clouds.WithCloudName("openstack"), clouds.WithLocations("/etc/openstack/clouds.yaml"))
		if err != nil {
			log.Fatalf("Failed to parse the clouds.yaml: %v", err)
		}

		pc, err := openstack.AuthenticatedClient(ctx, ao)
		if err != nil {
			log.Fatalf("Failed to make providerClient with AuthenticatedClient: %v", err)
		}

		_identityClient, err := openstack.NewIdentityV3(pc, eo)
		if err != nil {
			log.Fatalf("Failed to make _networkClient with NewNetworkV2: %v", err)
		}

		identityClient = _identityClient
	}
}
