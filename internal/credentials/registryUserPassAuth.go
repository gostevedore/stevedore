package credentials

// RegistryUserPassAuth containes user password auth for docker registry
type RegistryUserPassAuth struct {
	Username string `json:"docker_login_username"`
	Password string `json:"docker_login_password"`
}
