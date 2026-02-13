## 2024-05-24 - Symlink Protection in Template Replacement
**Vulnerability:** Symlink attack during project initialization. A cloned template could contain a symlink `README.md -> /etc/passwd`. The tool would then overwrite the system file when performing variable replacement.
**Learning:** Automatic variable replacement in files after cloning an untrusted repository is a high-risk operation if symlinks are followed.
**Prevention:** Always use `os.Lstat` to check for symlinks before reading from or writing to files that were recently created from an external source.

## 2025-01-30 - Comprehensive Symlink Protection
**Vulnerability:** Partial symlink protection only covered variable replacement, leaving file copy operations (`copyFile`, `copyDir`) vulnerable to overwriting arbitrary files via pre-existing symlinks at the destination.
**Learning:** Security fixes must be applied to all relevant code paths. Protecting only one stage of a multi-stage file operation (like project initialization) leaves other stages exposed to similar attacks.
**Prevention:** Audit all functions that perform file I/O (especially those using `os.Create` or `os.Open` for writing) and ensure they verify the destination isn't a symbolic link when targeting untrusted or recently created directories.
