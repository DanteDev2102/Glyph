#!/bin/sh

# CLI name (the executable is already in the root)
CLI_EXECUTABLE="glyph"

# Directory where the global command will be installed
GLOBAL_BIN_DIR="$HOME/.local/bin"

CONFIG_DIR="$HOME/.config/Glyph/"

echo "Installing '${CLI_EXECUTABLE}' as a global command..."

# Create the global bin directory if it doesn't exist
mkdir -p "$GLOBAL_BIN_DIR"

mkdir -p "$CONFIG_DIR"

cp ./config/repositories.toml "$CONFIG_DIR"
chmod +r "$HOME/.config/Glyph/repositories.toml"


# Check if the executable exists in the root
if [ -f "./${CLI_EXECUTABLE}" ]; then
  # Copy the executable to the global bin directory
  cp "./${CLI_EXECUTABLE}" "$GLOBAL_BIN_DIR/${CLI_EXECUTABLE}"

  # Give execute permissions to the global executable
  chmod +x "$GLOBAL_BIN_DIR/${CLI_EXECUTABLE}"

  echo "${CLI_EXECUTABLE} installed globally in ${GLOBAL_BIN_DIR}"
  echo "Automatically adding '${GLOBAL_BIN_DIR}' to your PATH..."

  # Detect the shell and add the line to the configuration file
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
      echo "Please manually add '${GLOBAL_BIN_DIR}' to your shell's PATH variable."
      exit 0
      ;;
  esac

  # Check if the line already exists in the configuration file
  if grep -qF "$PATH_LINE" "$CONFIG_FILE"; then
    echo "'${GLOBAL_BIN_DIR}' is already in your PATH in ${CONFIG_FILE}."
  else
    # Add the line to the end of the configuration file
    echo "$PATH_LINE" >> "$CONFIG_FILE"
    echo "'${GLOBAL_BIN_DIR}' has been added to your PATH in ${CONFIG_FILE}."
  fi

  echo "You might need to restart your terminal or run 'source $CONFIG_FILE' for the changes to take effect."
  echo "Done! You should be able to use the '${CLI_EXECUTABLE}' command in your terminal."

else
  echo "Error: The executable '${CLI_EXECUTABLE}' was not found in the project root."
  echo "Make sure to run this script after building the project."
  exit 1
fi

exit 0