import argparse

import os


def get_repo_name(file_path):
    # Split the file path using the '/' delimiter
    path_parts = file_path.split('/')

    # Find the index of 'go.mod' in the path_parts list
    try:
        index = path_parts.index('go.mod')

        # Return the element just before 'go.mod'
        if index > 0:
            return path_parts[index - 1]
    except ValueError:
        # 'go.mod' not found in the path
        pass

    return ''  # Return an empty string if 'go.mod' not found


file_path = "/home/aidan/code/go/test/repo1/go.mod"
repo_name = get_repo_name(file_path)


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
    version_mismatches = {}

    for i in range(len(go_mod_paths)):
        for j in range(i + 1, len(go_mod_paths)):
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
                    if dependency not in version_mismatches:
                        version_mismatches[dependency] = []
                    version_mismatches[dependency].append(f"DIFFERENT-VERSIONS: {dir_name1}: {version1} | {dir_name2}: {version2}")

    return version_mismatches


def update_versions(go_mod_paths, repo, new_version):
    for go_mod_path in go_mod_paths:
        with open(go_mod_path, 'r') as file:
            lines = file.readlines()

        with open(go_mod_path, 'w') as file:
            for line in lines:
                if line.strip().startswith(repo):
                    # Find the existing version number and replace it
                    parts = line.split()
                    if len(parts) >= 2:
                        existing_version = parts[1]
                        line = line.replace(existing_version, new_version, 1)
                        print(
                            f"REPO: {repo_name}  existing-version: {existing_version}   new-version: {new_version}    ")
                file.write(line)


# Main function
def main():
    # Create argument parser
    parser = argparse.ArgumentParser(description='Compare and update go.mod files.')
    parser.add_argument('directories', nargs='+', help='Directory paths to compare go.mod files.')
    parser.add_argument('--update-repo', help='Repository to update in go.mod files.')
    parser.add_argument('--update-version', help='New version to set in go.mod files.')

    # Parse command line arguments
    args = parser.parse_args()

    if args.update_repo and args.update_version:
        print(f"Updated {args.update_repo} to version {args.update_version}")

    # Check if at least two directories are provided
    if len(args.directories) < 2:
        print("Please provide at least two directory paths to compare.")
        exit(1)

    go_mod_paths = [os.path.join(dir_path, 'go.mod') for dir_path in args.directories]

    # Filter out paths that do not have a go.mod file
    go_mod_paths = [path for path in go_mod_paths if os.path.isfile(path)]

    if len(go_mod_paths) < 2:
        print("ERROR: At least two go.mod files are required.")
        exit(1)

    # Check if update arguments are provided
    if args.update_repo and args.update_version:
        update_versions(go_mod_paths, args.update_repo, args.update_version)

    else:
        version_mismatches = compare_go_mod(go_mod_paths)

        if version_mismatches:
            for dependency, mismatches in version_mismatches.items():
                print(f"{dependency}")
                for mismatch in mismatches:
                    print(mismatch)
        else:
            print("No version mismatches found.")


if __name__ == "__main__":
    main()
