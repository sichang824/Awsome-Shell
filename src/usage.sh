#!/usr/bin/env bash

# ==============================================================================

# Function to display usage information
_usage() {
    echo "Usage: $0 [command]"

    # Check if the main function exists
    if declare -F main >/dev/null; then
        # Extract the description and usage for the main function
        main_start=$(grep -n "^main()" "$0" | cut -d: -f1)
        main_start=$((main_start - 1))
        main_end=$(awk "NR<$main_start" "$0" | awk '/^}/ {print NR}' | tail -n 1)

        # If main_end is empty, set it to 1 (start of the file)
        if [[ -z "$main_end" ]]; then
            main_end=1
        fi

        main_description=$(sed -n "$main_end,$main_start p" "$0" | grep "^# Description:" | sed 's/^# Description: //')
        main_usage=$(sed -n "$main_end,$main_start p" "$0" | grep "^# Usage:" | sed 's/^# Usage: //')

        echo
        echo -e "Main:\t$main_description"
        if [[ -n "$main_usage" ]]; then
            while IFS= read -r example; do
                echo "  Usage: $example"
            done <<<"$main_usage"
        fi
    fi

    echo
    echo "Commands:"

    declare -F | awk '{print $3}' | while read -r func; do
        if [[ $func == entry_* ]]; then
            # Extract the description and usage from the function definition
            func_start=$(grep -n "^$func()" "$0" | cut -d: -f1)
            func_start=$((func_start - 1))
            func_end=$(awk "NR<$func_start" "$0" | awk '/^}/ {print NR}' | tail -n 1)

            # If func_end is empty, set it to 1 (start of the file)
            if [[ -z "$func_end" ]]; then
                func_end=1
            fi

            description=$(sed -n "$func_end,$func_start p" "$0" | grep "^# Description:" | sed 's/^# Description: //')
            usage_examples=$(sed -n "$func_end,$func_start p" "$0" | grep "^# Usage:" | sed 's/^# Usage: //')
            echo "  ${func#entry_}     $description"
            if [[ -n "$usage_examples" ]]; then
                while IFS= read -r example; do
                    echo "    Usage: $example"
                done <<<"$usage_examples"
            fi
        fi
    done
    exit 1
}

usage() {
    # Check if the function exists and call it
    if [[ $# -eq 0 ]]; then
        if declare -F main >/dev/null; then
            main "${@}"
        else
            echo "Error: main function not found."
            exit 1
        fi
    elif [[ "$1" == "-h" || "$1" == "--help" ]]; then
        _usage
    elif declare -F "entry_$1" >/dev/null; then
        "entry_$1" "${@:2}"
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
