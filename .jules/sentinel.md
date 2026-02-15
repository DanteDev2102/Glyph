## 2024-05-24 - Symlink Protection in Template Replacement
**Vulnerability:** Symlink attack during project initialization. A cloned template could contain a symlink `README.md -> /etc/passwd`. The tool would then overwrite the system file when performing variable replacement.
**Learning:** Automatic variable replacement in files after cloning an untrusted repository is a high-risk operation if symlinks are followed.
**Prevention:** Always use `os.Lstat` to check for symlinks before reading from or writing to files that were recently created from an external source.

## 2025-05-15 - Expanded Symlink Protection in Init Operations
**Vulnerability:** Symlink attacks during project initialization in copyDir, copyFile, and destination path handling. An attacker could use symlinks to leak sensitive files from the host or overwrite system files.
**Learning:** Hardening only the file replacement logic is insufficient if the initial copy or destination validation also follows symlinks.
**Prevention:** Consistently use os.Lstat to check all file/directory targets and sources involved in template extraction and initialization. Skip symlinks in sources and reject them in destinations.
