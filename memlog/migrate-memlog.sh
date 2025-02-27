#!/bin/bash
# migrate-memlog.sh - Script to migrate the existing memlog system to the new hierarchical structure

# Set the base directory
MEMLOG_DIR="$(pwd)"
echo "Starting memlog migration in: $MEMLOG_DIR"

# Create the new directory structure
echo "Creating new directory structure..."
mkdir -p "$MEMLOG_DIR/active"
mkdir -p "$MEMLOG_DIR/archived"
mkdir -p "$MEMLOG_DIR/shared"

# Move shared tracking files to the shared directory
echo "Moving shared tracking files..."
[ -f "$MEMLOG_DIR/changelog.md" ] && mv "$MEMLOG_DIR/changelog.md" "$MEMLOG_DIR/shared/"
[ -f "$MEMLOG_DIR/stability_checklist.md" ] && mv "$MEMLOG_DIR/stability_checklist.md" "$MEMLOG_DIR/shared/"
[ -f "$MEMLOG_DIR/url_debug_checklist.md" ] && mv "$MEMLOG_DIR/url_debug_checklist.md" "$MEMLOG_DIR/shared/"

# Process task log files
echo "Processing task log files..."
for task_file in "$MEMLOG_DIR"/*.log; do
  if [ -f "$task_file" ]; then
    filename=$(basename "$task_file")
    project_name=${filename%.tasks.log}
    project_name=${project_name%.*}  # Handle both formats: project.tasks.log and project-tasks.log
    
    echo "Processing $filename for project: $project_name"
    
    # Create active task file with standardized format
    active_file="$MEMLOG_DIR/active/$project_name.md"
    
    # Extract current content
    echo "# $project_name Active Tasks" > "$active_file"
    echo "" >> "$active_file"
    echo "## Current Sprint: Current" >> "$active_file"
    echo "" >> "$active_file"
    echo "Start Date: $(date +%Y-%m-%d)" >> "$active_file"
    echo "End Date: TBD" >> "$active_file"
    echo "" >> "$active_file"
    echo "## Active Tasks" >> "$active_file"
    echo "" >> "$active_file"
    
    # Extract tasks that are not completed
    grep -A 20 "Status: In Progress\|Status: Pending\|Status: Blocked" "$task_file" | 
      sed '/^$/d' | 
      sed '/^###/d' | 
      sed '/^Status: Completed/,$d' >> "$active_file"
    
    echo "" >> "$active_file"
    echo "## Recent Updates (Last 2 weeks)" >> "$active_file"
    echo "" >> "$active_file"
    
    # Extract recent updates (last 2 weeks)
    grep -A 10 "\[$(date +%Y-%m)" "$task_file" | head -n 20 >> "$active_file"
    
    echo "" >> "$active_file"
    echo "## Next Steps" >> "$active_file"
    echo "" >> "$active_file"
    
    # Extract next steps if they exist
    if grep -q "Next Steps" "$task_file"; then
      grep -A 10 "Next Steps" "$task_file" | tail -n +2 >> "$active_file"
    else
      echo "1. Review current tasks and prioritize" >> "$active_file"
      echo "2. Update task statuses" >> "$active_file"
    fi
    
    # Create archive directory for this project
    mkdir -p "$MEMLOG_DIR/archived/$project_name"
    
    # Create archive file for completed tasks
    archive_file="$MEMLOG_DIR/archived/$project_name/$(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 )).md"
    
    echo "# $project_name Archived Tasks - $(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 ))" > "$archive_file"
    echo "" >> "$archive_file"
    
    # Extract completed tasks
    if grep -q "Status: Completed" "$task_file"; then
      echo "## Completed Tasks" >> "$archive_file"
      echo "" >> "$archive_file"
      grep -A 30 "Status: Completed" "$task_file" | 
        sed '/^$/d' | 
        sed '/^###/d' | 
        sed '/^Status: In Progress/,$d' >> "$archive_file"
    else
      echo "No completed tasks found for archiving." >> "$archive_file"
    fi
  fi
done

# Create index.md
echo "Creating index.md..."
index_file="$MEMLOG_DIR/index.md"

echo "# Memlog System Index" > "$index_file"
echo "" >> "$index_file"
echo "Last Updated: $(date +%Y-%m-%d)" >> "$index_file"
echo "" >> "$index_file"
echo "## Active Projects" >> "$index_file"
echo "" >> "$index_file"
echo "| Project | Current Sprint | Priority Tasks | Status |" >> "$index_file"
echo "|---------|---------------|----------------|--------|" >> "$index_file"

# Add entries for each active project
for active_file in "$MEMLOG_DIR/active"/*.md; do
  if [ -f "$active_file" ]; then
    filename=$(basename "$active_file")
    project_name=${filename%.md}
    
    # Extract priority task if available
    priority_task=$(grep -A 1 "Priority: High" "$active_file" | grep -v "Priority" | head -n 1 | sed 's/^- \[ \] //')
    if [ -z "$priority_task" ]; then
      priority_task="None specified"
    fi
    
    # Extract status if available
    status=$(grep -A 1 "Status:" "$active_file" | head -n 1 | sed 's/Status: //')
    if [ -z "$status" ]; then
      status="Unknown"
    fi
    
    echo "| [$project_name](./active/$filename) | Current | $priority_task | $status |" >> "$index_file"
  fi
done

echo "" >> "$index_file"
echo "## Recently Completed Tasks" >> "$index_file"
echo "" >> "$index_file"
echo "| Project | Task | Completion Date | Archive Link |" >> "$index_file"
echo "|---------|------|----------------|-------------|" >> "$index_file"

# Add entries for recently completed tasks
for archive_dir in "$MEMLOG_DIR/archived"/*; do
  if [ -d "$archive_dir" ]; then
    project_name=$(basename "$archive_dir")
    
    for archive_file in "$archive_dir"/*.md; do
      if [ -f "$archive_file" ]; then
        archive_filename=$(basename "$archive_file")
        period=${archive_filename%.md}
        
        # Extract a completed task if available
        completed_task=$(grep -A 1 "Status: Completed" "$archive_file" | grep -v "Status" | head -n 1 | sed 's/^- \[x\] //')
        if [ -z "$completed_task" ]; then
          completed_task="None specified"
        fi
        
        # Use current date as completion date
        echo "| $project_name | $completed_task | $(date +%Y-%m-%d) | [$period](./archived/$project_name/$archive_filename) |" >> "$index_file"
      fi
    done
  fi
done

echo "" >> "$index_file"
echo "## Shared Resources" >> "$index_file"
echo "" >> "$index_file"
echo "- [Changelog](./shared/changelog.md)" >> "$index_file"
echo "- [Stability Checklist](./shared/stability_checklist.md)" >> "$index_file"
echo "- [URL Debug Checklist](./shared/url_debug_checklist.md)" >> "$index_file"

echo "Migration completed successfully!"
echo "Please review the new structure and make any necessary adjustments."
echo "You can now use the new memlog system as described in AI_RULES_AND_GUIDELINES.md."
