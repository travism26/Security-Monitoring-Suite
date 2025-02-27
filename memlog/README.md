# Memlog System

The Memlog System is a hierarchical task tracking and memory management system designed to efficiently organize project tasks, track progress, and maintain a historical record while optimizing token usage.

## Directory Structure

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

## Key Features

- **Hierarchical Organization**: Clear separation of active and archived tasks
- **Token Optimization**: Selective loading of only relevant files
- **Standardized Format**: Consistent structure across all task files
- **Automated Maintenance**: Scripts for common operations
- **Historical Record**: Complete project history in an accessible format

## Utility Scripts

The memlog system includes several utility scripts to help manage tasks:

### 1. custom-migrate.sh

Migrates from the old flat memlog structure to the new hierarchical structure, handling mixed project files by properly separating them.

```bash
cd memlog
./custom-migrate.sh
```

### 2. cleanup-after-migration.sh

Cleans up redundant files after migration by moving them to a backup directory.

```bash
cd memlog
./cleanup-after-migration.sh
```

### 3. migrate-memlog.sh

A simpler migration script for cases where projects are not mixed in the same file.

```bash
cd memlog
./migrate-memlog.sh
```

### 4. create-project.sh

Creates a new project task file with standardized format.

```bash
cd memlog
./create-project.sh [options] PROJECT_NAME
```

Options:

- `-d, --description DESC`: Add a project description
- `-s, --sprint SPRINT`: Set the initial sprint name/number (default: Sprint 1)
- `-p, --priority PRIORITY`: Set the default task priority (default: Medium)
- `-h, --help`: Show help message

Examples:

```bash
./create-project.sh new-project
./create-project.sh -d "API Service" api-service
./create-project.sh -s "Q1-Sprint" -p "High" new-project
```

### 3. archive-tasks.sh

Archives completed tasks from active files to archive files.

```bash
cd memlog
./archive-tasks.sh [options]
```

Options:

- `-p, --project PROJECT`: Specify a project to archive (default: all projects)
- `-d, --period PERIOD`: Specify the archive period (default: current quarter, e.g., 2025-Q1)
- `-a, --all`: Archive all completed tasks regardless of age
- `-h, --help`: Show help message

Examples:

```bash
./archive-tasks.sh
./archive-tasks.sh -p log-aggregator
./archive-tasks.sh -d 2024-Q4
./archive-tasks.sh -p siem-dashboard -d 2025-Q2
```

## Usage Guidelines

### For AI Assistants

1. **Before Starting Work**:

   - Check the index file (`memlog/index.md`) for project status overview
   - Load the relevant active task file for the specific project
   - Verify and update shared tracking files as needed

2. **During Task Execution**:

   - Update the active task file with progress, status changes, and new information
   - Add detailed timestamps for all updates
   - Cross-reference related tasks in other projects when applicable

3. **After Task Completion**:
   - Mark tasks as completed in the active file
   - Consider if completed tasks should be archived (based on age and relevance)
   - Update the index file to reflect current project status

### For Users

1. **Creating a New Project**:

   - Use the `create-project.sh` script to create a new project task file
   - Provide a descriptive name and optional details

2. **Archiving Completed Tasks**:

   - Periodically run the `archive-tasks.sh` script to move completed tasks to archive files
   - This helps keep active files focused and reduces token usage

3. **Finding Information**:
   - Start with the `index.md` file to get an overview of all projects
   - Navigate to specific project files in the `active` directory for current tasks
   - Check the `archived` directory for historical information

## File Format Standards

### Active Task Files

Active task files follow this structure:

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

## Recent Updates (Last 2 weeks)

[YYYY-MM-DD]

- Update 1
- Update 2

## Next Steps

1. Step 1
2. Step 2
```

### Archived Task Files

Archived task files follow this structure:

```markdown
# [Project Name] Archived Tasks - [Period]

## Completed Tasks

### 1. [Task Name] [Priority: High/Medium/Low]

Status: Completed

- [x] Subtask 1
- [x] Subtask 2

[YYYY-MM-DD]

- Completion details
- Performance metrics
```

## Benefits

1. **Reduced Token Usage**: By loading only relevant, current information
2. **Improved Organization**: Clear separation of active and archived tasks
3. **Better Information Retrieval**: Structured format makes finding information easier
4. **Scalability**: System can grow with the project without becoming unwieldy
5. **Historical Record**: Maintains complete project history in an accessible format
