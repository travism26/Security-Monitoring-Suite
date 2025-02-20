# User Data Isolation Implementation Tasks

## Problem Statement

Currently, the log aggregator system has a potential data isolation issue where users without organizations could potentially access each other's data. The system needs to be modified to ensure proper data isolation by consistently using user_id for data filtering, while maintaining but not actively using organization_id for future use.

## Goals

1. Ensure complete user data isolation
2. Remove dependency on organization_id for data filtering
3. Maintain organization_id in the database for future use
4. Implement consistent user_id based filtering across all queries

## Detailed Task Breakdown

### 1. Repository Layer Changes (log_repository.go)

#### 1.1 Store Method

- [x] Remove organization_id conditional logic
- [x] Update INSERT query to always include user_id
- [x] Modify parameter ordering in query
- [x] Update args slice construction

```sql
INSERT INTO logs (
    id, api_key, user_id, timestamp, host, message, level, metadata,
    process_count, total_cpu_percent, total_memory_usage, organization_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
```

#### 1.2 StoreBatch Method

- [x] Update bulk insert query construction
- [x] Modify valueStrings generation
- [x] Update valueArgs slice construction
- [x] Ensure consistent parameter ordering

#### 1.3 FindByID Method

- [x] Remove organization_id parameter
- [x] Add user_id parameter
- [x] Update query to filter by user_id

```sql
SELECT id, api_key, organization_id, timestamp, host, message, level, metadata,
       process_count, total_cpu_percent, total_memory_usage
FROM logs
WHERE id = $1 AND user_id = $2
```

#### 1.4 List Method

- [x] Remove organization_id parameter
- [x] Update CTE query to always filter by user_id

```sql
WITH recent_logs AS (
    SELECT *,
    ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
    FROM logs
    WHERE user_id = $1
    ORDER BY timestamp DESC
)
```

#### 1.5 ListByTimeRange Method

- [x] Remove organization_id parameter
- [x] Add user_id parameter
- [x] Update time range query to include user_id filter

```sql
WHERE user_id = $1 AND timestamp >= $2 AND timestamp <= $3
```

#### 1.6 CountByTimeRange Method

- [x] Remove organization_id parameter
- [x] Add user_id parameter
- [x] Update count query

```sql
SELECT COUNT(*) FROM logs WHERE user_id = $1 AND timestamp >= $2 AND timestamp <= $3
```

#### 1.7 ListByHost Method

- [x] Remove organization_id parameter
- [x] Add user_id parameter
- [x] Update host filtering query

```sql
WHERE user_id = $1 AND host = $2
```

#### 1.8 ListByLevel Method

- [x] Remove organization_id parameter
- [x] Add user_id parameter
- [x] Update level filtering query

```sql
WHERE user_id = $1 AND level = $2
```

### 2. Domain Layer Changes (log.go)

#### 2.1 LogRepository Interface

- [x] Update method signatures to remove organization_id parameters
- [x] Add user_id parameters where missing
- [x] Update interface documentation

```go
type LogRepository interface {
    FindByID(userID, id string) (*Log, error)
    List(userID string, limit, offset int) ([]*Log, error)
    ListByTimeRange(userID string, start, end time.Time, limit, offset int) ([]*Log, error)
    // ... update other methods
}
```

### 3. Service Layer Changes (log_service.go)

#### 3.1 Configuration Updates

- [x] Remove MultiTenancyEnabled config option
- [x] Update configuration documentation

#### 3.2 Method Updates

- [x] Update GetLog method to use user_id
- [x] Modify ListLogs to remove organization_id dependency
- [x] Update ListByTimeRange to use user_id
- [x] Modify all other service methods to use user_id

#### 3.3 Cache Updates

- [x] Update cache key generation to use user_id instead of organization_id
- [x] Modify cache invalidation logic
- [x] Update cache key patterns in CacheKeyGenerator

### 4. Mock Repository Updates (mock/log_repository.go)

#### 4.1 Interface Implementation

- [x] Update mock repository method signatures
- [x] Modify mock implementation logic
- [x] Update mock return values

#### 4.2 Test Updates

- [x] Update existing test cases
- [x] Add new test cases for user isolation
- [x] Verify error cases

### 5. Testing Requirements

#### 5.1 Unit Tests

- [x] Test user data isolation in repository layer
- [x] Verify cache behavior with new user-based keys
- [x] Test error handling for invalid user IDs

#### 5.2 Integration Tests

- [x] Test complete user data isolation
- [x] Verify multi-user scenarios
- [x] Test pagination with user-specific data

#### 5.3 Performance Tests

- [x] Benchmark query performance with user_id filtering
- [x] Test cache hit rates with new key structure
- [x] Verify batch operation performance

### 6. Migration Considerations

#### 6.1 Database

- [ ] Keep organization_id column
- [ ] Ensure all existing queries work with new structure
- [ ] Verify indexes for user_id queries

#### 6.2 Data Migration

- [ ] No data migration needed
- [ ] organization_id remains in schema
- [ ] Existing data remains unchanged

## Implementation Order

1. Start with repository layer changes
2. Update domain interface
3. Modify service layer
4. Update mock repository
5. Add/update tests
6. Verify query performance
7. Deploy changes

## Rollback Plan

- Keep original code in version control
- Document all changes
- Maintain ability to revert to organization-based filtering if needed

## Future Considerations

1. Organization support can be re-enabled in the future
2. Current changes preserve organization_id data
3. Query optimization might be needed for user_id based filtering
4. Consider adding combined user_id + organization_id indexes for future use
