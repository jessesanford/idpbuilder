#!/usr/bin/env python3
"""
Test script to verify terminal width detection in state-machine-viewer.py
"""

import shutil
import sys

def test_terminal_width():
    """Test that terminal width is being detected correctly"""
    try:
        width, height = shutil.get_terminal_size((80, 40))
        print(f"Terminal dimensions detected: {width} columns x {height} rows")

        # Test the centering function
        test_text = "TEST"
        centered = test_text.center(width)
        print(f"\nCentered text (using {width} columns):")
        print("|" + "=" * (width - 2) + "|")
        print("|" + centered[:width-2] + "|")
        print("|" + "=" * (width - 2) + "|")

        # Show different widths
        print("\nDynamic box widths based on terminal size:")
        print(f"  Small box (1/3 width):  {width // 3} columns")
        print(f"  Medium box (1/2 width): {width // 2} columns")
        print(f"  Large box (2/3 width):  {width * 2 // 3} columns")

        # Test with actual box drawing
        box_width = min(max(35, width // 3), width - 10)
        center_offset = max(0, (width - box_width - 2) // 2)

        print(f"\nBox centered with {box_width} columns width:")
        print(" " * center_offset + "┌" + "─" * box_width + "┐")
        print(" " * center_offset + "│" + "SAMPLE STATE".center(box_width) + "│")
        print(" " * center_offset + "└" + "─" * box_width + "┘")

        assert width > 0
        return True

    except Exception as e:
        print(f"Error detecting terminal size: {e}")
        assert False, f"Error detecting terminal size: {e}"

if __name__ == "__main__":
    success = test_terminal_width()
    sys.exit(0 if success else 1)