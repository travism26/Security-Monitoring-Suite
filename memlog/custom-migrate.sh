#!/bin/bash
# custom-migrate.sh - Script to migrate the existing memlog system to the new hierarchical structure
# This script handles mixed project files and properly separates them

# Set the base directory
MEMLOG_DIR="$(pwd)"
echo "Starting custom memlog migration in: $MEMLOG_DIR"

# Create the new directory structure
echo "Creating new directory structure..."
mkdir -p "$MEMLOG_DIR/active"
mkdir -p "$MEMLOG_DIR/archived"
mkdir -p "$MEMLOG_DIR/shared"

# Move shared tracking files to the shared directory
echo "Moving shared tracking files..."
[ -f "$MEMLOG_DIR/changelog.md" ] && cp "$MEMLOG_DIR/changelog.md" "$MEMLOG_DIR/shared/"
[ -f "$MEMLOG_DIR/stability_checklist.md" ] && cp "$MEMLOG_DIR/stability_checklist.md" "$MEMLOG_DIR/shared/"
[ -f "$MEMLOG_DIR/url_debug_checklist.md" ] && cp "$MEMLOG_DIR/url_debug_checklist.md" "$MEMLOG_DIR/shared/"

# Define project mapping
declare -A PROJECT_MAP
PROJECT_MAP["log-aggregator"]="Log Aggregator"
PROJECT_MAP["system-monitoring-agent"]="System Monitoring Agent"
PROJECT_MAP["system-monitoring-gateway"]="System Monitoring Gateway"
PROJECT_MAP["siem-dashboard"]="SIEM Dashboard"
PROJECT_MAP["mini-xdr"]="Mini-XDR System"
PROJECT_MAP["network-protocol-analyzer"]="Network Protocol Analyzer"
PROJECT_MAP["threat-detection-simulation"]="Threat Detection Simulation"
PROJECT_MAP["documentation"]="Documentation"
PROJECT_MAP["ingress-implementation"]="Ingress Implementation"

