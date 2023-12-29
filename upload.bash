#!/bin/bash

# Check if the directory or file path is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <path_to_plugin_directory_or_file>"
    exit 1
fi

PLUGIN_PATH="$1"

# URL of the plugin upload endpoint
URL="http://localhost:1770/api/plugins/upload?install=true"

# Check if the path exists
if [ ! -e "$PLUGIN_PATH" ]; then
    echo "Error: Path not found!"
    exit 1
fi

cd "$PLUGIN_DIR"
# Build the plugin
echo "Building the plugin..."
go build -buildmode=plugin -o nodes.so *go
if [ $? -ne 0 ]; then
    echo "Error: Failed to build the plugin."
    exit 1
fi


# Create a ZIP file containing only the built .so file
ZIP_FILE="nodes-v1.0.0.zip"
zip "$ZIP_FILE" "nodes.so"
if [ ! -f "$ZIP_FILE" ]; then
    echo "Error: Failed to create ZIP file."
    exit 1
fi

# Check if the zip operation was successful
if [ ! -f "$ZIP_FILE" ]; then
    echo "Error: Failed to create ZIP file."
    exit 1
fi

# Using curl to upload the zip file
echo "Uploading $ZIP_FILE..."
curl -X POST -F "file=@$ZIP_FILE" $URL

# Optional: Remove the ZIP file after upload
# rm "$ZIP_FILE"
