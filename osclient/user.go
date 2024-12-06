package osclient

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	blockstorage "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/quotasets"
	compute "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/quotasets"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	network "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
	OBIIIConfig "github.com/ritsec/ops-bot-iii/config"
)

// networkOpts is a struct for the quotas to be applied in network service
type networkOpts struct {
	FloatingIP        int
	Network           int
	Port              int
	Router            int
	Subnet            int
	SecurityGroup     int
	SecurityGroupRule int
}

func newNetworkOpts() networkOpts {
	return networkOpts{
		// FloatingIP specifies the number of floating IPs the user can allocate.
		FloatingIP: 0,
		// Network specifies the number of networks the user can create.
		Network: 10,
		// Port specifies the number of ports the user can create.
		Port: 50,
		// Router specifies the number of routers the user can create.
		Router: 1,
		// Subnet specifies the number of subnets the user can create.
		Subnet: 20,
		// SecurityGroup specifies the number of security groups the user can create.
		SecurityGroup: 10,
		// SecurityGroupRule specifies the number of security group rules the user can create.
		SecurityGroupRule: -1,
	}
}

// QuotaOpts defines the quota options for a compute resource in OpenStack.
type quotaOpts struct {
	InjectedFileContentBytes int
	InjectedFilePathBytes    int
	InjectedFiles            int
	KeyPairs                 int
	RAM                      int
	Cores                    int
	Instances                int
	ServerGroups             int
	ServerGroupMembers       int
}

func newQuotaOpts() quotaOpts {
	return quotaOpts{
		// InjectedFileContentBytes specifies the number of bytes allowed for injected file content.
		InjectedFileContentBytes: 10240,
		// InjectedFilePathBytes specifies the number of bytes allowed for the path of injected files.
		InjectedFilePathBytes: 255,
		// InjectedFiles specifies the number of injected files allowed.
		InjectedFiles: 5,
		// KeyPairs specifies the number of key pairs allowed.
		KeyPairs: 10,
		// RAM specifies the amount of RAM (in megabytes) allowed.
		RAM: 51200,
		// Cores specifies the number of CPU cores allowed.
		Cores: 20,
		// Instances specifies the number of instances allowed.
		Instances: 10,
		// ServerGroups specifies the number of server groups allowed.
		ServerGroups: 10,
		// ServerGroupMembers specifies the number of members allowed per server group.
		ServerGroupMembers: 10,
	}
}

// StorageOpts defines the quota options for block storage resources in OpenStack.
type StorageOpts struct {
	Volumes   int
	Snapshots int
	Gigabytes int
}

func newStorageOpts() StorageOpts {
	return StorageOpts{
		// Volumes specifies the number of volumes allowed.
		Volumes: 10,
		// Snapshots specifies the number of snapshots allowed.
		Snapshots: 10,
		// Gigabytes specifies the amount of storage (in gigabytes) allowed.
		Gigabytes: 250,
	}
}

// Returns the username portion of the email
func extractUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return strings.ToLower(parts[0])
	}
	return ""
}

// Checks to see if the user exists on openstack
func CheckUserExists(email string) (bool, error) {
	username := extractUsername(email)

	ctx := context.Background()

	listOpts := users.ListOpts{
		DomainID: "default",
	}

	allPages, err := users.List(identityClient, listOpts).AllPages(ctx)
	if err != nil {
		return false, err
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return false, err
	}

	for _, user := range allUsers {
		if user.Name == username {
			return true, nil
		}
	}

	return false, nil
}

// Returns a password that should be used temporary to create the account or resetting the account
func tempPass() string {
	randomNumber := rand.Int()
	hash := md5.Sum([]byte(fmt.Sprint(randomNumber)))
	return fmt.Sprintf("%x", hash)[:24]
}

