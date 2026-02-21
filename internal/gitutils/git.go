package gitutils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/manifoldco/promptui"
)

// CloneRepo clones a repository using go-git.
func CloneRepo(url, path, branch, tag string) error {
	cloneOptions := &git.CloneOptions{
		URL:          url,
		Depth:        1,
		SingleBranch: true,
	}

	if tag != "" {
		cloneOptions.ReferenceName = plumbing.NewTagReferenceName(tag)
	} else if branch != "" {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}

	_, err := git.PlainClone(path, false, cloneOptions)
	if err == transport.ErrAuthenticationRequired || err == transport.ErrAuthorizationFailed {
		auth, authErr := RequestAuth()
		if authErr != nil {
			return authErr
		}
		cloneOptions.Auth = auth
		_, err = git.PlainClone(path, false, cloneOptions)
	}

	return err
}

// ValidateRepo checks if a repository is accessible.
func ValidateRepo(url string) (transport.AuthMethod, error) {
	// Try listing remotes without auth first
	_, err := git.NewRemote(nil, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	}).List(&git.ListOptions{})

	if err == nil {
		return nil, nil
	}

	if err == transport.ErrAuthenticationRequired || err == transport.ErrAuthorizationFailed {
		fmt.Println("Repository requires authentication to validate access.")
		return RequestAuth()
	}

	return nil, err
}

// validateNotEmpty checks if the input is not empty or just whitespace.
func validateNotEmpty(input string) error {
	if len(strings.TrimSpace(input)) == 0 {
		return errors.New("input cannot be empty")
	}
	return nil
}

// RequestAuth prompts the user for authentication credentials.
func RequestAuth() (transport.AuthMethod, error) {
	prompt := promptui.Select{
		Label: "Select authentication method",
		Items: []string{"SSH Key", "Username and Token/Password"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	if result == "SSH Key" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not get home directory: %v", err)
		}
		prompt := promptui.Prompt{
			Label:   "Path to SSH Key",
			Default: home + "/.ssh/id_rsa",
			Validate: func(input string) error {
				if _, err := os.Stat(input); err != nil {
					return fmt.Errorf("file does not exist: %v", err)
				}
				return nil
			},
		}
		keyPath, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		promptPassphrase := promptui.Prompt{
			Label: "SSH Key Passphrase (leave empty if none)",
			Mask:  '*',
		}
		passphrase, err := promptPassphrase.Run()
		if err != nil {
			return nil, err
		}

		publicKeys, err := ssh.NewPublicKeysFromFile("git", keyPath, passphrase)
		if err != nil {
			return nil, err
		}
		return publicKeys, nil
	} else {
		promptUser := promptui.Prompt{
			Label:    "Username",
			Validate: validateNotEmpty,
		}
		username, err := promptUser.Run()
		if err != nil {
			return nil, err
		}

		promptPass := promptui.Prompt{
			Label:    "Token/Password",
			Mask:     '*',
			Validate: validateNotEmpty,
		}
		password, err := promptPass.Run()
		if err != nil {
			return nil, err
		}

		return &http.BasicAuth{
			Username: username,
			Password: password,
		}, nil
	}
}
