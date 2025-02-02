# dev Notes

# Curl commands to test the API

1. Register a new user

```bash
curl -k -X POST https://security.dev/gateway/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "firstName": "Test",
    "lastName": "User"
  }'
```

2. Login (store the token)

```bash
TOKEN=$(curl -X POST https://security.dev/gateway/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' | jq -r '.token')
```

3. Get the user profile (use the token)

```bash
curl https://security.dev/gateway/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

4. Update the user profile (use the token)

```bash
curl -X PATCH https://security.dev/gateway/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Updated",
    "lastName": "Name"
  }'
```

5. Request password reset:

```bash
curl -X POST https://security.dev/gateway/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'
```

6. Health check endpoint

```bash
curl https://security.dev/health
```

## For admin operations (requires admin token):

1. Get all users

```bash
curl https://security.dev/gateway/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

2. Create new user (admin only):

```bash
curl -X POST https://security.dev/gateway/api/v1/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "password123",
    "firstName": "New",
    "lastName": "User",
    "role": "user"
  }'
```

Note: Add -k flag to bypass SSL verification if using self-signed certificates for local development. For better output formatting, you can pipe the response through jq: | jq '.'
