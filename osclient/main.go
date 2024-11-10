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
)

func init() {
	log.Print("In the init before bool")
	if OBIIIConfig.OpenstackEnabled {
		log.Print("In the init after bool")
		ctx := context.Background()
		ao, eo, tlsConfig, err := clouds.Parse()
		if err != nil {
			log.Fatalf("Failed to parse the clouds.yaml: %v", err)
		}

		providerClient, err := config.NewProviderClient(ctx, ao, config.WithTLSConfig(tlsConfig))
		if err != nil {
			log.Fatalf("Failed to make providerClient with NewProviderClient: %v", err)
		}

		_identityClient, err := openstack.NewIdentityV3(providerClient, eo)
		if err != nil {
			log.Fatalf("Failed to make _networkClient with NewNetworkV2: %v", err)
		}

		identityClient = _identityClient
	}
}
