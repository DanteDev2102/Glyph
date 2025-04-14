#!/bin/sh

# CLI name (the globally installed executable)
CLI_EXECUTABLE="glyph"

# Directory where the global command was installed
GLOBAL_BIN_DIR="$HOME/.local/bin"

echo "Uninstalling the global command '${CLI_EXECUTABLE}'..."

# Remove the executable from the global bin directory
if [ -f "$GLOBAL_BIN_DIR/${CLI_EXECUTABLE}" ]; then
  rm "$GLOBAL_BIN_DIR/${CLI_EXECUTABLE}"
  echo "Executable '${CLI_EXECUTABLE}' removed from ${GLOBAL_BIN_DIR}."
else
  echo "Executable '${CLI_EXECUTABLE}' not found in ${GLOBAL_BIN_DIR}."
fi

echo "Removing '${GLOBAL_BIN_DIR}' from your PATH in the shell configuration files..."

# Detect the shell and remove the line from the configuration file
SHELL_NAME=$(basename "$SHELL")
CONFIG_FILE=""
PATH_LINE="export PATH=\"\$PATH:${GLOBAL_BIN_DIR}\""

case "$SHELL_NAME" in
  bash)
    CONFIG_FILE="$HOME/.bashrc"
    ;;
  zsh)
    CONFIG_FILE="$HOME/.zshrc"
    ;;
  *)
    echo "Could not automatically detect your shell."
    echo "Please manually remove the line that adds '${GLOBAL_BIN_DIR}' to your PATH in your shell's configuration."
    exit 0
    ;;
esac

# Remove the PATH line from the configuration file
if sed -i '/export PATH="\$PATH:'"$GLOBAL_BIN_DIR"'"/d' "$CONFIG_FILE"; then
  echo "The line containing '${GLOBAL_BIN_DIR}' in your PATH was removed from ${CONFIG_FILE}."
else
  echo "The line containing '${GLOBAL_BIN_DIR}' in your PATH was not found in ${CONFIG_FILE}."
fi

rm -rf "$HOME/.config/Glyph"

echo "You might need to restart your terminal or run 'source $CONFIG_FILE' for the changes to take effect."
echo "Uninstallation complete!"

exit 0