#!/bin/bash

# Cleanup Backups Helper Script
# Helps users manage backup folder size by removing old backups

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Get the directory where this script lives
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Assume project root is parent of tools/
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

BACKUPS_DIR="$PROJECT_ROOT/backups"

# Check if backups directory exists
if [[ ! -d "$BACKUPS_DIR" ]]; then
    echo -e "${YELLOW}No backups directory found at: $BACKUPS_DIR${NC}"
    echo -e "${BLUE}No backups to clean up.${NC}"
    exit 0
fi

echo -e "${CYAN}${BOLD}"
cat << "EOF"
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║           Backup Cleanup Helper                      ║
║       Manage Software Factory Backups                ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}Current Backup Status${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"

# Show current backups
echo -e "\n${CYAN}Current backups:${NC}"
if ls "$BACKUPS_DIR"/* >/dev/null 2>&1; then
    du -h "$BACKUPS_DIR"/* 2>/dev/null | sort -k2 || echo -e "${YELLOW}No backups found${NC}"
else
    echo -e "${YELLOW}No backups found${NC}"
fi

echo ""

# Calculate total size
TOTAL_SIZE=$(du -sh "$BACKUPS_DIR" 2>/dev/null | cut -f1 || echo "0")
echo -e "${BOLD}Total backup size: ${CYAN}$TOTAL_SIZE${NC}"

# Count backups
BACKUP_COUNT=$(find "$BACKUPS_DIR" -maxdepth 1 \( -type d -o -name "*.tar.gz" \) | grep -v "^$BACKUPS_DIR$" | wc -l || echo "0")
echo -e "${BOLD}Number of backups: ${CYAN}$BACKUP_COUNT${NC}"
echo ""

if [[ $BACKUP_COUNT -eq 0 ]]; then
    echo -e "${GREEN}No backups to clean up${NC}"
    exit 0
fi

echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}Cleanup Options${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${CYAN}1.${NC} Delete all backups"
echo -e "${CYAN}2.${NC} Keep last N backups (delete older ones)"
echo -e "${CYAN}3.${NC} Delete backups older than N days"
echo -e "${CYAN}4.${NC} Cancel"
echo ""
read -p "Choose option (1-4): " OPTION

case $OPTION in
    1)
        echo ""
        echo -e "${RED}${BOLD}⚠️  WARNING: This will delete ALL backups!${NC}"
        read -p "Are you sure you want to delete ALL backups? (yes/no): " CONFIRM
        if [[ "$CONFIRM" =~ ^[Yy]([Ee][Ss])?$ ]]; then
            echo -e "${YELLOW}Deleting all backups...${NC}"
            rm -rf "$BACKUPS_DIR"/*
            echo -e "${GREEN}✓ All backups deleted${NC}"
        else
            echo -e "${BLUE}Cancelled - no backups deleted${NC}"
        fi
        ;;
    2)
        echo ""
        read -p "How many recent backups to keep? " KEEP_COUNT

        if ! [[ "$KEEP_COUNT" =~ ^[0-9]+$ ]]; then
            echo -e "${RED}Invalid number. Cancelled.${NC}"
            exit 1
        fi

        echo -e "${YELLOW}Deleting old backups (keeping $KEEP_COUNT most recent)...${NC}"

        # Get list of backups sorted by modification time (newest first)
        # Delete all except the N most recent
        find "$BACKUPS_DIR" -maxdepth 1 \( -type d -o -name "*.tar.gz" \) | \
            grep -v "^$BACKUPS_DIR$" | \
            sort -r | \
            tail -n +$((KEEP_COUNT + 1)) | \
            while read -r backup; do
                echo -e "  Deleting: ${backup#$BACKUPS_DIR/}"
                rm -rf "$backup"
            done

        echo -e "${GREEN}✓ Kept $KEEP_COUNT most recent backups${NC}"
        ;;
    3)
        echo ""
        read -p "Delete backups older than how many days? " DAYS

        if ! [[ "$DAYS" =~ ^[0-9]+$ ]]; then
            echo -e "${RED}Invalid number. Cancelled.${NC}"
            exit 1
        fi

        echo -e "${YELLOW}Deleting backups older than $DAYS days...${NC}"

        find "$BACKUPS_DIR" -maxdepth 1 \( -type d -o -name "*.tar.gz" \) -mtime +$DAYS | \
            grep -v "^$BACKUPS_DIR$" | \
            while read -r backup; do
                echo -e "  Deleting: ${backup#$BACKUPS_DIR/}"
                rm -rf "$backup"
            done

        echo -e "${GREEN}✓ Deleted backups older than $DAYS days${NC}"
        ;;
    4)
        echo -e "${BLUE}Cancelled - no changes made${NC}"
        exit 0
        ;;
    *)
        echo -e "${RED}Invalid option. Cancelled.${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}New Backup Status${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"

# Show updated status
NEW_SIZE=$(du -sh "$BACKUPS_DIR" 2>/dev/null | cut -f1 || echo "0")
NEW_COUNT=$(find "$BACKUPS_DIR" -maxdepth 1 \( -type d -o -name "*.tar.gz" \) | grep -v "^$BACKUPS_DIR$" | wc -l || echo "0")

echo -e "${BOLD}New total size: ${CYAN}$NEW_SIZE${NC}"
echo -e "${BOLD}Remaining backups: ${CYAN}$NEW_COUNT${NC}"
echo ""

if [[ $NEW_COUNT -gt 0 ]]; then
    echo -e "${CYAN}Remaining backups:${NC}"
    ls -lht "$BACKUPS_DIR" | head -n $((NEW_COUNT + 1)) | tail -n $NEW_COUNT
fi

echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}${BOLD}✓ Cleanup complete!${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
