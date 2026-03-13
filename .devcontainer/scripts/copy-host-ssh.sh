#!/usr/bin/env bash
set -euo pipefail

SOURCE_DIR="${1:-/tmp/host-ssh}"
TARGET_DIR="${HOME}/.ssh"

if [[ ! -d "${SOURCE_DIR}" ]]; then
  echo "[copy-host-ssh] Source '${SOURCE_DIR}' was not found. Skipping copy."
  exit 0
fi

mkdir -p "${TARGET_DIR}"
chmod 700 "${TARGET_DIR}"

shopt -s dotglob nullglob
for item in "${SOURCE_DIR}"/*; do
  [[ -e "${item}" ]] || continue
  cp -R "${item}" "${TARGET_DIR}/"
done

# OpenSSH expects strict permissions for private keys and config files.
find "${TARGET_DIR}" -type d -exec chmod 700 {} +
find "${TARGET_DIR}" -type f -name "*.pub" -exec chmod 644 {} +
find "${TARGET_DIR}" -type f -name "known_hosts" -exec chmod 644 {} +
find "${TARGET_DIR}" -type f ! -name "*.pub" ! -name "known_hosts" -exec chmod 600 {} +

echo "[copy-host-ssh] Host SSH files copied to '${TARGET_DIR}'."
