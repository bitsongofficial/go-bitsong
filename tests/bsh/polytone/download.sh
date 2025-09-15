#!/bin/bash
BASE_URL="https://github.com/DA0-DA0/polytone/releases/download/v1.1.0/"
# Directory to store WASM files
POLYTONE_WASM_DIR="bin/"
mkdir -p "$POLYTONE_WASM_DIR"

POLYTONE_CONTRACTS=(
  "polytone_listener.wasm"
  "polytone_note.wasm"
  "polytone_proxy.wasm"
  "polytone_voice.wasm"
  "polytone_tester.wasm"
  )

# Download each WASM file
for CONTRACT in "${POLYTONE_CONTRACTS[@]}"; do
  FILE_URL="$BASE_URL$CONTRACT"
  FILE_PATH="$POLYTONE_WASM_DIR/$CONTRACT"
  
  # Construct the curl command
  CURL_CMD="curl -v -L -o '$FILE_PATH' '$FILE_URL'"
  
  # Attempt the download
  eval "$CURL_CMD"
  
  # Check the file size after download
  FILE_SIZE=$(stat -c%s "$FILE_PATH" 2>/dev/null)
  
  if [ $? -eq 0 ] && [ $FILE_SIZE -gt 0 ]; then
    echo "$CONTRACT downloaded successfully (Size: $FILE_SIZE bytes)."
  else
    echo "Failed to download $CONTRACT or the file is empty. **Please try copying and running the above curl command manually to troubleshoot.**"
  fi
  echo
done