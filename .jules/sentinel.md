## 2024-05-24 - Symlink Protection in Template Replacement
**Vulnerability:** Symlink attack during project initialization. A cloned template could contain a symlink `README.md -> /etc/passwd`. The tool would then overwrite the system file when performing variable replacement.
**Learning:** Automatic variable replacement in files after cloning an untrusted repository is a high-risk operation if symlinks are followed.
**Prevention:** Always use `os.Lstat` to check for symlinks before reading from or writing to files that were recently created from an external source.

## 2025-05-14 - Enhanced Symlink Protection in File Operations
**Vulnerability:** Information disclosure and arbitrary file overwrite during template copying. `copyDir` and `copyFile` followed symlinks at both source and destination.
**Learning:** Hardening individual file operations is not enough; all paths in a recursive copy or file creation must be validated to prevent jumping out of the intended directory via symlinks.
**Prevention:** Skip symlinks at the source to prevent reading unauthorized files, and reject symlinks at the destination to prevent overwriting existing files outside the target directory. Use os.Lstat for all existence/mode checks.

## 2026-02-17 - Symlink Protection in Configuration Parser
**Vulnerability:** Symlink attack on the global configuration file. An attacker could pre-create the configuration file as a symlink to a sensitive file, causing the tool to overwrite it when saving templates.
**Learning:** Hardening the CLI's project initialization is insufficient if the configuration persistence layer remains vulnerable to the same class of attacks.
**Prevention:** Always use `os.Lstat` to verify that a file is not a symbolic link before performing read or write operations on shared or predictable configuration paths.

## 2026-03-20 - Credential Validation and Information Disclosure
**Vulnerability:** Lack of input validation in authentication prompts and information disclosure via stack traces. Users could provide empty credentials or invalid SSH key paths, and configuration errors would trigger a panic, leaking internal details.
**Learning:** Authentication flows should always validate inputs to fail early and securely. Global initialization in CLI tools must handle errors gracefully instead of panicking to avoid exposing internal state to users.
**Prevention:** Use validation functions for all user inputs in interactive prompts. Replace `panic` with controlled exits and sanitized error messages in the application's entry points.

## 2026-04-12 - Secure Configuration File Permissions
**Vulnerability:** Configuration file created with world-readable permissions (0644).
**Learning:** Default file creation permissions in many languages/environments are overly permissive for configuration files that may contain sensitive user information or preferences.
**Prevention:** Explicitly set restricted permissions (e.g., 0600) when creating or writing to configuration files to ensure only the owner can access them.

## 2026-05-15 - Template Name Validation
**Vulnerability:** Lack of validation for template names could lead to configuration corruption (e.g., overwriting the `[config]` section) or CLI command shadowing.
**Learning:** User-provided keys in a configuration file that also dictate CLI subcommands must be strictly validated against a whitelist of characters and a blacklist of reserved words.
**Prevention:** Enforce strict naming conventions (alphanumeric, hyphens, underscores) and reject reserved keywords used by the application's configuration or CLI framework.

## 2026-06-18 - Preservation of File Permissions during Template Initialization
**Vulnerability:** Insecure file permissions during project scaffolding. When cloning or copying templates, file permissions (including the executable bit and restricted access modes like 0600) were lost, defaulting to system-wide permissive modes (e.g., 0644).
**Learning:** Default file creation and copying operations in many standard libraries do not automatically replicate source permissions, which can expose sensitive configuration files or break scripts in the generated project.
**Prevention:** Always retrieve the source file's `os.FileMode` using `os.Lstat` and explicitly apply it to the destination using `os.Chmod` or by passing `info.Mode().Perm()` to file creation functions like `os.OpenFile`.
