# Memlog System Reorganization Proposal

## Current Issues

1. **File Size Growth**: Task log files are becoming too large as they accumulate all historical tasks and updates.
2. **Token Inefficiency**: Loading large task files consumes unnecessary tokens during AI interactions.
3. **Lack of Archiving Structure**: No clear mechanism for archiving completed tasks or separating active from historical information.
4. **Information Retrieval Challenges**: Finding specific information in large files is becoming difficult.

## Proposed Solution: Hierarchical Memlog System

### 1. Directory Structure Reorganization

```
memlog/
├── active/                 # Active tasks and current status
│   ├── log-aggregator.md
│   ├── system-monitoring-gateway.md
│   ├── siem-dashboard.md
│   └── ...
├── archived/               # Archived completed tasks by project
│   ├── log-aggregator/
│   │   ├── 2024-Q1.md
│   │   └── 2024-Q2.md
│   ├── system-monitoring-gateway/
│   │   └── 2024-Q1.md
│   └── ...
├── shared/                 # Shared tracking files
│   ├── changelog.md
│   ├── stability_checklist.md
│   └── url_debug_checklist.md
└── index.md                # Master index of all projects and their status
```

### 2. Task File Structure Standardization

#### Active Task Files (`active/*.md`)

```markdown
# [Project Name] Active Tasks

## Current Sprint: [Sprint Name/Number]

Start Date: YYYY-MM-DD
End Date: YYYY-MM-DD

## Active Tasks

### 1. [Task Name] [Priority: High/Medium/Low]

Status: In Progress/Pending/Blocked

- [ ] Subtask 1
- [x] Subtask 2
- [ ] Subtask 3

### 2. [Task Name] [Priority: High/Medium/Low]

Status: In Progress/Pending/Blocked

- [ ] Subtask 1
- [ ] Subtask 2

## Recent Updates (Last 2 weeks)

[YYYY-MM-DD]

- Update 1
- Update 2

[YYYY-MM-DD]

- Update 3
- Update 4

## Next Steps

1. Step 1
2. Step 2
3. Step 3
```

#### Archived Task Files (`archived/[project]/[period].md`)

```markdown
# [Project Name] Archived Tasks - [Period]

## Sprint: [Sprint Name/Number]

Start Date: YYYY-MM-DD
End Date: YYYY-MM-DD

### 1. [Task Name] [Priority: High/Medium/Low]

Status: Completed

- [x] Subtask 1
- [x] Subtask 2
- [x] Subtask 3

[YYYY-MM-DD]

- Completion details
- Performance metrics
- Issues encountered and resolutions

### 2. [Task Name] [Priority: High/Medium/Low]

Status: Completed

- [x] Subtask 1
- [x] Subtask 2

[YYYY-MM-DD]

- Completion details
- Performance metrics
```

#### Index File (`index.md`)

```markdown
# Memlog System Index

Last Updated: YYYY-MM-DD

## Active Projects

| Project                                      | Current Sprint | Priority Tasks               | Status      |
| -------------------------------------------- | -------------- | ---------------------------- | ----------- |
| [Log Aggregator](./active/log-aggregator.md) | Sprint 5       | Multi-tenancy Implementation | In Progress |
| [SIEM Dashboard](./active/siem-dashboard.md) | Sprint 3       | Alert Visualization          | Blocked     |
| ...                                          | ...            | ...                          | ...         |

## Recently Completed Tasks

| Project        | Task                     | Completion Date | Archive Link                                    |
| -------------- | ------------------------ | --------------- | ----------------------------------------------- |
| Log Aggregator | Performance Optimization | 2024-01-24      | [2024-Q1](./archived/log-aggregator/2024-Q1.md) |
| SIEM Dashboard | Build Issues             | 2025-02-27      | [2025-Q1](./archived/siem-dashboard/2025-Q1.md) |
| ...            | ...                      | ...             | ...                                             |

## Shared Resources

- [Changelog](./shared/changelog.md)
- [Stability Checklist](./shared/stability_checklist.md)
- [URL Debug Checklist](./shared/url_debug_checklist.md)
```

### 3. Archiving Process

1. **Regular Archiving**: At the end of each sprint/month/quarter (configurable):

   - Move completed tasks from active files to appropriate archive files
   - Update the index with newly archived tasks
   - Keep only active and recent (last 2 weeks) updates in active files

2. **Automated Scripts**: Create simple scripts to assist with archiving:

```bash
# Example archiving script (conceptual)
#!/bin/bash
# archive_tasks.sh

PROJECT=$1
PERIOD=$2  # e.g., "2024-Q1"
ARCHIVE_FILE="memlog/archived/$PROJECT/$PERIOD.md"

# Ensure archive directory exists
mkdir -p "memlog/archived/$PROJECT"

# Extract completed tasks from active file
grep -A20 "Status: Completed" "memlog/active/$PROJECT.md" >> "$ARCHIVE_FILE"

# Remove completed tasks from active file
sed -i '/Status: Completed/,/^$/d' "memlog/active/$PROJECT.md"

# Update index
echo "Archived completed tasks for $PROJECT to $PERIOD"
```

### 4. Token Optimization Strategies

1. **Selective Loading**: AI should only load:

   - The index file for overall context
   - The specific active task file for the project being worked on
   - Relevant shared tracking files

2. **Summarization**: Include a "Summary" section at the top of each file with key points

3. **Reference Links**: Use reference links to point to archived information rather than including it inline

4. **Task Isolation**: When working on a specific task, create temporary "focused" files that only contain relevant information

### 5. Implementation Plan

1. **Create Directory Structure**: Set up the new directory hierarchy
2. **Migrate Existing Data**: Move current task information to the new structure
3. **Update AI Guidelines**: Modify AI_RULES_AND_GUIDELINES.md to reflect the new memlog system
4. **Create Helper Scripts**: Develop scripts for archiving and maintenance
5. **Test and Refine**: Test the new system with various AI interactions and refine as needed

## Benefits

1. **Reduced Token Usage**: By loading only relevant, current information
2. **Improved Organization**: Clear separation of active and archived tasks
3. **Better Information Retrieval**: Structured format makes finding information easier
4. **Scalability**: System can grow with the project without becoming unwieldy
5. **Historical Record**: Maintains complete project history in an accessible format

## Next Steps

1. Review this proposal and make any necessary adjustments
2. Create a migration plan for existing memlog data
3. Update the AI_RULES_AND_GUIDELINES.md with the new memlog system instructions
4. Implement the new directory structure and file formats
