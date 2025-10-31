#!/usr/bin/env python3
"""Test script to demonstrate terminal width detection and adaptation"""

import subprocess
import os

def test_width(width):
    print(f"\n{'='*80}")
    print(f"Testing with terminal width: {width} columns")
    print(f"{'='*80}\n")

    env = os.environ.copy()
    env['COLUMNS'] = str(width)

    # Run the viewer with the specified width
    result = subprocess.run(
        ['python3', 'tools/state-machine-viewer.py', '--view', 'focused'],
        env=env,
        capture_output=True,
        text=True
    )

    # Show first 30 lines
    lines = result.stdout.split('\n')[:30]
    for line in lines:
        print(line)

    # Count the actual width used
    max_width = max(len(line.replace('\033[', '').replace('[0m', '').replace('[93m', '').replace('[95m', '').replace('[91m', '').replace('[1m', '').replace('[96m', '').replace('[32m', '').replace('[2m', ''))
                    for line in lines if line)
    print(f"\nActual max content width: {max_width} characters")

# Test with different widths
for width in [60, 80, 100, 120, 150]:
    test_width(width)