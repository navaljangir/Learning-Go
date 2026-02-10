# API Examples with Curl Commands

Base URL: `http://localhost:8080`

---

## Authentication Endpoints

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser1",
    "email": "test1@example.com",
    "password": "SecurePass123!",
    "full_name": "Test User"
  }'
```

**Validation Rules:**
- `username`: 3-30 chars, starts with letter, only letters/numbers/underscores, no spaces
- `email`: valid email format, max 255 chars
- `password`: 8-72 chars, must contain uppercase, lowercase, number, and special character
- `full_name`: 2-100 chars (required)

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "eabe56d5-6a1c-48da-8351-b0129d9813ec",
      "username": "testuser1",
      "email": "test1@example.com",
      "full_name": "Test User",
      "created_at": "2026-02-09T12:05:10.188692504+05:30",
      "updated_at": "2026-02-09T12:05:10.188692504+05:30"
    },
    "expires_at": 1770705310
  }
}
```

---

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser1",
    "password": "SecurePass123!"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "eabe56d5-6a1c-48da-8351-b0129d9813ec",
      "username": "testuser1",
      "email": "test1@example.com",
      "full_name": "Test User",
      "created_at": "2026-02-09T06:35:10Z",
      "updated_at": "2026-02-09T06:35:10Z"
    },
    "expires_at": 1770705315
  }
}
```

---

## User Profile Endpoints

### Get Profile
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "eabe56d5-6a1c-48da-8351-b0129d9813ec",
    "username": "testuser1",
    "email": "test1@example.com",
    "full_name": "Test User",
    "created_at": "2026-02-09T06:35:10Z",
    "updated_at": "2026-02-09T06:35:10Z"
  }
}
```

---

## Todo Endpoints

### Create Todo - Minimal (Only Required Fields)
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer [TOKEN]" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy groceries",
    "priority": "medium"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "2a97f648-c496-4703-90ef-b094ac34c5ba",
    "title": "Buy groceries",
    "description": "",
    "completed": false,
    "priority": "medium",
    "created_at": "2026-02-09T12:05:20.912501402+05:30",
    "updated_at": "2026-02-09T12:05:20.912501402+05:30",
    "is_overdue": false
  }
}
```

---

### Create Todo - With All Optional Fields (Including List)
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer [TOKEN]" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write API docs and user guide",
    "priority": "high",
    "due_date": "2026-02-08T11:30:00Z",
    "completed": false,
    "list_id": "62ff611b-b155-47f2-9476-cd4a0cad400c"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "3b08f759-d507-5814-a1f0-c205bd45d6cb",
    "list_id": "62ff611b-b155-47f2-9476-cd4a0cad400c",
    "list_name": "Work Projects",
    "title": "Complete project documentation",
    "description": "Write API docs and user guide",
    "completed": false,
    "priority": "high",
    "due_date": "2026-02-08T11:30:00Z",
    "created_at": "2026-02-09T12:05:20.912501402+05:30",
    "updated_at": "2026-02-09T12:05:20.912501402+05:30",
    "is_overdue": true
  }
}
```

---

### Create Todo - Already Completed (Import from another system)
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer [TOKEN]" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Setup development environment",
    "description": "Install Go, MySQL, and configure project",
    "priority": "high",
    "completed": true,
    "completed_at": "2026-02-08T09:00:00Z"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "4c19g860-e618-6925-b2g1-d316ce56e7dc",
    "title": "Setup development environment",
    "description": "Install Go, MySQL, and configure project",
    "completed": true,
    "priority": "high",
    "created_at": "2026-02-09T12:05:20.912501402+05:30",
    "updated_at": "2026-02-09T12:05:20.912501402+05:30",
    "completed_at": "2026-02-08T09:00:00Z",
    "is_overdue": false
  }
}
```

---

### Create Todo (In a specific list)
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer [TOKEN]" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete API documentation",
    "priority": "high",
    "list_id": "62ff611b-b155-47f2-9476-cd4a0cad400c"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "97c4fa49-925a-477a-b963-027b46788df7",
    "title": "Complete API documentation",
    "description": "",
    "completed": false,
    "priority": "high",
    "created_at": "2026-02-09T12:05:30.764271555+05:30",
    "updated_at": "2026-02-09T12:05:30.764271555+05:30",
    "is_overdue": false
  }
}
```

