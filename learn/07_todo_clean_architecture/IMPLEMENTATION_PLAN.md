# Todo App - Implementation Plan

## Current State

**Backend:** Go (Gin) with clean architecture
**Frontend:** Coming later
**Database:** MySQL with migrations

### Existing Routes (`/api/v1`)
```
POST   /auth/register          - User signup
POST   /auth/login             - User login
GET    /users/profile          - Get user profile
PUT    /users/profile          - Update profile
GET    /todos                  - List user's todos
POST   /todos                  - Create todo
GET    /todos/:id              - Get single todo
PUT    /todos/:id              - Update todo
PATCH  /todos/:id/toggle       - Toggle completion
DELETE /todos/:id              - Delete todo
```

### Current Schema
```sql
todos (id, user_id, title, description, completed, priority, due_date, ...)
users (id, email, password, ...)
```

---

## Planned Changes

### Core Concept: Global Todos + Optional Lists
Transform from **only individual todos** to **global todos + organized lists**

```
User
 ├── Global Todos (list_id = NULL) ← Uncategorized/quick todos
 └── Lists (Optional organization)
      ├── "Work Projects"
      ├── "Shopping"
      └── "Personal"
           └── Todos (list_id = specific list)
```

---

## How It Works: Global Todos + Optional Lists

Todos can exist in two ways:
- **Uncategorized** (`list_id = NULL`) - Quick todos, no list needed
- **In a list** (`list_id = some-uuid`) - Organized todos

Users can:
- Create todos without assigning to a list
- Create lists when they want organization
- Move todos between lists or back to uncategorized
- Delete a list → **permanently deletes the list and all its todos**

---

## Migration

1. Create `todo_lists` table
2. Add `list_id` column to `todos` (nullable)
3. Existing todos remain uncategorized by default

---

## New Database Schema

### 1. Create `todo_lists` Table
```sql
CREATE TABLE todo_lists (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_lists (user_id, created_at)
);
```

### 2. Modify `todos` Table
```sql
ALTER TABLE todos ADD COLUMN list_id CHAR(36) NULL AFTER user_id;
ALTER TABLE todos ADD CONSTRAINT fk_todos_list_id
    FOREIGN KEY (list_id) REFERENCES todo_lists(id) ON DELETE CASCADE;
CREATE INDEX idx_todos_list_id ON todos(list_id);
```

**Key Details:**
- `list_id` is **NULLABLE** (NULL = global todos)
- `ON DELETE CASCADE`: Deleting list deletes all its todos
- No default list created automatically

---

## Migration Commands

### Migration File: `000003_add_todo_lists.up.sql`
```sql
-- Create todo_lists table
CREATE TABLE IF NOT EXISTS todo_lists (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_lists_user_id FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_user_lists (user_id, created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add list_id column to todos (nullable for global todos)
ALTER TABLE todos ADD COLUMN list_id CHAR(36) NULL AFTER user_id;

-- Add foreign key with CASCADE on delete (deletes todos with list)
ALTER TABLE todos ADD CONSTRAINT fk_todos_list_id
    FOREIGN KEY (list_id) REFERENCES todo_lists(id) ON DELETE CASCADE;

-- Add index for list queries
CREATE INDEX idx_todos_list_id ON todos(list_id);
```

### Migration File: `000003_add_todo_lists.down.sql`
```sql
-- Remove foreign key and index
ALTER TABLE todos DROP FOREIGN KEY fk_todos_list_id;
DROP INDEX idx_todos_list_id ON todos;

-- Remove list_id column
ALTER TABLE todos DROP COLUMN list_id;

-- Drop todo_lists table
DROP TABLE IF EXISTS todo_lists;
```

---

## New API Routes

### List Management (`/api/v1/lists`)
```
GET    /lists                  - Get all user's lists
POST   /lists                  - Create new list
GET    /lists/:id              - Get list with todos
PUT    /lists/:id              - Update list (rename)
DELETE /lists/:id              - Delete list and all its todos
POST   /lists/:id/duplicate    - Duplicate list + todos
```

### Updated Todo Routes (`/api/v1/todos`)
```
GET    /todos                  - List all todos (filter: ?list_id=xxx or ?global=true)
POST   /todos                  - Create todo (optional list_id in body)
GET    /todos/:id              - Get single todo
PUT    /todos/:id              - Update todo (can change list_id)
PATCH  /todos/:id/toggle       - Toggle completion
DELETE /todos/:id              - Delete todo
PATCH  /todos/move             - Move multiple todos to list or global
```

