# Changelog

All notable changes to the System Monitoring Gateway will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Implemented Authentication & Authorization system
  - Added JWT token validation middleware with tenant claims
  - Created API key validation system with tenant context
  - Implemented API key management endpoints:
    - Generate new API keys
    - List tenant's API keys
    - Revoke API keys
    - Rotate API keys
  - Added tenant context to request pipeline
  - Implemented middleware for protecting routes
  - Set up authentication flow:
    - JWT authentication for API key management
    - API key authentication for metrics endpoints

### Changed

### Deprecated

### Removed

### Fixed

### Security
