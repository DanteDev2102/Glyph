package config

const (
	CreateUse         string = "create --repo <repository-url> --name <template-name> --summary <short-description> --description <long-description>"
	CreateSummary     string = "Create new template command"
	CreateDescription string = `
Creates a new, user-defined command within Glyph that leverages remote Git repositories as scaffolding templates for project initialization.

This command allows you to define a new subcommand that, when invoked, will clone a specified Git repository. This cloned repository then serves as the foundation for creating new projects with a consistent and pre-configured structure.

You will need to provide the following information when creating a new command:

--name -n: The name of the new Glyph subcommand you wish to create. This will be used with the init command (e.g., glyph init <your-new-command>).
--repo -r: The URL of the Git repository containing the scaffolding template.
--branch -b (optional): The specific branch of the repository to clone. If not specified, the default branch will be used.
--tag -t (optional): A specific tag within the repository to checkout. If both --branch and --tag are provided, --tag will take precedence.
--summary -s (optional): A short, one-line description of the template.
--description -d (optional): A more detailed explanation of what this template provides.

The configuration for this new command will be stored in the ~/.config/Glyph/repositories.toml file, making it available for subsequent use with the init command.

Example:
	glyph create --name fiber-template --repo git@github.com/MyUser/fiber-template.git
	`

	InitUse         string = "init <template-name> <destination-directory>"
	InitSummary     string = "Initializes a new project using a custom template."
	InitDescription string = `
This command clones a pre-configured Git repository into the specified
destination directory, based on a template you've previously defined
using the 'create' command.

Arguments:
	<template_name> 		The name of the custom template to use. Glyph
    				will look up the repository details for this name
                  			in its configuration.

	<destination_directory> The path where you want to clone the template
	                      	repository and initialize your new project. This
	                      	directory will be created if it doesn't exist.

Examples:
	glyph init node-basic my-new-app
	glyph init react-ui ~/Projects/frontend

Notes:
	Before using 'init', you must first define your custom templates using
	the 'create' command. Glyph reads the template configurations from
	~/.config/Glyph/repositories.toml.
	`

	RmUse                = "rm <template-name>"
	RmSummary     string = "The delete command reads the Glyph configuration file "
	RmDescription string = `
Deletes a custom template configuration from Glyph.

This command removes the specified template entry from the
~/.config/Glyph/repositories.toml configuration file. Once deleted,
the template will no longer be available for use with the 'init' command.

<template-name> is the name of the custom template you wish to remove.

Examples:
	glyph delete my-old-template
	`

	SetUse         string = "set <old-template-name>"
	SetSummary     string = "Modify an existing custom template configuration."
	SetDescription string = `
Modifies an existing custom template configuration.

This command allows you to change the settings of a previously created
template in the ~/.config/Glyph/repositories.toml file. You can update
properties like the repository URL, branch, tag, summary, description,
and even rename the template.

<old-template-name> is the current name of the template you want to update.

Flags:
	-b, --branch string      The new branch name for the template.
	-d, --description string The new detailed description for the template.
	-n, --name string        The new name for the template.
	-r, --repo string        The new URL of the Git repository for the template.
	-s, --summary string     The new short description for the template.
	-t, --tag string         The new tag name for the template.

Examples:
	glyph update existing-template --repo [https://new-repo.com/project.git](https://new-repo.com/project.git)
	glyph update my-old-name --name my-new-name --summary "Updated description"
	`
)
