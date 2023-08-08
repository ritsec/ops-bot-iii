package helpers

// AtRole returns a string that mentions a role
func AtRole(role string) string {
	return "<@&" + role + ">"
}

// AtUser returns a string that mentions a user
func AtUser(user string) string {
	return "<@" + user + ">"
}
