# Set Command

The set command reads the ~/.config/Glyph/repositories.toml file to find and modify an existing custom template. You can change any of its configured keys, including the name.

**Usage:**

```sh
glyph set <old_template_name> [--name <new_template_name>] [--repo <new_repository_url>] [--branch <new_branch_name>] [--tag <new_tag_name>] [--summary "<new_short_description>"] [--description "<new_detailed_description>"]
```

**Options:**

<old_template_name>: (Required) The current name of the template you want to update.
--name, -n: (Optional) The new name for the template.
--repo, -r: (Optional) The new URL of the Git repository.
--branch, -b: (Optional) The new branch name.
--tag, -t: (Optional) The new tag name.
--summary, -s: (Optional) The new short description.
--description, -d: (Optional) The new detailed description.

**Example:**

To update the repository URL of the node-basic template:

```sh
glyph set node-basic --repo https://github.com/nodejs/node-starter
```

To rename the react-feature template to react-ui and update its description:

```sh
glyph update react-feature --name react-ui --description "React application with the latest UI features."
```