---

### List All Todos
```bash
curl -X GET "http://localhost:8080/api/v1/todos" \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "todos": [
      {
        "id": "b7a8f38d-5c7e-462f-8655-2f8fd19144d7",
        "title": "Finish project",
        "description": "",
        "completed": false,
        "priority": "high",
        "created_at": "2026-02-09T06:35:15Z",
        "updated_at": "2026-02-09T06:35:15Z",
        "is_overdue": false
      },
      {
        "id": "2a97f648-c496-4703-90ef-b094ac34c5ba",
        "title": "Buy groceries",
        "description": "Milk, eggs, bread",
        "completed": false,
        "priority": "medium",
        "created_at": "2026-02-09T06:35:13Z",
        "updated_at": "2026-02-09T06:35:13Z",
        "is_overdue": false
      }
    ],
    "total": 2,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

---

### Get Single Todo
```bash
curl -X GET http://localhost:8080/api/v1/todos/2a97f648-c496-4703-90ef-b094ac34c5ba \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "2a97f648-c496-4703-90ef-b094ac34c5ba",
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "completed": false,
    "priority": "medium",
    "created_at": "2026-02-09T06:35:13Z",
    "updated_at": "2026-02-09T06:35:13Z",
    "is_overdue": false
  }
}
```

---

### Update Todo - Partial Update (Only Title)
```bash
curl -X PUT http://localhost:8080/api/v1/todos/2a97f648-c496-4703-90ef-b094ac34c5ba \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer [TOKEN]" \
  -d '{
    "title": "Buy groceries and fruits"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "2a97f648-c496-4703-90ef-b094ac34c5ba",
    "title": "Buy groceries and fruits",
    "description": "Milk, eggs, bread",
    "completed": false,
    "priority": "medium",
    "created_at": "2026-02-09T06:35:13Z",
    "updated_at": "2026-02-09T12:05:38.342140887+05:30",
    "is_overdue": false
  }
}
```

---

### Update Todo - Multiple Fields (with due_date)
```bash
curl -X PUT http://localhost:8080/api/v1/todos/2a97f648-c496-4703-90ef-b094ac34c5ba \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer [TOKEN]" \
  -d '{
    "title": "Buy groceries, fruits, and vegetables",
    "description": "Milk, eggs, bread, apples, carrots",
    "priority": "high",
    "due_date": "2026-02-08T12:30:00Z"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "2a97f648-c496-4703-90ef-b094ac34c5ba",
    "title": "Buy groceries, fruits, and vegetables",
    "description": "Milk, eggs, bread, apples, carrots",
    "completed": false,
    "priority": "high",
    "due_date": "2026-02-08T12:30:00Z",
    "created_at": "2026-02-09T06:35:13Z",
    "updated_at": "2026-02-09T12:05:55.342140887+05:30",
    "is_overdue": true
  }
}
```

---

### Update Todo - Mark as Completed with Custom Date
```bash
curl -X PUT http://localhost:8080/api/v1/todos/2a97f648-c496-4703-90ef-b094ac34c5ba \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer [TOKEN]" \
  -d '{
    "completed": true,
    "completed_at": "2026-02-08T06:30:00Z"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "2a97f648-c496-4703-90ef-b094ac34c5ba",
    "title": "Buy groceries, fruits, and vegetables",
    "description": "Milk, eggs, bread, apples, carrots",
    "completed": true,
    "priority": "high",
    "due_date": "2026-02-08T12:30:00Z",
    "created_at": "2026-02-09T06:35:13Z",
    "updated_at": "2026-02-09T12:05:57.342140887+05:30",
    "completed_at": "2026-02-08T06:30:00Z",
    "is_overdue": false
  }
}
```

---

### Update Todo - Mark as Incomplete
```bash
curl -X PUT http://localhost:8080/api/v1/todos/2a97f648-c496-4703-90ef-b094ac34c5ba \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer [TOKEN]" \
  -d '{
    "completed": false
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "2a97f648-c496-4703-90ef-b094ac34c5ba",
    "title": "Buy groceries, fruits, and vegetables",
    "description": "Milk, eggs, bread, apples, carrots",
    "completed": false,
    "priority": "high",
    "due_date": "2026-02-08T12:30:00Z",
    "created_at": "2026-02-09T06:35:13Z",
    "updated_at": "2026-02-09T12:05:59.342140887+05:30",
    "is_overdue": true
  }
}
```

---

### Toggle Todo Completion
```bash
curl -X PATCH http://localhost:8080/api/v1/todos/2a97f648-c496-4703-90ef-b094ac34c5ba/toggle \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "2a97f648-c496-4703-90ef-b094ac34c5ba",
    "title": "Buy groceries and fruits",
    "description": "Milk, eggs, bread",
    "completed": true,
    "priority": "medium",
    "created_at": "2026-02-09T06:35:13Z",
    "updated_at": "2026-02-09T12:05:41.17774522+05:30",
    "completed_at": "2026-02-09T12:05:41.17774522+05:30",
    "is_overdue": false
  }
}
```

---

### Delete Todo
```bash
curl -X DELETE http://localhost:8080/api/v1/todos/b7a8f38d-5c7e-462f-8655-2f8fd19144d7 \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "todo deleted successfully"
  }
}
```

---

## List Management Endpoints

### Create List
```bash
curl -X POST http://localhost:8080/api/v1/lists \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer [TOKEN]" \
  -d '{
    "name": "Work Projects"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "62ff611b-b155-47f2-9476-cd4a0cad400c",
    "user_id": "ee293d6f-a11d-4f80-8eea-efd01a047383",
    "name": "Work Projects",
    "created_at": "2026-02-09T12:05:12.302840949+05:30",
    "updated_at": "2026-02-09T12:05:12.302840949+05:30"
  }
}
```

---

### Get All Lists
```bash
curl -X GET http://localhost:8080/api/v1/lists \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "lists": [
      {
        "id": "03ed90f6-6998-41b3-9cd1-bd0b2dfc2db8",
        "user_id": "ee293d6f-a11d-4f80-8eea-efd01a047383",
        "name": "Personal Tasks",
        "created_at": "2026-02-09T06:35:14Z",
        "updated_at": "2026-02-09T06:35:14Z"
      },
      {
        "id": "62ff611b-b155-47f2-9476-cd4a0cad400c",
        "user_id": "ee293d6f-a11d-4f80-8eea-efd01a047383",
        "name": "Work Projects",
        "created_at": "2026-02-09T06:35:12Z",
        "updated_at": "2026-02-09T06:35:12Z"
      }
    ],
    "total": 2
  }
}
```

---

### Get List with Todos
```bash
curl -X GET http://localhost:8080/api/v1/lists/62ff611b-b155-47f2-9476-cd4a0cad400c \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "62ff611b-b155-47f2-9476-cd4a0cad400c",
    "user_id": "ee293d6f-a11d-4f80-8eea-efd01a047383",
    "name": "Work Projects",
    "created_at": "2026-02-09T06:35:12Z",
    "updated_at": "2026-02-09T06:35:12Z",
    "todos": [
      {
        "id": "fc34ab49-12f3-4536-bd43-3d62760786a3",
        "title": "Review pull requests",
        "description": "",
        "completed": false,
        "priority": "medium",
        "created_at": "2026-02-09T06:35:46Z",
        "updated_at": "2026-02-09T06:35:46Z",
        "is_overdue": false
      },
      {
        "id": "97c4fa49-925a-477a-b963-027b46788df7",
        "title": "Complete API documentation",
        "description": "",
        "completed": false,
        "priority": "high",
        "created_at": "2026-02-09T06:35:44Z",
        "updated_at": "2026-02-09T06:35:44Z",
        "is_overdue": false
      }
    ]
  }
}
```

---

### Update List (Rename)
```bash
curl -X PUT http://localhost:8080/api/v1/lists/62ff611b-b155-47f2-9476-cd4a0cad400c \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer [TOKEN]" \
  -d '{
    "name": "Work Projects (Updated)"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "62ff611b-b155-47f2-9476-cd4a0cad400c",
    "user_id": "ee293d6f-a11d-4f80-8eea-efd01a047383",
    "name": "Work Projects (Updated)",
    "created_at": "2026-02-09T06:35:12Z",
    "updated_at": "2026-02-09T12:05:50.301021756+05:30"
  }
}
```

---

### Duplicate List (with all todos)
```bash
curl -X POST http://localhost:8080/api/v1/lists/62ff611b-b155-47f2-9476-cd4a0cad400c/duplicate \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "e3f6185d-5766-40e5-a9a4-e27d5da8524c",
    "user_id": "ee293d6f-a11d-4f80-8eea-efd01a047383",
    "name": "Work Projects (Updated) (Copy)",
    "created_at": "2026-02-09T12:05:52.726317947+05:30",
    "updated_at": "2026-02-09T12:05:52.726317947+05:30",
    "todos": [
      {
        "id": "47cf02cd-50cb-4e45-ba84-889ea2320987",
        "title": "Review pull requests",
        "description": "",
        "completed": false,
        "priority": "medium",
        "created_at": "2026-02-09T12:05:52.730020476+05:30",
        "updated_at": "2026-02-09T12:05:52.730020476+05:30",
        "is_overdue": false
      },
      {
        "id": "aa0c091e-f837-464e-bfca-33adb394faa9",
        "title": "Complete API documentation",
        "description": "",
        "completed": false,
        "priority": "high",
        "created_at": "2026-02-09T12:05:52.73348959+05:30",
        "updated_at": "2026-02-09T12:05:52.73348959+05:30",
        "is_overdue": false
      }
    ]
  }
}
```

---

### Generate Share Link (Step 1: Owner generates a link)
```bash
curl -X POST http://localhost:8080/api/v1/lists/62ff611b-b155-47f2-9476-cd4a0cad400c/share \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "share_url": "/api/v1/lists/import/62ff611bb15547f29476cd4a0cad400ca1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "share_token": "62ff611bb15547f29476cd4a0cad400ca1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6"
  }
}
```

**Note:** No request body needed. The token is an HMAC-signed encoding of the list ID — no database storage required. The same list always produces the same token, so calling this endpoint multiple times is safe.

---

### Import Shared List (Step 2: Friend imports with the token)
```bash
curl -X POST http://localhost:8080/api/v1/lists/import/62ff611bb15547f29476cd4a0cad400ca1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6 \
  -H "Authorization: Bearer [FRIEND_TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4-e5f6-4789-b0c1-d2e3f4a5b6c7",
    "user_id": "f8d2b1c4-3e6a-4f9b-8c7d-1a2b3c4d5e6f",
    "name": "Work Projects (shared)",
    "created_at": "2026-02-09T12:05:55.123456789+05:30",
    "updated_at": "2026-02-09T12:05:55.123456789+05:30",
    "todos": [
      {
        "id": "b2c3d4e5-f6a7-4890-c1d2-e3f4a5b6c7d8",
        "title": "Review pull requests",
        "description": "",
        "completed": false,
        "priority": "medium",
        "created_at": "2026-02-09T12:05:55.234567890+05:30",
        "updated_at": "2026-02-09T12:05:55.234567890+05:30",
        "is_overdue": false
      },
      {
        "id": "c3d4e5f6-a7b8-4901-d2e3-f4a5b6c7d8e9",
        "title": "Complete API documentation",
        "description": "",
        "completed": false,
        "priority": "high",
        "created_at": "2026-02-09T12:05:55.345678901+05:30",
        "updated_at": "2026-02-09T12:05:55.345678901+05:30",
        "is_overdue": false
      }
    ]
  }
}
```

**How sharing works:**
1. Owner calls `POST /api/v1/lists/:id/share` → gets a share token (HMAC-signed, no DB storage)
2. Owner sends the URL to a friend
3. Friend calls `POST /api/v1/lists/import/:token` → list + todos are copied to their account
4. The import creates a completely independent copy — changes by either user won't affect the other
5. You cannot import your own list (use duplicate instead)

---

### Delete List (permanently deletes list and all todos)
```bash
curl -X DELETE http://localhost:8080/api/v1/lists/03ed90f6-6998-41b3-9cd1-bd0b2dfc2db8 \
  -H "Authorization: Bearer [TOKEN]"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "list deleted successfully"
  }
}
```

---

## Move Todos (Bulk Operations)

### Move todos from one list to another list
```bash
curl -X PATCH http://localhost:8080/api/v1/todos/move \
  -H "Authorization: Bearer [TOKEN]" \
  -H "Content-Type: application/json" \
  -d '{
    "todo_ids": ["1f98fc7e-e0fe-4857-887a-c01724805298", "b4e392db-9398-4169-adb2-0fd557533bc6"],
    "list_id": "f4c8c99e-541a-4aa6-840e-591f40c814a4"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "todos moved successfully"
  }
}
```

---

### Move todos from list to global (uncategorize)
```bash
curl -X PATCH http://localhost:8080/api/v1/todos/move \
  -H "Authorization: Bearer [TOKEN]" \
  -H "Content-Type: application/json" \
  -d '{
    "todo_ids": ["9fd5bbea-f08a-431f-9048-953d0de78a2e"],
    "list_id": null
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "todos moved successfully"
  }
}
```

---

### Move todos from global to a list (organize)
```bash
curl -X PATCH http://localhost:8080/api/v1/todos/move \
  -H "Authorization: Bearer [TOKEN]" \
  -H "Content-Type: application/json" \
  -d '{
    "todo_ids": ["a865d78d-f1ca-49fd-9b7f-614f8dec1a69"],
    "list_id": "edfa03f1-baf7-496d-8520-0d5b6f574049"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "todos moved successfully"
  }
}
```

---

## Field Reference

### CreateTodoRequest Fields

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `title` | string | **Yes** | Todo title (1-255 chars) | `"Buy groceries"` |
| `description` | string | No | Detailed description (max 2000 chars) | `"Need milk and eggs"` |
| `priority` | string | **Yes** | Priority: `low`, `medium`, `high`, `urgent` | `"high"` |
| `due_date` | ISO 8601 | No | When todo is due | `"2026-02-10T11:30:00Z"` |
| `completed` | boolean | No | Create as completed (default: false) | `true` |
| `completed_at` | ISO 8601 | No | Completion date (only if completed=true) | `"2026-02-08T06:30:00Z"` |
| `list_id` | UUID string | No | Assign to a list (null = global) | `"550e8400-e29b-..."` |

### UpdateTodoRequest Fields

All fields are **optional** - only provided fields will be updated.

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `title` | string | No | Update title (1-255 chars) | `"Buy groceries"` |
| `description` | string | No | Update description (max 2000 chars) | `"Need milk and eggs"` |
| `priority` | string | No | Update priority: `low`, `medium`, `high`, `urgent` | `"high"` |
| `due_date` | ISO 8601 | No | Update due date (null = remove) | `"2026-02-10T11:30:00Z"` |
| `completed` | boolean | No | Update completion status | `true` |
| `completed_at` | ISO 8601 | No | Update completion date | `"2026-02-08T06:30:00Z"` |

### Priority Values

- `low` - Low priority task
- `medium` - Medium priority task (default)
- `high` - High priority task
- `urgent` - Urgent/critical task

### Date Format

All dates must be in **ISO 8601** format with timezone:
- `2026-02-09T06:35:00Z` (UTC)
- `2026-02-09T12:05:00+05:30` (with timezone offset)

### Notes

1. **Partial Updates**: PUT endpoints only update fields you provide
2. **Null Values**: Use `null` (not empty string `""`) for optional fields
3. **Completed Logic**:
   - If `completed=true` without `completed_at`, current time is used
   - If `completed=false`, `completed_at` is cleared
4. **Overdue Calculation**: A todo is overdue if `due_date` is in the past AND `completed=false`
