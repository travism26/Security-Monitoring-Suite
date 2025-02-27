#!/bin/bash
# cleanup-after-migration.sh - Script to clean up redundant files after migration

# Set the base directory
MEMLOG_DIR="$(pwd)"
echo "Starting cleanup in: $MEMLOG_DIR"

# Check if the new structure exists
if [ ! -d "$MEMLOG_DIR/active" ] || [ ! -d "$MEMLOG_DIR/archived" ] || [ ! -d "$MEMLOG_DIR/shared" ]; then
    echo "Error: New directory structure not found. Please run the migration script first."
    exit 1
fi

# Check if index.md exists
if [ ! -f "$MEMLOG_DIR/index.md" ]; then
    echo "Error: index.md not found. Please run the migration script first."
    exit 1
fi

# Create a backup directory
BACKUP_DIR="$MEMLOG_DIR/pre_migration_backup"
mkdir -p "$BACKUP_DIR"
echo "Created backup directory: $BACKUP_DIR"

# Move original log files to backup
echo "Moving original log files to backup..."
for log_file in "$MEMLOG_DIR"/*.log; do
    if [ -f "$log_file" ]; then
        filename=$(basename "$log_file")
        echo "  Backing up: $filename"
        mv "$log_file" "$BACKUP_DIR/"
    fi
done

# Move original shared tracking files to backup (only if they exist in shared directory)
echo "Moving original shared tracking files to backup..."
for shared_file in changelog.md stability_checklist.md url_debug_checklist.md; do
    if [ -f "$MEMLOG_DIR/$shared_file" ] && [ -f "$MEMLOG_DIR/shared/$shared_file" ]; then
        echo "  Backing up: $shared_file"
        mv "$MEMLOG_DIR/$shared_file" "$BACKUP_DIR/"
    fi
done

# Move reorganization proposal to backup
if [ -f "$MEMLOG_DIR/reorganization-proposal.md" ]; then
    echo "  Backing up: reorganization-proposal.md"
    mv "$MEMLOG_DIR/reorganization-proposal.md" "$BACKUP_DIR/"
fi

echo "Cleanup completed successfully!"
echo "All original files have been moved to: $BACKUP_DIR"
echo ""
echo "You can verify the new structure is working correctly, and if everything is fine,"
echo "you can remove the backup directory with:"
echo "  rm -rf $BACKUP_DIR"
echo ""
echo "Or keep it as a backup for reference."
