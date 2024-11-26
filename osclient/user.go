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

	networkOpts := network.UpdateOpts{
		FloatingIP:        gophercloud.IntToPointer(0),
		Network:           gophercloud.IntToPointer(10),
		Port:              gophercloud.IntToPointer(50),
		Router:            gophercloud.IntToPointer(1),
		Subnet:            gophercloud.IntToPointer(20),
		SecurityGroup:     gophercloud.IntToPointer(10),
		SecurityGroupRule: gophercloud.IntToPointer(-1),
	}

	projectID := project.ID
	_, err = network.Update(ctx, networkClient, projectID, networkOpts).Extract()
	if err != nil {
		return "", "", err
	}

	quotaOpts := compute.UpdateOpts{
		InjectedFileContentBytes: gophercloud.IntToPointer(10240),
		InjectedFilePathBytes:    gophercloud.IntToPointer(255),
		InjectedFiles:            gophercloud.IntToPointer(5),
		KeyPairs:                 gophercloud.IntToPointer(10),
		RAM:                      gophercloud.IntToPointer(51200),
		Cores:                    gophercloud.IntToPointer(20),
		Instances:                gophercloud.IntToPointer(10),
		ServerGroups:             gophercloud.IntToPointer(10),
		ServerGroupMembers:       gophercloud.IntToPointer(10),
	}

	_, err = compute.Update(ctx, computeClient, projectID, quotaOpts).Extract()
	if err != nil {
		return "", "", err
	}

	storageOpts := blockstorage.UpdateOpts{
		Volumes:   gophercloud.IntToPointer(10),
		Snapshots: gophercloud.IntToPointer(10),
		Gigabytes: gophercloud.IntToPointer(250),
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
