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
      "created_at": "2026-02-05T16:37:42.188692504+05:30",
      "updated_at": "2026-02-05T16:37:42.188692504+05:30"
    },
    "expires_at": 1770376062
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
      "created_at": "2026-02-05T11:07:42Z",
      "updated_at": "2026-02-05T11:07:42Z"
    },
    "expires_at": 1770376067
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
    "created_at": "2026-02-05T11:07:42Z",
    "updated_at": "2026-02-05T11:07:42Z"
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
    "created_at": "2026-02-05T16:38:12.912501402+05:30",
    "updated_at": "2026-02-05T16:38:12.912501402+05:30",
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
    "due_date": "2024-01-20T17:00:00Z",
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
    "due_date": "2024-01-20T17:00:00Z",
    "created_at": "2026-02-05T16:38:12.912501402+05:30",
    "updated_at": "2026-02-05T16:38:12.912501402+05:30",
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
    "completed_at": "2024-01-10T14:30:00Z"
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
    "created_at": "2026-02-05T16:38:12.912501402+05:30",
    "updated_at": "2026-02-05T16:38:12.912501402+05:30",
    "completed_at": "2024-01-10T14:30:00Z",
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
    "created_at": "2026-02-05T16:38:43.764271555+05:30",
    "updated_at": "2026-02-05T16:38:43.764271555+05:30",
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
        "created_at": "2026-02-05T11:08:15Z",
        "updated_at": "2026-02-05T11:08:15Z",
        "is_overdue": false
      },
      {
        "id": "2a97f648-c496-4703-90ef-b094ac34c5ba",
        "title": "Buy groceries",
        "description": "Milk, eggs, bread",
        "completed": false,
        "priority": "medium",
        "created_at": "2026-02-05T11:08:13Z",
        "updated_at": "2026-02-05T11:08:13Z",
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
    "created_at": "2026-02-05T11:08:13Z",
    "updated_at": "2026-02-05T11:08:13Z",
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
    "created_at": "2026-02-05T11:08:13Z",
    "updated_at": "2026-02-05T16:38:38.342140887+05:30",
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
    "due_date": "2024-01-16T18:00:00Z"
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
    "due_date": "2024-01-16T18:00:00Z",
    "created_at": "2026-02-05T11:08:13Z",
    "updated_at": "2026-02-05T16:40:00.342140887+05:30",
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
    "completed_at": "2024-01-15T12:00:00Z"
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
    "due_date": "2024-01-16T18:00:00Z",
    "created_at": "2026-02-05T11:08:13Z",
    "updated_at": "2026-02-05T16:41:00.342140887+05:30",
    "completed_at": "2024-01-15T12:00:00Z",
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
    "due_date": "2024-01-16T18:00:00Z",
    "created_at": "2026-02-05T11:08:13Z",
    "updated_at": "2026-02-05T16:42:00.342140887+05:30",
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
    "created_at": "2026-02-05T11:08:13Z",
    "updated_at": "2026-02-05T16:38:41.17774522+05:30",
    "completed_at": "2026-02-05T16:38:41.17774522+05:30",
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
    "created_at": "2026-02-05T16:38:12.302840949+05:30",
    "updated_at": "2026-02-05T16:38:12.302840949+05:30"
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
        "created_at": "2026-02-05T11:08:14Z",
        "updated_at": "2026-02-05T11:08:14Z"
      },
      {
        "id": "62ff611b-b155-47f2-9476-cd4a0cad400c",
        "user_id": "ee293d6f-a11d-4f80-8eea-efd01a047383",
        "name": "Work Projects",
        "created_at": "2026-02-05T11:08:12Z",
        "updated_at": "2026-02-05T11:08:12Z"
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
    "created_at": "2026-02-05T11:08:12Z",
    "updated_at": "2026-02-05T11:08:12Z",
    "todos": [
      {
        "id": "fc34ab49-12f3-4536-bd43-3d62760786a3",
        "title": "Review pull requests",
        "description": "",
        "completed": false,
        "priority": "medium",
        "created_at": "2026-02-05T11:08:46Z",
        "updated_at": "2026-02-05T11:08:46Z",
        "is_overdue": false
      },
      {
        "id": "97c4fa49-925a-477a-b963-027b46788df7",
        "title": "Complete API documentation",
        "description": "",
        "completed": false,
        "priority": "high",
        "created_at": "2026-02-05T11:08:44Z",
        "updated_at": "2026-02-05T11:08:44Z",
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
    "created_at": "2026-02-05T11:08:12Z",
    "updated_at": "2026-02-05T16:39:17.301021756+05:30"
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
    "created_at": "2026-02-05T16:39:19.726317947+05:30",
    "updated_at": "2026-02-05T16:39:19.726317947+05:30",
    "todos": [
      {
        "id": "47cf02cd-50cb-4e45-ba84-889ea2320987",
        "title": "Review pull requests",
        "description": "",
        "completed": false,
        "priority": "medium",
        "created_at": "2026-02-05T16:39:19.730020476+05:30",
        "updated_at": "2026-02-05T16:39:19.730020476+05:30",
        "is_overdue": false
      },
      {
        "id": "aa0c091e-f837-464e-bfca-33adb394faa9",
        "title": "Complete API documentation",
        "description": "",
        "completed": false,
        "priority": "high",
        "created_at": "2026-02-05T16:39:19.73348959+05:30",
        "updated_at": "2026-02-05T16:39:19.73348959+05:30",
        "is_overdue": false
      }
    ]
  }
}
```

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
| `due_date` | ISO 8601 | No | When todo is due | `"2024-01-20T17:00:00Z"` |
| `completed` | boolean | No | Create as completed (default: false) | `true` |
| `completed_at` | ISO 8601 | No | Completion date (only if completed=true) | `"2024-01-15T12:00:00Z"` |
| `list_id` | UUID string | No | Assign to a list (null = global) | `"550e8400-e29b-..."` |

### UpdateTodoRequest Fields

All fields are **optional** - only provided fields will be updated.

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `title` | string | No | Update title (1-255 chars) | `"Buy groceries"` |
| `description` | string | No | Update description (max 2000 chars) | `"Need milk and eggs"` |
| `priority` | string | No | Update priority: `low`, `medium`, `high`, `urgent` | `"high"` |
| `due_date` | ISO 8601 | No | Update due date (null = remove) | `"2024-01-20T17:00:00Z"` |
| `completed` | boolean | No | Update completion status | `true` |
| `completed_at` | ISO 8601 | No | Update completion date | `"2024-01-15T12:00:00Z"` |

### Priority Values

- `low` - Low priority task
- `medium` - Medium priority task (default)
- `high` - High priority task
- `urgent` - Urgent/critical task

### Date Format

All dates must be in **ISO 8601** format with timezone:
- `2024-01-20T17:00:00Z` (UTC)
- `2024-01-20T17:00:00+05:30` (with timezone offset)

### Notes

1. **Partial Updates**: PUT endpoints only update fields you provide
2. **Null Values**: Use `null` (not empty string `""`) for optional fields
3. **Completed Logic**:
   - If `completed=true` without `completed_at`, current time is used
   - If `completed=false`, `completed_at` is cleared
4. **Overdue Calculation**: A todo is overdue if `due_date` is in the past AND `completed=false`
