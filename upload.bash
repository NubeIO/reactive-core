#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <path_to_plugin_directory_or_file> <plugin_name>"
    exit 1
fi

PLUGIN_PATH="$1"
PLUGIN_NAME="$2.so"

# URL of the plugin upload endpoint
URL="http://localhost:1770/api/plugins/upload?install=true"

# Check if the path exists and is a directory
if [ ! -d "$PLUGIN_PATH" ]; then
    echo "Error: Directory path not found!"
    exit 1
fi

cd "$PLUGIN_PATH"

# Build the plugin
echo "Building the plugin..."
go build -ldflags "-s -w" -buildmode=plugin -o "$PLUGIN_NAME"  *go

if [ $? -ne 0 ]; then
    echo "Error: Failed to build the plugin."
    exit 1
fi

# Create a ZIP file containing only the built .so file
ZIP_FILE="${PLUGIN_NAME%.*}.zip"
zip "$ZIP_FILE" "$PLUGIN_NAME"
if [ ! -f "$ZIP_FILE" ]; then
    echo "Error: Failed to create ZIP file."
    exit 1
fi

# Using curl to upload the zip file
echo "Uploading $ZIP_FILE..."
curl -X POST -F "file=@$ZIP_FILE" $URL

rm "$ZIP_FILE"
rm "$PLUGIN_NAME"