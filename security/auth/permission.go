package auth

func (c *Claims) HasRole(role string) bool {
	return c.Role == role
}

func (c *Claims) HasPerm(p string) bool {
	for _, perm := range c.Perms {
		if perm == p {
			return true
		}
	}
	return false
}
