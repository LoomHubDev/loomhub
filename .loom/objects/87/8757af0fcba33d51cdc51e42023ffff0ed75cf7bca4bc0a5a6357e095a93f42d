# LoomHub — Testing Strategy

---

## Test Layers

### Unit Tests

Each package has `_test.go` files alongside the source. Focus:

| Package | What to Test |
|---------|-------------|
| `internal/auth` | JWT generation/validation, bcrypt, permission checks |
| `internal/store` | CRUD operations against in-memory SQLite |
| `internal/sync` | Negotiate logic, send/receive correctness |
| `internal/render` | Syntax highlighting, markdown output |
| `internal/webhook` | HMAC signing, payload formatting |

All store tests use a fresh in-memory SQLite database per test:

```go
func TestUserCreate(t *testing.T) {
    db := testutil.NewTestDB(t)
    s := store.NewUserStore(db)

    user, err := s.Create(context.Background(), store.CreateUserParams{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    })

    require.NoError(t, err)
    assert.Equal(t, "testuser", user.Username)
}
```

### Integration Tests

In `test/integration/`, test full HTTP request/response cycles:

```go
func TestSendReceive(t *testing.T) {
    srv := testutil.NewTestServer(t)

    // Create user + loom
    token := srv.CreateUser("testuser", "test@example.com")
    srv.CreateLoom(token, "my-app")

    // Send operations
    resp := srv.Send(token, "testuser/my-app", sendPayload)
    assert.Equal(t, 200, resp.StatusCode)

    // Receive and verify
    receiveResp := srv.Receive(token, "testuser/my-app", 0)
    assert.Len(t, receiveResp.Operations, len(sendPayload.Operations))
}
```

### End-to-End Tests

Test the actual `loom` CLI against a running LoomHub server:

```go
func TestCLISendReceive(t *testing.T) {
    srv := testutil.StartServer(t)

    // Init local loom project and send
    dir := t.TempDir()
    testutil.RunLoom(t, dir, "init")
    testutil.RunLoom(t, dir, "hub", "add", "origin", srv.URL+"/testuser/my-app")

    os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644)
    testutil.RunLoom(t, dir, "checkpoint", "initial")
    testutil.RunLoom(t, dir, "send", "origin")

    // Receive into a fresh project and verify
    dir2 := t.TempDir()
    testutil.RunLoom(t, dir2, "init")
    testutil.RunLoom(t, dir2, "hub", "add", "origin", srv.URL+"/testuser/my-app")
    testutil.RunLoom(t, dir2, "receive", "origin")

    // Verify content matches
    content, _ := os.ReadFile(filepath.Join(dir2, "main.go"))
    assert.Equal(t, "package main", string(content))
}
```

---

## Test Helpers

`test/testutil/helpers.go`:

```go
// NewTestDB creates a fresh in-memory SQLite database with hub schema
func NewTestDB(t *testing.T) *sql.DB

// NewTestServer creates a full test server with in-memory storage
func NewTestServer(t *testing.T) *TestServer

// NewTestLoomDB creates a fresh Loom database for a loom
func NewTestLoomDB(t *testing.T) *sql.DB
```

---

## What to Test for Each Feature

When adding a feature, test:

1. **Happy path** — Feature works as expected
2. **Auth** — Unauthenticated access is rejected for protected endpoints
3. **Permissions** — Users can't access/modify resources they shouldn't
4. **Validation** — Invalid input returns proper errors
5. **Edge cases** — Empty states, large payloads, concurrent access

---

## Running Tests

```bash
# All tests
go test ./...

# Verbose
go test -v ./...

# Specific package
go test ./internal/store/...

# Integration only
go test ./test/integration/

# With race detector
go test -race ./...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```
