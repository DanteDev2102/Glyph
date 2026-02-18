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
		home, _ := os.UserHomeDir()
		prompt := promptui.Prompt{
			Label:   "Path to SSH Key",
			Default: home + "/.ssh/id_rsa",
			Validate: func(input string) error {
				if _, err := os.Stat(input); os.IsNotExist(err) {
					return errors.New("ssh key file does not exist")
				}
				return nil
			},
		}
		keyPath, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		promptPass := promptui.Prompt{
			Label: "SSH Key Passphrase (optional)",
			Mask:  '*',
		}
		passphrase, err := promptPass.Run()
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
			Label: "Username",
			Validate: func(input string) error {
				if len(strings.TrimSpace(input)) == 0 {
					return errors.New("username cannot be empty")
				}
				return nil
			},
		}
		username, err := promptUser.Run()
		if err != nil {
			return nil, err
		}

		promptPass := promptui.Prompt{
			Label: "Token/Password",
			Mask:  '*',
			Validate: func(input string) error {
				if len(strings.TrimSpace(input)) == 0 {
					return errors.New("token/password cannot be empty")
				}
				return nil
			},
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
