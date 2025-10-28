#!/bin/bash

# Simple YAML parser for setup configuration
# This is a helper script that can be sourced by setup-noninteractive.sh

# Function to parse YAML and export variables
parse_yaml() {
    local config_file="$1"
    local prefix="${2:-}"
    
    # Use awk for more reliable parsing
    awk -F': ' '
    BEGIN { 
        current_section = ""
        current_subsection = ""
    }
    
    # Skip comments and empty lines
    /^[[:space:]]*#/ { next }
    /^[[:space:]]*$/ { next }
    
    # Handle sections (no leading spaces)
    /^[a-zA-Z]/ {
        current_section = $1
        gsub(/:$/, "", current_section)  # Remove trailing colon from section name
        gsub(/:/, "_", current_section)  # Replace any colons in section name with underscore
        gsub(/-/, "_", current_section)  # Replace hyphens with underscore for valid bash vars
        current_subsection = ""
        if (NF > 1) {
            # Single line key: value
            gsub(/^[[:space:]]+/, "", $2)
            gsub(/[[:space:]]+$/, "", $2)
            gsub(/"/, "", $2)
            gsub(/[[:space:]]*#.*$/, "", $2)  # Remove inline comments
            print prefix current_section "=\"" $2 "\""
        }
        next
    }
    
    # Handle subsections (2 spaces)
    /^  [a-zA-Z]/ {
        gsub(/^[[:space:]]+/, "", $1)
        current_subsection = $1
        gsub(/:/, "_", current_subsection)  # Replace any colons in subsection name with underscore
        gsub(/-/, "_", current_subsection)  # Replace hyphens with underscore for valid bash vars
        if (NF > 1) {
            gsub(/^[[:space:]]+/, "", $2)
            gsub(/[[:space:]]+$/, "", $2)
            gsub(/"/, "", $2)
            gsub(/[[:space:]]*#.*$/, "", $2)  # Remove inline comments
            print prefix current_section "_" current_subsection "=\"" $2 "\""
        }
        next
    }
    
    # Handle array items (start with -)
    /^[[:space:]]*-/ {
        item = $0
        gsub(/^[[:space:]]*-[[:space:]]*/, "", item)
        gsub(/^"/, "", item)
        gsub(/"$/, "", item)
        gsub(/[[:space:]]*#.*$/, "", item)  # Remove inline comments
        
        if (current_subsection != "") {
            var_name = prefix current_section "_" current_subsection "_items"
        } else {
            var_name = prefix current_section "_items"
        }
        
        # Append to array
        print var_name "+=\" " item "\""
    }
    ' "$config_file"
}

# Test if run directly
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    if [ -z "$1" ]; then
        echo "Usage: $0 <yaml-file>"
        exit 1
    fi
    
    echo "# Parsed variables from $1:"
    parse_yaml "$1"
fi