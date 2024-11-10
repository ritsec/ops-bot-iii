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
	if OBIIIConfig.OpenstackEnabled {
		ctx := context.Background()
		ao, eo, tlsConfig, err := clouds.Parse(clouds.WithCloudName("openstack"), clouds.WithLocations("/etc/openstack/clouds.yaml"))
		if err != nil {
			log.Fatalf("Failed to parse the clouds.yaml: %v", err)
		}

		tlsConfig.InsecureSkipVerify = true

		providerClient, err := config.NewProviderClient(ctx, ao, config.WithTLSConfig(tlsConfig))
		if err != nil {
			log.Fatalf("Failed to make providerClient with NewProviderClient: %v\ntlsConfig: %#v\n", err, tlsConfig)
		}

		_identityClient, err := openstack.NewIdentityV3(providerClient, eo)
		if err != nil {
			log.Fatalf("Failed to make _networkClient with NewNetworkV2: %v", err)
		}

		identityClient = _identityClient
	}
}
