## 2024-05-24 - Symlink Protection in Template Replacement
**Vulnerability:** Symlink attack during project initialization. A cloned template could contain a symlink `README.md -> /etc/passwd`. The tool would then overwrite the system file when performing variable replacement.
**Learning:** Automatic variable replacement in files after cloning an untrusted repository is a high-risk operation if symlinks are followed.
**Prevention:** Always use `os.Lstat` to check for symlinks before reading from or writing to files that were recently created from an external source.

## 2025-05-14 - Enhanced Symlink Protection in File Operations
**Vulnerability:** Information disclosure and arbitrary file overwrite during template copying. `copyDir` and `copyFile` followed symlinks at both source and destination.
**Learning:** Hardening individual file operations is not enough; all paths in a recursive copy or file creation must be validated to prevent jumping out of the intended directory via symlinks.
**Prevention:** Skip symlinks at the source to prevent reading unauthorized files, and reject symlinks at the destination to prevent overwriting existing files outside the target directory. Use os.Lstat for all existence/mode checks.
