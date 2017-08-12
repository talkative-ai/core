package models

// User model for the Workbench User, authenticated by Google
// TODO: This will change as we allow sign-ups to AUM via other modes than just Google
type User struct {
	Sub     string
	Email   string
	Name    string
	Picture string
}
