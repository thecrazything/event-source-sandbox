package login


func LoginUser(username string, password string) (string, error) {
	// Simulate user login logic
	println("Logging in user:", username)
	println("Password:", password)

	// Here you would typically check the user's credentials against a database
	// For this example, we'll just print the details and return a mock token
	token := "mock-token-for-" + username

	return token, nil
}