// Creats the account on openstack and returns username and password for the user to get
func Create(email string) (string, string, error) {
	ctx := context.Background()
	username := extractUsername(email)

	// Create Project
	createOpts := projects.CreateOpts{
		DomainID: "default",
		Enabled:  gophercloud.Enabled,
		Name:     username,
	}

	project, err := projects.Create(ctx, identityClient, createOpts).Extract()
	if err != nil {
		return "", "", err
	}

	defaultNetworkOpts := newNetworkOpts()
	networkOpts := network.UpdateOpts{
		FloatingIP:        gophercloud.IntToPointer(defaultNetworkOpts.FloatingIP),
		Network:           gophercloud.IntToPointer(defaultNetworkOpts.Network),
		Port:              gophercloud.IntToPointer(defaultNetworkOpts.Port),
		Router:            gophercloud.IntToPointer(defaultNetworkOpts.Router),
		Subnet:            gophercloud.IntToPointer(defaultNetworkOpts.Subnet),
		SecurityGroup:     gophercloud.IntToPointer(defaultNetworkOpts.SecurityGroup),
		SecurityGroupRule: gophercloud.IntToPointer(defaultNetworkOpts.SecurityGroupRule),
	}

	projectID := project.ID
	_, err = network.Update(ctx, networkClient, projectID, networkOpts).Extract()
	if err != nil {
		return "", "", err
	}

	defaultQuotaOpts := newQuotaOpts()
	quotaOpts := compute.UpdateOpts{
		InjectedFileContentBytes: gophercloud.IntToPointer(defaultQuotaOpts.InjectedFileContentBytes),
		InjectedFilePathBytes:    gophercloud.IntToPointer(defaultQuotaOpts.InjectedFilePathBytes),
		InjectedFiles:            gophercloud.IntToPointer(defaultQuotaOpts.InjectedFiles),
		KeyPairs:                 gophercloud.IntToPointer(defaultQuotaOpts.KeyPairs),
		RAM:                      gophercloud.IntToPointer(defaultQuotaOpts.RAM),
		Cores:                    gophercloud.IntToPointer(defaultQuotaOpts.Cores),
		Instances:                gophercloud.IntToPointer(defaultQuotaOpts.Instances),
		ServerGroups:             gophercloud.IntToPointer(defaultQuotaOpts.ServerGroups),
		ServerGroupMembers:       gophercloud.IntToPointer(defaultQuotaOpts.ServerGroupMembers),
	}

	_, err = compute.Update(ctx, computeClient, projectID, quotaOpts).Extract()
	if err != nil {
		return "", "", err
	}

	defaultStorageOpts := newStorageOpts()
	storageOpts := blockstorage.UpdateOpts{
		Volumes:   gophercloud.IntToPointer(defaultStorageOpts.Volumes),
		Snapshots: gophercloud.IntToPointer(defaultStorageOpts.Snapshots),
		Gigabytes: gophercloud.IntToPointer(defaultStorageOpts.Gigabytes),
	}

	_, err = blockstorage.Update(ctx, storageClient, projectID, storageOpts).Extract()
	if err != nil {
		return "", "", err
	}

	password := tempPass()
	userOpts := users.CreateOpts{
		Name:             username,
		DefaultProjectID: projectID,
		Description:      "",
		DomainID:         "default",
		Enabled:          gophercloud.Enabled,
		Password:         password,
		Extra: map[string]any{
			"email": email,
		},
	}

	user, err := users.Create(ctx, identityClient, userOpts).Extract()
	if err != nil {
		return "", "", err
	}

	userID := user.ID
	err = roles.Assign(ctx, identityClient, OBIIIConfig.Openstack.MemberID, roles.AssignOpts{
		UserID:    userID,
		ProjectID: projectID,
	}).ExtractErr()
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

// Resets the account's password and returns the new temporary password for the user to login and change
func Reset(email string) (string, string, error) {
	ctx := context.Background()

	userID, err := GetUserID(email)
	if err != nil {
		return "", "", err
	}

	newPassword := tempPass()

	changePasswordOpts := users.UpdateOpts{
		Password: newPassword,
	}

	_, err = users.Update(ctx, identityClient, userID, changePasswordOpts).Extract()
	if err != nil {
		return "", "", err
	}

	return extractUsername(email), newPassword, nil
}

func GetUserID(email string) (string, error) {
	ctx := context.Background()
	username := extractUsername(email)

	listOpts := users.ListOpts{
		DomainID: "default",
	}

	allPages, err := users.List(identityClient, listOpts).AllPages(ctx)
	if err != nil {
		return "", err
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return "", err
	}

	for _, user := range allUsers {
		if user.Name == username {
			return user.ID, nil
		}
	}

	return "", fmt.Errorf("user id not found")
}
