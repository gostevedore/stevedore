package credentials

// UserPasswordAuth containes user password auth for docker registry
type UserPasswordAuth struct {
	DEPRECATEDUsername string `json:"docker_login_username"`
	DEPRECATEDPassword string `json:"docker_login_password"`
	Username           string `json:"username"`
	Password           string `json:"password"`
}
