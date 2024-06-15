#!/usr/bin/awk -f

BEGIN {
    FS = "\n"
    print "Usage: " ARGV[1] " [command]"
    main_found = 0
    in_main = 0
    in_entry = 0
    main_description = ""
    main_usage = ""
    all_usages = ""
}

/^# Description: / {
    description = substr($0, 15)
}

/^# Usage: / {
    usage = "    " substr($0, 9)
    all_usages = all_usages usage "\n"
}

function process_entry() {
    commands_info = commands_info sprintf("  \033[1;32m%s\033[0m \033[2;3m%s\033[0m\n", entry_name, description)
    if (all_usages != "") {
        commands_info = commands_info sprintf("%s", all_usages)
    }
}

function process_main() {
    main_description = description
    if (all_usages != "") {
        main_usage = main_usage all_usages
    }
}

function reset_usage_variables() {
    description = ""
    usage = ""
    all_usages = ""
}

# Main function detection
/^(function )?main\(/ {
    main_found = 1
    in_main = 1
    next
}

# Entry function detection
/^(function )?entry_/ {
    entry_name = $0
    sub(/^(function )?entry_/, "", entry_name)
    sub(/\(\) \{/, "", entry_name)
    in_entry = 1
    next
}

/^\}/ {
    if (in_main) {
        in_main = 0
        process_main()
        reset_usage_variables()
    } else if (in_entry) {
        in_entry = 0
        process_entry()
        reset_usage_variables()
    }
}

END {
    if (main_found) {
        print "\nMain: " main_description
        if (main_usage != "") {
            printf main_usage
        }
    }
    if (commands_info != "") {
        print "\nCommands:"
        printf "%s", commands_info
    }
}