# Process dedicated project files
echo "Processing dedicated project files..."
for task_file in "$MEMLOG_DIR"/*.log; do
    if [ -f "$task_file" ]; then
        filename=$(basename "$task_file")
        
        # Skip tasks.log as it will be processed separately
        if [ "$filename" == "tasks.log" ]; then
            continue
        fi
        
        # Extract project name from filename
        project_id=${filename%.tasks.log}
        project_id=${project_id%.*}  # Handle both formats: project.tasks.log and project-tasks.log
        
        # Skip if not a recognized project
        if [ -z "${PROJECT_MAP[$project_id]}" ]; then
            echo "Skipping unknown project file: $filename"
            continue
        fi
        
        project_name="${PROJECT_MAP[$project_id]}"
        echo "Processing $filename for project: $project_id ($project_name)"
        
        # Create active task file with standardized format
        active_file="$MEMLOG_DIR/active/$project_id.md"
        
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
        grep -A 50 "Status: In Progress\|Status: Pending\|Status: Blocked\|Status: Partially Completed" "$task_file" | 
            sed '/^$/d' | 
            sed '/^###/d' | 
            sed '/^Status: Completed/,$d' >> "$active_file" || true
        
        echo "" >> "$active_file"
        echo "## Recent Updates (Last 2 weeks)" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract recent updates (last 2 weeks)
        grep -A 10 "\[$(date +%Y-%m)" "$task_file" | head -n 20 >> "$active_file" || true
        
        echo "" >> "$active_file"
        echo "## Next Steps" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract next steps if they exist
        if grep -q "Next Steps" "$task_file"; then
            grep -A 10 "Next Steps" "$task_file" | tail -n +2 >> "$active_file" || true
        else
            echo "1. Review current tasks and prioritize" >> "$active_file"
            echo "2. Update task statuses" >> "$active_file"
        fi
        
        # Create archive directory for this project
        mkdir -p "$MEMLOG_DIR/archived/$project_id"
        
        # Create archive file for completed tasks
        archive_file="$MEMLOG_DIR/archived/$project_id/$(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 )).md"
        
        echo "# $project_name Archived Tasks - $(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 ))" > "$archive_file"
        echo "" >> "$archive_file"
        
        # Extract completed tasks
        if grep -q "Status: Completed" "$task_file"; then
            echo "## Completed Tasks" >> "$archive_file"
            echo "" >> "$archive_file"
            grep -A 50 "Status: Completed" "$task_file" | 
                sed '/^$/d' | 
                sed '/^###/d' | 
                sed '/^Status: In Progress/,$d' >> "$archive_file" || true
        else
            echo "No completed tasks found for archiving." >> "$archive_file"
        fi
    fi
done

# Process the mixed tasks.log file
echo "Processing mixed tasks.log file..."
if [ -f "$MEMLOG_DIR/tasks.log" ]; then
    # Extract Mini-XDR System tasks
    if grep -q "Mini-XDR System" "$MEMLOG_DIR/tasks.log"; then
        echo "Extracting Mini-XDR System tasks..."
        project_id="mini-xdr"
        project_name="${PROJECT_MAP[$project_id]}"
        
        # Create active task file
        active_file="$MEMLOG_DIR/active/$project_id.md"
        
        echo "# $project_name Active Tasks" > "$active_file"
        echo "" >> "$active_file"
        echo "## Current Sprint: Current" >> "$active_file"
        echo "" >> "$active_file"
        echo "Start Date: $(date +%Y-%m-%d)" >> "$active_file"
        echo "End Date: TBD" >> "$active_file"
        echo "" >> "$active_file"
        echo "## Active Tasks" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract Mini-XDR tasks that are not completed
        sed -n '/# Mini-XDR System/,/# Network Protocol Analyzer/p' "$MEMLOG_DIR/tasks.log" | 
            grep -A 50 "Status: In Progress\|Status: Pending\|Status: Blocked\|Status: Planned" | 
            sed '/^$/d' | 
            sed '/^###/d' | 
            sed '/^Status: Completed/,$d' >> "$active_file" || true
        
        echo "" >> "$active_file"
        echo "## Recent Updates (Last 2 weeks)" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract recent updates
        sed -n '/# Mini-XDR System/,/# Network Protocol Analyzer/p' "$MEMLOG_DIR/tasks.log" | 
            grep -A 10 "\[$(date +%Y-%m)" | head -n 20 >> "$active_file" || true
        
        echo "" >> "$active_file"
        echo "## Next Steps" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract next steps
        sed -n '/# Mini-XDR System/,/# Network Protocol Analyzer/p' "$MEMLOG_DIR/tasks.log" | 
            grep -A 10 "Next steps:" | tail -n +2 >> "$active_file" || true
        
        # Create archive directory and file
        mkdir -p "$MEMLOG_DIR/archived/$project_id"
        archive_file="$MEMLOG_DIR/archived/$project_id/$(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 )).md"
        
        echo "# $project_name Archived Tasks - $(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 ))" > "$archive_file"
        echo "" >> "$archive_file"
        echo "## Completed Tasks" >> "$archive_file"
        echo "" >> "$archive_file"
        
        # Extract completed tasks
        sed -n '/# Mini-XDR System/,/# Network Protocol Analyzer/p' "$MEMLOG_DIR/tasks.log" | 
            grep -A 50 "Status: Completed" | 
            sed '/^$/d' | 
            sed '/^###/d' | 
            sed '/^Status: In Progress/,$d' >> "$archive_file" || true
    fi
    
    # Extract Network Protocol Analyzer tasks
    if grep -q "Network Protocol Analyzer" "$MEMLOG_DIR/tasks.log"; then
        echo "Extracting Network Protocol Analyzer tasks..."
        project_id="network-protocol-analyzer"
        project_name="${PROJECT_MAP[$project_id]}"
        
        # Create active task file
        active_file="$MEMLOG_DIR/active/$project_id.md"
        
        echo "# $project_name Active Tasks" > "$active_file"
        echo "" >> "$active_file"
        echo "## Current Sprint: Current" >> "$active_file"
        echo "" >> "$active_file"
        echo "Start Date: $(date +%Y-%m-%d)" >> "$active_file"
        echo "End Date: TBD" >> "$active_file"
        echo "" >> "$active_file"
        echo "## Active Tasks" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract Network Protocol Analyzer tasks that are not completed
        sed -n '/# Network Protocol Analyzer/,/# Mini-XDR System Implementation Plan/p' "$MEMLOG_DIR/tasks.log" | 
            grep -A 50 "Status: Pending\|Status: Planned" | 
            sed '/^$/d' | 
            sed '/^###/d' | 
            sed '/^Status: Completed/,$d' >> "$active_file" || true
        
        echo "" >> "$active_file"
        echo "## Recent Updates (Last 2 weeks)" >> "$active_file"
        echo "" >> "$active_file"
        
        # No recent updates found for this project
        echo "No recent updates." >> "$active_file"
        
        echo "" >> "$active_file"
        echo "## Next Steps" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract next steps
        echo "1. Begin implementation of packet capture engine" >> "$active_file"
        echo "2. Set up project structure with Go modules" >> "$active_file"
        echo "3. Implement packet capture using gopacket" >> "$active_file"
        
        # Create archive directory and file
        mkdir -p "$MEMLOG_DIR/archived/$project_id"
        archive_file="$MEMLOG_DIR/archived/$project_id/$(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 )).md"
        
        echo "# $project_name Archived Tasks - $(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 ))" > "$archive_file"
        echo "" >> "$archive_file"
        echo "## Completed Tasks" >> "$archive_file"
        echo "" >> "$archive_file"
        echo "No completed tasks yet." >> "$archive_file"
    fi
    
    # Extract SIEM Dashboard tasks
    if grep -q "SIEM Dashboard Tasks" "$MEMLOG_DIR/tasks.log"; then
        echo "Extracting SIEM Dashboard tasks..."
        project_id="siem-dashboard"
        project_name="${PROJECT_MAP[$project_id]}"
        
        # Create active task file
        active_file="$MEMLOG_DIR/active/$project_id.md"
        
        echo "# $project_name Active Tasks" > "$active_file"
        echo "" >> "$active_file"
        echo "## Current Sprint: Current" >> "$active_file"
        echo "" >> "$active_file"
        echo "Start Date: $(date +%Y-%m-%d)" >> "$active_file"
        echo "End Date: TBD" >> "$active_file"
        echo "" >> "$active_file"
        echo "## Active Tasks" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract SIEM Dashboard tasks that are not completed
        sed -n '/# SIEM Dashboard Tasks/,/$/p' "$MEMLOG_DIR/tasks.log" | 
            grep -A 50 "Status: Pending\|Status: In Progress" | 
            sed '/^$/d' | 
            sed '/^###/d' | 
            sed '/^Status: Completed/,$d' >> "$active_file" || true
        
        echo "" >> "$active_file"
        echo "## Recent Updates (Last 2 weeks)" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract recent updates
        sed -n '/# SIEM Dashboard Tasks/,/$/p' "$MEMLOG_DIR/tasks.log" | 
            grep -A 10 "\[$(date +%Y-%m)" | head -n 20 >> "$active_file" || true
        
        echo "" >> "$active_file"
        echo "## Next Steps" >> "$active_file"
        echo "" >> "$active_file"
        
        # Add next steps
        echo "1. Add build step to deployment scripts" >> "$active_file"
        echo "2. Implement remaining security features" >> "$active_file"
        
        # Create archive directory and file
        mkdir -p "$MEMLOG_DIR/archived/$project_id"
        archive_file="$MEMLOG_DIR/archived/$project_id/$(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 )).md"
        
        echo "# $project_name Archived Tasks - $(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 ))" > "$archive_file"
        echo "" >> "$archive_file"
        echo "## Completed Tasks" >> "$archive_file"
        echo "" >> "$archive_file"
        
        # Extract completed tasks
        sed -n '/# SIEM Dashboard Tasks/,/$/p' "$MEMLOG_DIR/tasks.log" | 
            grep -A 50 "Status: Completed" | 
            sed '/^$/d' | 
            sed '/^###/d' | 
            sed '/^Status: In Progress/,$d' >> "$archive_file" || true
    fi
    
    # Extract Shared Infrastructure tasks
    if grep -q "Shared Infrastructure Components" "$MEMLOG_DIR/tasks.log"; then
        echo "Extracting Shared Infrastructure tasks..."
        project_id="shared-infrastructure"
        project_name="Shared Infrastructure"
        
        # Create active task file
        active_file="$MEMLOG_DIR/active/$project_id.md"
        
        echo "# $project_name Active Tasks" > "$active_file"
        echo "" >> "$active_file"
        echo "## Current Sprint: Current" >> "$active_file"
        echo "" >> "$active_file"
        echo "Start Date: $(date +%Y-%m-%d)" >> "$active_file"
        echo "End Date: TBD" >> "$active_file"
        echo "" >> "$active_file"
        echo "## Active Tasks" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract Shared Infrastructure tasks
        sed -n '/## Shared Infrastructure Components/,/## Next Steps/p' "$MEMLOG_DIR/tasks.log" | 
            grep -A 50 "Status: Planned" | 
            sed '/^$/d' | 
            sed '/^###/d' | 
            sed '/^Status: Completed/,$d' >> "$active_file" || true
        
        echo "" >> "$active_file"
        echo "## Recent Updates (Last 2 weeks)" >> "$active_file"
        echo "" >> "$active_file"
        
        # No recent updates found for this project
        echo "No recent updates." >> "$active_file"
        
        echo "" >> "$active_file"
        echo "## Next Steps" >> "$active_file"
        echo "" >> "$active_file"
        
        # Extract next steps
        echo "1. Implement OAuth2/OIDC authentication" >> "$active_file"
        echo "2. Set up distributed database cluster" >> "$active_file"
        echo "3. Set up centralized logging" >> "$active_file"
        
        # Create archive directory and file
        mkdir -p "$MEMLOG_DIR/archived/$project_id"
        archive_file="$MEMLOG_DIR/archived/$project_id/$(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 )).md"
        
        echo "# $project_name Archived Tasks - $(date +%Y)-Q$(( ($(date +%-m)-1)/3+1 ))" > "$archive_file"
        echo "" >> "$archive_file"
        echo "## Completed Tasks" >> "$archive_file"
        echo "" >> "$archive_file"
        echo "No completed tasks yet." >> "$archive_file"
    fi
fi

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
        project_id=${filename%.md}
        
        # Extract project name from file
        project_name=$(head -n 1 "$active_file" | sed 's/# \(.*\) Active Tasks/\1/')
        
        # Extract priority task if available
        priority_task=$(grep -A 1 "Priority: High" "$active_file" | grep -v "Priority" | head -n 1 | sed 's/^- \[ \] //')
        if [ -z "$priority_task" ]; then
            priority_task="None specified"
        fi
        
        # Extract status if available
        status=$(grep -A 1 "Status:" "$active_file" | head -n 1 | sed 's/Status: //')
        if [ -z "$status" ]; then
            status="In Progress"
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
        project_id=$(basename "$archive_dir")
        
        for archive_file in "$archive_dir"/*.md; do
            if [ -f "$archive_file" ]; then
                archive_filename=$(basename "$archive_file")
                period=${archive_filename%.md}
                
                # Extract project name from file
                project_name=$(head -n 1 "$archive_file" | sed 's/# \(.*\) Archived Tasks.*/\1/')
                
                # Extract a completed task if available
                completed_task=$(grep -A 1 "Status: Completed" "$archive_file" | grep -v "Status" | head -n 1 | sed 's/^- \[x\] //')
                if [ -z "$completed_task" ]; then
                    completed_task="Tasks archived"
                fi
                
                # Use current date as completion date
                echo "| $project_name | $completed_task | $(date +%Y-%m-%d) | [$period](./archived/$project_id/$archive_filename) |" >> "$index_file"
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

echo "Custom migration completed successfully!"
echo "Please review the new structure and make any necessary adjustments."
echo "You can now use the new memlog system as described in AI_RULES_AND_GUIDELINES.md."
