#!/usr/bin/env bash

# ==============================================================================

# Function to display usage information
_usage() {
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    awk -f "$script_dir/usage.awk" "$0"
    exit 0
}

# Entry function to display development usage information
# Usage: ./<SCRIPT_NAME> usage_dev
# Description: Display detailed usage information for developers
_usage_dev() {
    echo "Development usage information for $0:"
    echo "  This script can be called with various commands to perform specific tasks."
    echo "  Use the following commands during development:"
    echo "    - entry_<command>: Executes the specified command."
    echo "    - main: The main function if defined."
    echo "    - -h, --help: Display this help message."
    exit 0
}

# Entry function to display usage template information
# Usage: ./<SCRIPT_NAME> usage_template
# Description: Display a template for writing entry functions with usage and description
_usage_template() {
    echo "Usage template information for $0:"
    echo "  This script uses the following template for commands:"
    echo "    - Entry function to <ACTION_DESCRIPTION>"
    echo "    - Usage: ./<SCRIPT_NAME> <COMMAND>"
    echo "    - Description: <DETAILED_DESCRIPTION>"
    echo
    echo "Examples:"
    echo
    echo "# Entry function to: <Messages>"
    echo "# Usage: ./<SCRIPT_NAME> <COMMAND>"
    echo "# Description: <Display into usage>"
    echo "entry_<COMMAND>() {"
    echo "    echo \"This is a template\""
    echo "}"
    echo
    echo "# Main function example"
    echo "# Description: Main function to execute the default behavior"
    echo "main() {"
    echo "    echo \"This is the main function\""
    echo "}"
    echo
    echo "Note: Replace <ACTION_DESCRIPTION>, <SCRIPT_NAME>, <COMMAND>, and <DETAILED_DESCRIPTION> with appropriate values."
    exit 0
}

usage() {
    if [[ "$1" == "usage_dev" ]]; then
        _usage_dev
    elif [[ "$1" == "usage_template" ]]; then
        _usage_template
    elif [[ "$1" == "-h" || "$1" == "--help" ]]; then
        _usage
    elif declare -F "entry_$1" >/dev/null; then
        "entry_$1" "${@:2}"
    elif [[ $# -eq 0 ]]; then
        if declare -F main >/dev/null; then
            main "${@}"
        else
            echo "Error: main function not found."
            exit 1
        fi
    else
        if declare -F main >/dev/null; then
            main "${@}"
        else
            echo "Error: Unknown command '$1'"
            exit 1
        fi
    fi
}

# ==============================================================================
