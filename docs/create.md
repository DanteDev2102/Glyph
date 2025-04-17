# Create Command

The `create` command allows you to define new custom commands within Glyph, enabling the cloning of a remote repository to be used as a scaffolding template for project initialization. The configuration for these templates is stored in the `~/.config/Glyph/repositories.toml` file.

**Usage:**

```sh
glyph create --name <template_name> --repo <repository_url> [--branch <branch_name>] [--tag <tag_name>] [--summary "<short_description>"] [--description "<detailed_description>"]
```

**Options:**

- `--name`, `-n`: _(Required)_ A unique name for your template. This will be used as the subcommand with `glyph init`.
- `--repo`, `-r`: _(Required)_ The URL of the Git repository containing the scaffolding template.
- `--branch`, `-b`: _(Optional)_ The specific branch of the repository to clone. Defaults to the repository's default branch.
- `--tag`, `-t`: _(Optional)_ A specific tag within the repository to checkout. If both `--branch` and `--tag` are provided, `--tag` takes precedence.
- `--summary`, `-s`: _(Optional)_ A short, one-line description of the template.
- `--description`, `-d`: _(Optional)_ A more detailed explanation of what this template provides.

**Examples:**

1.  **Creating a basic Node.js template:**

```sh
glyph create --name node-basic --repo https://github.com/expressjs/express-generator --summary "Basic Express.js application" --description "A simple Express.js application structure."
```

2.  **Creating a specific branch template for a React project:**

```sh
glyph create --name react-feature --repo https://github.com/facebook/create-react-app --branch feature-ui --summary "React app with the feature-ui branch" --description "A React application initialized from the 'feature-ui' branch."
```

3. **Creating a template using a specific tag:**

```sh
glyph create --name vue-legacy --repo https://github.com/vuejs/vue-cli --tag v3.0.0 --summary "Vue CLI v3.0.0 template" --description "A Vue.js project initialized with Vue CLI version 3.0.0."
```
