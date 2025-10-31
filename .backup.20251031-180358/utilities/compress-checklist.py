#!/usr/bin/env python3
"""
Compress SF 3.0 checklist by reducing completed items to single lines.
Preserves: Phase 6 "Validate integration workflows" item (lines 810-822)
"""

import re
import sys

def compress_completed_item(lines, start_idx):
    """
    Compress a completed item to a single line, preserving all key information.

    Returns: (compressed_line, end_idx)
    """
    # First line is the checkbox with task
    first_line = lines[start_idx].rstrip()

    # Extract key info from subsequent lines
    validation_results = None
    summary = None

    i = start_idx + 1
    while i < len(lines):
        line = lines[i].rstrip()

        # Stop at next checkbox or section header
        if line.startswith('- [') or line.startswith('##') or line.startswith('####'):
            break

        # Extract only the most important information for single-line format
        if '**Validation Results**:' in line:
            validation_results = line.split('**Validation Results**:')[1].strip()
        elif '**Summary**:' in line:
            summary = line.split('**Summary**:')[1].strip()
        elif '**Results**:' in line and not validation_results:
            validation_results = line.split('**Results**:')[1].strip()

        i += 1

    # Build compressed line - just keep the first line since it already has date, owner, checkmark
    compressed = first_line

    # Only add summary if it provides additional value not in first line
    if summary and len(summary) < 150 and summary not in first_line:
        compressed += f" [{summary}]"
    elif validation_results and len(validation_results) < 100 and validation_results not in first_line:
        compressed += f" [{validation_results}]"

    return compressed, i - 1

def should_preserve(line_num, line_content):
    """Check if this item should be preserved as multi-line"""
    # Preserve "Validate integration workflows" item (around line 810)
    if 'Validate integration workflows' in line_content or 'validate.*integration.*workflow' in line_content.lower():
        return True
    return False

def compress_checklist(input_file, output_file):
    """Main compression function"""
    with open(input_file, 'r') as f:
        lines = f.readlines()

    compressed = []
    i = 0

    while i < len(lines):
        line = lines[i].rstrip()

        # Check if this is a completed item
        if line.startswith('- [X]') or line.startswith('- [x]'):
            # Check if should preserve
            if should_preserve(i, line):
                # Keep as-is - copy this item and all its details
                compressed.append(lines[i])
                i += 1
                # Copy all detail lines until next item or section
                while i < len(lines):
                    next_line = lines[i].rstrip()
                    if next_line.startswith('- [') or next_line.startswith('##') or next_line.startswith('####'):
                        break
                    compressed.append(lines[i])
                    i += 1
            else:
                # Compress to single line
                compressed_line, end_idx = compress_completed_item(lines, i)
                compressed.append(compressed_line + '\n')
                i = end_idx + 1
        else:
            # Not a completed item - keep as-is
            compressed.append(lines[i])
            i += 1

    # Write output
    with open(output_file, 'w') as f:
        f.writelines(compressed)

    return len(lines), len(compressed)

if __name__ == '__main__':
    input_file = 'docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md'
    output_file = 'docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md.compressed'

    original_lines, compressed_lines = compress_checklist(input_file, output_file)

    print(f"✅ Compression complete:")
    print(f"   Original: {original_lines} lines")
    print(f"   Compressed: {compressed_lines} lines")
    print(f"   Reduction: {original_lines - compressed_lines} lines ({100*(original_lines - compressed_lines)/original_lines:.1f}%)")
    print(f"   Output: {output_file}")

    if compressed_lines > 25000:
        print(f"⚠️  WARNING: Still above 25k lines target")
        sys.exit(1)
    else:
        print(f"✅ Target met: {compressed_lines} < 25000 lines")
