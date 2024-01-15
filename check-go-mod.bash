#!/bin/bash

# Colors
GREEN='\033[0;32m' # OK
ORANGE='\033[0;33m' # N/A
RED='\033[0;31m' # ERROR
NC='\033[0m' # No Color

# Check if at least two directories are provided
if [ "$#" -lt 2 ]; then
    echo "Please provide at least two directory paths to compare."
    exit 1
fi

# Initialize arrays to store dependencies and their versions
declare -A dependencies

# Iterate over each directory
for dir in "$@"; do
    # Check if go.mod file exists
    if [ -f "$dir/go.mod" ]; then
        echo "Checking $dir"
        # Extract and store dependencies and their versions
        while read -r line; do
            if [[ "$line" =~ ^[[:space:]]*require[[:space:]]*([^[:space:]]+)[[:space:]]*v([^[:space:]]+) ]]; then
                dependency="${BASH_REMATCH[1]}"
                version="${BASH_REMATCH[2]}"
                if [ -z "${dependencies["$dependency"]}" ]; then
                    dependencies["$dependency"]="$version ($dir)"
                else
                    dependencies["$dependency"]="${dependencies["$dependency"]} $version ($dir)"
                fi
            fi
        done < "$dir/go.mod"
    else
        echo -e "${RED}ERROR:${NC} $dir does not contain a go.mod file."
        exit 1
    fi
done

# Debug: Print dependencies
echo "Debug: Dependencies"
for dependency in "${!dependencies[@]}"; do
    echo "$dependency: ${dependencies["$dependency"]}"
done

# Compare dependencies and versions
for dependency in "${!dependencies[@]}"; do
    versions="${dependencies["$dependency"]}"
    if [ -z "$versions" ]; then
        echo -e "${ORANGE}NA:${NC} $dependency"
    else
        echo -e "${GREEN}OK:${NC} $dependency $versions"
    fi
done
