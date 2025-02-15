# SIEM Dashboard Data Source Mapping

This document outlines the data sources and integration points for each UI component in the SIEM dashboard.

## 1. Event Log Component

### Data Source

- **Primary**: Log Aggregator (PostgreSQL)
- **Access Method**: RESTful API endpoints
- **Database**: log-aggregator-postgres-depl

### Data Flow

1. System metrics and logs ingested through Kafka
2. Processed by Log Aggregator service
3. Stored in PostgreSQL database
4. Retrieved via API endpoints

### Data Types

- System metrics (CPU, Memory usage)
- Process information
- Metadata and tags
- Timestamps and identifiers

### Required API Endpoints

```typescript
// Event Log API Interface
interface EventLogAPI {
  // Get paginated events with optional filters
  GET /api/v1/events?
    page={number}&
    limit={number}&
    startDate={ISO8601}&
    endDate={ISO8601}&
    severity={string}&
    type={string}

  // Get event details by ID
  GET /api/v1/events/{eventId}

  // Get event types for filtering
  GET /api/v1/events/types
}
```

## 2. Threat Summary Component

### Data Sources

- **Primary**: Threat Detection Simulation (TDS)
- **Secondary**: Mini-XDR System
- **Storage**: Log Aggregator PostgreSQL

### Data Flow

1. TDS generates synthetic security events
2. Events processed through API Gateway
3. Stored in Log Aggregator's PostgreSQL
4. Retrieved via API endpoints

### Data Types

- Malware incidents
- Phishing attempts
- DDoS attacks
- Severity levels (LOW, MEDIUM, HIGH, CRITICAL)

### Required API Endpoints

```typescript
// Threat Summary API Interface
interface ThreatAPI {
  // Get threat summary statistics
  GET /api/v1/threats/summary

  // Get threat details by category
  GET /api/v1/threats/{category}?
    timeframe={string}&
    severity={string}

  // Get threat incident details
  GET /api/v1/threats/incidents/{incidentId}
}
```

## 3. System Health Component

### Data Sources

- **Primary**: System Monitoring Gateway
- **Secondary**: Log Aggregator
- **Storage**: MongoDB (monitoring-gateway-mongo-depl)

### Data Flow

1. Real-time metrics collected from monitoring agents
2. Processed through monitoring gateway
3. Stored in MongoDB
4. Retrieved via API endpoints

### Metrics

- Firewall status
- IDS status
- Log server health
- Email filter status

### Required API Endpoints

```typescript
// System Health API Interface
interface SystemHealthAPI {
  // Get overall system health status
  GET /api/v1/health/status

  // Get component-specific health
  GET /api/v1/health/components/{componentId}

  // Get health history
  GET /api/v1/health/history?
    component={string}&
    timeframe={string}
}
```

## 4. Network Traffic Analysis

### Data Sources

- **Primary**: Network Protocol Analyzer
- **Secondary**: Log Aggregator
- **Storage**: Time-series data in PostgreSQL

### Data Flow

1. Network metrics collected by monitoring agents
2. Processed through monitoring gateway
3. Stored in time-series format
4. Retrieved via API endpoints

### Metrics

- Traffic volume
- Protocol distribution
- Source/destination patterns
- Anomaly indicators

### Required API Endpoints

```typescript
// Network Traffic API Interface
interface NetworkTrafficAPI {
  // Get real-time traffic metrics
  GET /api/v1/network/metrics/realtime

  // Get historical traffic data
  GET /api/v1/network/metrics/history?
    startTime={ISO8601}&
    endTime={ISO8601}&
    resolution={string}

  // Get traffic anomalies
  GET /api/v1/network/anomalies
}
```

## 5. Alert Management

### Data Source

- **Primary**: Log Aggregator's Alert Service
- **Storage**: PostgreSQL (alerts_schema)

### Data Flow

1. Alerts generated based on configured thresholds
2. Stored in PostgreSQL database
3. Retrieved via API endpoints
4. Real-time updates via WebSocket

### Data Types

- Alert severity
- Status (Open/Resolved)
- Related logs
- Timestamps

### Required API Endpoints

```typescript
// Alert Management API Interface
interface AlertAPI {
  // Get all alerts with filtering
  GET /api/v1/alerts?
    status={string}&
    severity={string}&
    startDate={ISO8601}&
    endDate={ISO8601}

  // Get alert details
  GET /api/v1/alerts/{alertId}

  // Update alert status
  PATCH /api/v1/alerts/{alertId}

  // Get alert statistics
  GET /api/v1/alerts/statistics
}
```

## 6. Identity & Access Management

### Data Source

- **Primary**: System Monitoring Gateway
- **Storage**: MongoDB (monitoring-gateway-mongo-depl)

### Data Flow

1. Authentication handled by gateway
2. User data stored in MongoDB
3. Session management via JWT
4. Retrieved via API endpoints

### Data Types

- User profiles
- Authentication tokens
- Role permissions
- Access logs

### Required API Endpoints

```typescript
// IAM API Interface
interface IAMApi {
  // User authentication
  POST /api/v1/auth/login
  POST /api/v1/auth/logout

  // User management
  GET /api/v1/users/me
  PATCH /api/v1/users/me

  // Role management
  GET /api/v1/roles
  GET /api/v1/roles/{roleId}/permissions
}
```

## Database Schema Overview

### PostgreSQL Tables

```sql
-- Events table
CREATE TABLE events (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP WITH TIME ZONE,
  type VARCHAR(50),
  source VARCHAR(100),
  severity VARCHAR(20),
  description TEXT,
  metadata JSONB
);

-- Alerts table
CREATE TABLE alerts (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP WITH TIME ZONE,
  severity VARCHAR(20),
  status VARCHAR(20),
  title VARCHAR(200),
  description TEXT,
  related_events JSONB
);

-- Network metrics table
CREATE TABLE network_metrics (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP WITH TIME ZONE,
  metric_type VARCHAR(50),
  value NUMERIC,
  metadata JSONB
);
```

### MongoDB Collections

```javascript
// Users collection
{
  _id: ObjectId,
  email: String,
  password: String,
  firstName: String,
  lastName: String,
  role: String,
  lastLogin: Date,
  metadata: Object
}

// System health collection
{
  _id: ObjectId,
  componentId: String,
  status: String,
  lastCheck: Date,
  metrics: Object,
  history: Array
}
```

## Integration Points

### Kafka Topics

1. `system.metrics` - Real-time system metrics
2. `security.events` - Security-related events
3. `network.traffic` - Network traffic data
4. `system.alerts` - Alert notifications

### WebSocket Channels

1. `/ws/alerts` - Real-time alert updates
2. `/ws/metrics` - Real-time metric updates
3. `/ws/events` - Real-time event stream

### Service Dependencies

1. Log Aggregator -> PostgreSQL
2. Monitoring Gateway -> MongoDB
3. TDS -> API Gateway -> Kafka
4. Network Analyzer -> Log Aggregator

## Authentication & Authorization

### JWT Token Structure

```javascript
{
  "sub": "user_id",
  "role": "user_role",
  "permissions": ["permission1", "permission2"],
  "exp": timestamp,
  "iat": timestamp
}
```

### Required Headers

```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

## Error Handling

All API endpoints should return standardized error responses:

```javascript
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": {} // Optional additional information
  }
}
```

Common error codes:

- `UNAUTHORIZED` - Authentication required
- `FORBIDDEN` - Insufficient permissions
- `NOT_FOUND` - Resource not found
- `VALIDATION_ERROR` - Invalid input
- `INTERNAL_ERROR` - Server error
