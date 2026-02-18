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

## 2025-05-15 - Support for Encrypted SSH Keys
**Vulnerability:** Weak security practice due to lack of passphrase support. The tool previously only supported unencrypted SSH keys, which could encourage users to use less secure authentication methods.
**Learning:** Hardcoding an empty passphrase in SSH authentication logic significantly reduces the security posture of the tool by making encrypted keys unusable.
**Prevention:** Always provide an option for users to enter a passphrase when using SSH key authentication, and ensure it is collected securely using masked inputs.
