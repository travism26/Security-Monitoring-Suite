# AI Rules and Guidelines for Security-Monitoring-Suite Development

This document outlines the rules, best practices, and guidelines that AI must follow when assisting with the development of the Security-Monitoring-Suite project.

## Project Context

The Security-Monitoring-Suite is a comprehensive security monitoring system that provides:

- Real-time system monitoring
- Log aggregation
- Threat detection
- Security analytics capabilities

Key architectural features:

- Multi-tenant architecture supporting data isolation between organizations
- Event-driven design using Kafka for real-time data processing
- Microservices architecture with Go and TypeScript services
- Cloud-native deployment using Kubernetes

## Tech Stack

- Frontend: Next.js, React, TypeScript, Tailwind CSS, Shadcn UI
- Backend: Go (Log Aggregator, Monitoring Agent), TypeScript/Express (Gateway)
- Data Storage: PostgreSQL, MongoDB, Kafka
- Infrastructure: Kubernetes, Docker

## Code Style and Structure

### Directory Structure

```
log-aggregator/
├── cmd/                # Application entrypoints
├── internal/
    ├── domain/        # Core domain types
    ├── service/       # Business logic
    ├── handler/       # HTTP handlers
    ├── repository/    # Data access
    ├── kafka/         # Kafka integration
    └── config/        # Configuration

system-monitoring-gateway/
├── src/
    ├── routes/        # Express routes
    ├── services/      # Business logic
    ├── middleware/    # Express middleware
    ├── models/        # MongoDB models
    ├── types/         # TypeScript types
    └── kafka/         # Kafka integration

siem-dashboard/
├── app/
    ├── components/    # React components
    ├── hooks/         # Custom React hooks
    ├── contexts/      # React contexts
    ├── services/      # API services
    └── types/         # TypeScript types

system-monitoring-agent/
├── cmd/                # Application entrypoints
├── internal/
    ├── agent/         # Core agent functionality
    ├── metrics/       # Metrics collection
    ├── exporter/      # Data export handlers
    ├── threat/        # Threat analysis
    └── config/        # Configuration
```

### Naming Conventions

Go:

- Use PascalCase for exported names
- Use camelCase for internal names
- Use snake_case for file names
- Prefix test files with \_test.go

TypeScript/JavaScript:

- Use PascalCase for React components
- Use camelCase for functions and variables
- Use kebab-case for file names
- Suffix test files with .test.ts(x)

### TypeScript Usage

- Use TypeScript for all JavaScript code
- Define interfaces for all data structures
- Use discriminated unions for message types:

```typescript
interface BaseMessage {
  type: string;
}

interface MetricsMessage extends BaseMessage {
  type: "metrics";
  data: SystemMetrics;
}

interface AlertMessage extends BaseMessage {
  type: "alert";
  data: AlertData;
}

type Message = MetricsMessage | AlertMessage;
```

- Use proper type narrowing with type guards
- Avoid any type; use unknown for truly unknown types
- Use proper error handling with custom error types

### React Component Structure

- Use functional components with TypeScript interfaces
- Implement proper error boundaries
- Use React Context for global state
- Follow proper cleanup in useEffect hooks

Example:

```typescript
interface Props {
  data: SystemMetrics;
  onRefresh: () => void;
}

export function MetricsDisplay({ data, onRefresh }: Props): JSX.Element {
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const interval = setInterval(onRefresh, 5000);
    return () => clearInterval(interval);
  }, [onRefresh]);

  return (
    // JSX
  );
}
```

### UI and Styling

- Use Shadcn UI components with Tailwind CSS
- Follow consistent spacing and layout patterns
- Implement responsive design
- Use proper color schemes for alerts and status indicators
- Document new Shadcn component installations

## 1. Memlog System

### Directory Structure

The memlog system follows a hierarchical structure to optimize organization and token usage:

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

### Key Principles

- **Selective Loading**: Only load the index file and relevant active task files to minimize token usage
- **Regular Archiving**: Move completed tasks to archive files organized by project and time period
- **Standardized Format**: Follow consistent formatting for all task files
- **Cross-Referencing**: Use links between files rather than duplicating information

### Required Actions

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

4. **For New Projects**:
   - Create a new active task file following the standardized format
   - Add the project to the index file
   - Set up initial task structure and priorities

### File Maintenance

- Keep active files focused on current and upcoming tasks
- Maintain only the last 2 weeks of updates in active files
- Archive completed tasks quarterly or when files become too large
- Ensure all files follow the standardized format for consistency

## 2. Task Breakdown and Execution

- Break down all user instructions into clear, numbered steps
- Include both actions and reasoning for each step
- Flag potential issues before they arise
- Verify the completion of each step before proceeding to the next
- If errors occur, document them, return to previous steps, and retry as needed

## 3. Credential Management

- Explain the purpose of each credential when requesting from users
- Guide users to obtain any missing credentials
- Always test the validity of credentials before using them
- Never store credentials in plaintext; use environment variables
- Implement proper refresh procedures for expiring credentials
- Provide guidance on secure credential storage methods

## 4. Error Handling and Reporting

- Implement detailed and actionable error reporting
- Log errors with context and timestamps
- Provide users with clear steps for error recovery
- Track error history to identify patterns
- Implement escalation procedures for unresolved issues
- Ensure all systems have robust error handling mechanisms

## 5. Third-Party Services Integration

- Verify that the user has completed all setup requirements for each service
- Check all necessary permissions and settings
- Test service connections before using them in workflows
- Document version requirements and service dependencies
- Prepare contingency plans for potential service outages or failures

## 6. Testing and Quality Assurance

- Write unit tests for all business logic
- Implement integration tests for API endpoints
- Write end-to-end tests for critical flows
- Test error handling and edge cases
- Use proper mocking for external dependencies
- Maintain high test coverage and document it in the stability_checklist.md

## 7. Security Best Practices

- Implement proper authentication and authorization
- Use secure communication protocols (HTTPS)
- Handle sensitive data properly
- Follow security best practices for each technology
- Implement proper CORS and CSP policies
- Sanitize and validate all user inputs
- Follow the principle of least privilege

## 8. Performance Optimization

- Optimize database queries for efficiency
- Implement caching strategies where appropriate
- Minimize network requests and payload sizes
- Use asynchronous operations for I/O-bound tasks
- Profile the application to identify bottlenecks

## 9. Git Usage

Commit Message Prefixes:

- "fix:" for bug fixes
- "feat:" for new features
- "perf:" for performance improvements
- "docs:" for documentation changes
- "style:" for formatting changes
- "refactor:" for code refactoring
- "test:" for adding missing tests
- "chore:" for maintenance tasks

Rules:

- Use lowercase for commit messages
- Keep summary line under 50 characters
- Include detailed description for complex changes
- Reference issue numbers when applicable

## 10. Documentation

- Maintain clear README files
- Document API endpoints and data flows
- Keep configuration files well-documented
- Document complex business logic
- Maintain changelog for version updates
- Document new Shadcn component installations
- Update documentation with each feature update

Remember, these rules and guidelines must be followed without exception. Always refer back to this document when making decisions or providing assistance during the development process.
