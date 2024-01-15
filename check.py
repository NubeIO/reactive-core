import os

# Function to read and parse a go.mod file
def parse_go_mod(file_path):
    dependencies = {}
    with open(file_path, 'r') as f:
        lines = f.readlines()
        read_dependency_lines = False
        for line in lines:
            if read_dependency_lines and line.strip():
                parts = line.split()
                if len(parts) >= 2:
                    dependency = parts[0]
                    version = parts[1].strip()
                    dependencies[dependency] = version
            elif line.startswith('require'):
                read_dependency_lines = True
    return dependencies

# Function to compare two go.mod files
def compare_go_mod(go_mod_paths):
    version_mismatches = []

    for i in range(len(go_mod_paths)):
        for j in range(i+1, len(go_mod_paths)):
            go_mod_path1 = go_mod_paths[i]
            go_mod_path2 = go_mod_paths[j]

            dir_name1 = os.path.basename(os.path.dirname(go_mod_path1))
            dir_name2 = os.path.basename(os.path.dirname(go_mod_path2))

            dependencies1 = parse_go_mod(go_mod_path1)
            dependencies2 = parse_go_mod(go_mod_path2)

            # Find dependencies present in both files
            common_dependencies = set(dependencies1.keys()) & set(dependencies2.keys())

            # Compare versions for common dependencies
            for dependency in common_dependencies:
                version1 = dependencies1[dependency]
                version2 = dependencies2[dependency]
                if version1 != version2:
                    version_mismatches.append(f"ERROR: {dependency}: {dir_name1}: {version1} | {dir_name2}: {version2}")

    return version_mismatches

# Main function
def main():
    # Check if at least two directories are provided
    if len(os.sys.argv) < 3:
        print("Please provide at least two directory paths to compare.")
        exit(1)

    go_mod_paths = [os.path.join(dir_path, 'go.mod') for dir_path in os.sys.argv[1:]]

    # Filter out paths that do not have a go.mod file
    go_mod_paths = [path for path in go_mod_paths if os.path.isfile(path)]

    if len(go_mod_paths) < 2:
        print("ERROR: At least two go.mod files are required.")
        exit(1)

    version_mismatches = compare_go_mod(go_mod_paths)

    if version_mismatches:
        for mismatch in version_mismatches:
            print(mismatch)

if __name__ == "__main__":
    main()