### Move Todos API
**Endpoint:** `PATCH /api/v1/todos/move`

**Handles all move scenarios:**
- List → Another List
- List → Global (uncategorize)
- Global → List (organize)

**Request Body:**
```json
{
  "todo_ids": ["uuid1", "uuid2", "uuid3"],
  "list_id": "list-uuid"  // or null for moving to global
}
```

---

## Implementation Checklist

### Phase 1: Database & Entities
- [ ] Create migration `000003_add_todo_lists.up.sql`
- [ ] Create migration `000003_add_todo_lists.down.sql`
- [ ] Create `TodoList` entity (`domain/entity/todo_list.go`)
- [ ] Update `Todo` entity with `ListID *uuid.UUID` (nullable pointer)

### Phase 2: Repository Layer
- [ ] Create `TodoListRepository` interface (`domain/repository/`)
- [ ] Create SQLC queries for lists (`internal/repository/queries/todo_lists.sql`)
- [ ] Implement repository (`internal/repository/sqlc_impl/todo_list_repository.go`)
- [ ] Update todo queries to support filtering by list_id/global

### Phase 3: Service Layer
- [ ] Create `TodoListService` interface (`domain/service/`)
- [ ] Implement `TodoListServiceImpl` with business logic:
  - Create/read/update/delete lists
  - Duplicate list (copy all todos)
  - Move todos between lists
  - Move todos to global (set list_id = NULL)
  - Get global todos (WHERE list_id IS NULL)

### Phase 4: DTOs & Handlers
- [ ] Create `ListDTO` structs (`internal/dto/list_dto.go`)
- [ ] Create `TodoListHandler` (`api/handler/todo_list_handler.go`)
- [ ] Update `TodoHandler` to support list_id parameter
- [ ] Add list routes to router

### Phase 5: Testing
- [ ] Unit tests for list service
- [ ] Integration tests for list endpoints
- [ ] Test edge cases:
  - Delete list → permanently deletes all todos (CASCADE)
  - Duplicate list with many todos
  - Filter global vs list todos

---

## Key Design Decisions

1. **Global Todos**: `list_id = NULL` for uncategorized todos
2. **No Forced Structure**: Users can use app without creating lists
3. **Soft Deletes**: `deleted_at` for lists (can restore later)
4. **ON DELETE CASCADE**: Deleting list permanently deletes all its todos
5. **Simple Schema**: No color, no position (keep it minimal)

---

## Example Usage Flow

```
1. User signs up → No lists yet
2. User creates todos → list_id = NULL (global todos)
3. User creates "Work" list → Now can organize
4. User moves 3 todos → From global to "Work" list
5. User creates "Shopping" list with 2 todos
6. User duplicates "Work" → Creates "Work (Copy)" with same todos
7. User deletes "Shopping" → Permanently deletes list + its 2 todos
8. User filters: /todos?global=true → Shows only global todos
```

---

---

## Frontend Overview

### Tech Stack
- **React 19** - UI framework
- **React Router 7** - Routing
- **React Query** - Server state management & caching
- **Zustand** - Auth token storage
- **Tailwind CSS** - Styling
- **Native Fetch** - HTTP requests

---

### Pages

| Route | Purpose |
|-------|---------|
| `/login` | User login |
| `/register` | User registration |
| `/todos` | Main dashboard (filtered by sidebar selection) |
| `/lists` | All lists overview (grid of cards) |
| `/lists/:id` | Individual list detail page |

**How filtering works:**
- Sidebar has "Global Todos" and list items
- Clicking sidebar items filters the `/todos` page
- No separate `/todos/global` route needed

---

### Features

**Todos:**
- Create, edit, delete todos
- Toggle completion
- Move multiple todos between lists or to global
- Filter by status (all/active/completed)

**Lists:**
- Create, rename, delete lists
- Duplicate list (copies all todos)
- Delete list permanently deletes all its todos

**State Management:**
- Server data (todos, lists) → React Query
- Auth token → Zustand + localStorage
- UI state (modals, filters) → Component state

---

## Next Steps

1. [DONE] Run migrations
2. [DONE] Generate SQLC code
3. Complete backend implementation
4. Test API endpoints
5. Build frontend
