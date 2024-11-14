package osclient

import (
	"context"
	"log"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
)

// Returns the username portion of the email
func extractUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return strings.ToLower(parts[0])
	}
	return ""
}

func CheckUserExists(email string) (bool, error) {
	username := extractUsername(email)

	ctx := context.Background()
	_, err := users.Get(ctx, identityClient, username).Extract()
	if err != nil {
		log.Printf("%T", err)
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			// User does not exist
			return false, nil
		}
		// A different error happened
		return false, err
	}
	// User exists
	return true, nil
}

// func Createuser(email string) {
// 	username := extractUsername(email)

// 	ctx := context.Background()

// }

// Returns a randomly generated password
// func generateTempPassworod() string {
// 	// Generate random bytes
// 	b := make([]byte, 16) // 16 bytes = 128 bits
// 	_, err := io.ReadFull(rand.Reader, b)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Compute MD5 hash
// 	hash := md5.New()
// 	hash.Write(b)
// 	md5sum := hex.EncodeToString(hash.Sum(nil))

// 	// Return first 24 characters of the hash
// 	if len(md5sum) > 24 {
// 		return md5sum[:24]
// 	}
// 	return md5sum
// }
