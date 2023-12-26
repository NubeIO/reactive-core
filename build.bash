#!/bin/bash

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

ORIGINAL_DIR=$(pwd)


# Function to print messages in green
print_green() {
    echo -e "${GREEN}$1${NC}"
}

# Function to print messages in orange/yellow
print_orange() {
    echo -e "${YELLOW}$1${NC}"
}

# Function to print messages in red
print_red() {
    echo -e "${RED}$1${NC}"
}

to_camel() {
    local text="$1"
    local camelCaseText

    # Replace spaces with underscores, then use 'tr' to uppercase the first letter of each word
    camelCaseText=$(echo "$text" | tr ' ' '_' | sed -r 's/(^|_)([a-z])/\U\2/g')

    # Remove underscores
    camelCaseText=$(echo "$camelCaseText" | tr -d '_')

    echo "$camelCaseText"
}



lowercase_first() {
    local text="$1"
    # Extract the first character and the rest of the string
    local firstChar="${text:0:1}"
    local restOfString="${text:1}"
    # Convert the first character to lowercase and concatenate with the rest of the string
    echo "${firstChar,,}$restOfString"
}


cleanup() {
    rm -f "$ORIGINAL_DIR/stderr.txt"
    print_green "stop and clean up delete: $ORIGINAL_DIR/stderr.txt"
}

trap cleanup EXIT


check_go_version() {
    required_version="1.21"
    current_version=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
    print_green "golang version: $current_version"

    if [ -z "$current_version" ]; then
        print_orange "Go is not installed. Please install Go version $required_version or newer."
    fi

    if [ "$current_version" != "$required_version" ]; then
        print_orange "Go version is $current_version Please upgrade to Go version $required_version or newer."
    fi
}



# Function to display usage help
usage() {
    echo "Usage:"
    echo "  $0 new --project-name PROJECT_NAME --github-account GITHUB_ACCOUNT --plugin-name PLUGIN_NAME"
    echo "  $0 build --build-path BUILD_PATH --plugin-name PLUGIN_NAME --output-path OUTPUT_PATH"
    echo "  $0 setup --build-path BUILD_PATH"
    echo
    echo "Options for new:"
    echo "  --project-name     Name of the project"
    echo "  --github-account   GitHub account name"
    echo "  --plugin-name      Plugin name"
    echo
    echo "Options for build:"
    echo "  --build-path       Path to build the project"
    echo "  --plugin-name      Plugin name"
    echo "  --output-path      Path to send the output"
    echo
    echo "Options for setup:"
    echo "Compare the go.mod file in BUILD_PATH with a reference list of dependencies."
    echo "  --build-path       Path to build the project."
    echo "Help:"
    echo "  -h, --help         Show this help message"
    echo
}


# Function to compare go.mod dependencies
compare_dependencies() {
    local build_path="$1"
    local go_mod_file="$build_path/go.mod"

    # Check if go.mod exists
    if [ ! -f "$go_mod_file" ]; then
        echo -e "${RED}Error: go.mod not found in $build_path${NC}"
        return 1
    fi

    print_green "$go_mod_file"

    # Reference dependencies
    declare -A ref_deps=(
        ["github.com/NubeIO/lib-schema"]="v0.2.18"
        ["github.com/NubeIO/reactive"]="v0.0.1"
        ["github.com/gin-gonic/gin"]="v1.9.1"
        ["github.com/google/uuid"]="v1.5.0"
        ["github.com/mitchellh/mapstructure"]="v1.5.0"
        ["github.com/sirupsen/logrus"]="v1.9.3"
        ["github.com/spf13/cobra"]="v1.8.0"
        ["github.com/spf13/viper"]="v1.18.2"
        ["golang.org/x/sync"]="v0.5.0"
        ["github.com/NubeIO/schema"]="v0.0.1" # indirect
        ["github.com/bytedance/sonic"]="v1.9.1" # indirect
    )

    # Compare with the actual go.mod
    local index=1
    for dep in "${!ref_deps[@]}"; do
        local version=${ref_deps[$dep]}
        local current_version
        current_version=$(grep "$dep" "$go_mod_file" | awk '{print $2}')

        if [[ -z "$current_version" ]]; then
            echo -e "$index. $dep | ${YELLOW}missing${NC}"
        elif [[ "$current_version" != "$version" ]]; then
            echo -e "$index. $dep | ${RED}different version \"current: $current_version needed: $version\"${NC}"
        else
            echo -e "$index. $dep | ${GREEN}ok${NC}"
        fi
        ((index++))
    done
}

# Function to set up a new project
setup_project() {
    local project_name=""
    local github_account=""
    local plugin_name=""

    while [[ $# -gt 0 ]]; do
        case "$1" in
            --project-name)
                project_name="$2"
                shift 2
                ;;
            --github-account)
                github_account="$2"
                shift 2
                ;;
            --plugin-name)
                plugin_name="$2"
                shift 2
                ;;
            *)
                print_red "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done

    if [[ -z "$project_name" || -z "$github_account" || -z "$plugin_name" ]]; then
        print_red "Missing arguments for new project"
        usage
        exit 1
    fi

    check_go_version
    export_name=$(to_camel "$plugin_name")
    strut_name=$(lowercase_first "$export_name")
    echo "$strut_name"
    mkdir -p "$project_name"
    cat <<EOF >"$project_name/go.mod"
module github.com/$github_account/$project_name

go 1.19

require (
    github.com/NubeIO/reactive v0.0.2
    github.com/NubeIO/schema v0.0.1
)
EOF

    cat <<EOF >"$project_name/main.go"
package main

type $strut_name struct {}

func (n *$strut_name) Start() {}

func (n *$strut_name) Enable() bool {
    return true
}

// exports
var $export_name $strut_name

EOF

    print_green "Project setup complete."
}

# Function to build and run the project
run_build() {
    local build_path=""
    local plugin_name=""
    local output_path=""

    while [[ $# -gt 0 ]]; do
        case "$1" in
            --build-path)
                build_path="$2"
                shift 2
                ;;
            --plugin-name)
                plugin_name="$2"
                shift 2
                ;;
            --output-path)
                output_path="$2"
                shift 2
                ;;
            *)
                echo "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done

    if [[ -z "$build_path" || -z "$plugin_name" || -z "$output_path" ]]; then
        echo "Missing arguments for build"
        usage
        exit 1
    fi
    check_go_version

    to_camel "$plugin_name"

    cd "$build_path" || exit
    local build_command="go build -buildmode=plugin -o \"$output_path/plugins/$plugin_name.so\" *.go"

    # Print the build command
    print_green "$build_command"

    # Execute the build command
    eval "$build_command" 2>stderr.txt

    if [ $? -ne 0 ]; then
        print_red "Error during build: $(<stderr.txt)"
        exit 1
    fi
    cd "$output_path" || exit
    go run main.go server
}

# Check for help option
if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    usage
    exit 0
fi

# Main logic
case $1 in
    new)
        shift
        setup_project "$@"
        ;;
    build)
        shift
        run_build "$@"
        ;;
     setup)
         shift
         if [ "$#" -ne 2 ] || [ "$1" != "--build-path" ]; then
             echo "Usage: $0 setup --build-path BUILD_PATH"
             exit 1
         fi
         compare_dependencies "$2"
         ;;
     *)
         usage
         exit 1
         ;;
esac
