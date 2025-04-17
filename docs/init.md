# Init Command

The `init` command leverages the custom templates configured in the ~/.config/Glyph/repositories.toml file. For each defined template, a corresponding subcommand is dynamically created under init, allowing you to quickly clone and initialize new projects based on your custom scaffolding.

**Usage:**

```sh
glyph init <template_name> <destination_directory>
```

**Arguments:**

<template_name>: The name of the custom template you want to use (as defined with the create command).
<destination_directory>: The path where you want to clone the template repository and initialize your new project.

**Examples:**

Initializing a new Node.js project using the node-basic template:

```sh
glyph init node-basic my-new-node-app
```

Initializing a React project using the react-ui template:

```sh
glyph init react-ui frontend
```
