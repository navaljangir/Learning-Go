# List Sharing Functionality

## Overview

The list sharing feature allows users to share their todo lists with other users by creating independent copies of the list and all its todos. When a list is shared, a completely new list is created for the target user with copies of all todos.

## API Endpoint

### Share a List

**Endpoint:** `POST /api/v1/lists/:id/share`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "target_user_id": "uuid-of-target-user",
  "custom_name": "Optional custom name for the shared list"
}
```

**Response:** `201 Created`
```json
{
  "id": "new-list-uuid",
  "user_id": "target-user-uuid",
  "name": "Shared List Name",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "todos": [
    {
      "id": "new-todo-uuid",
      "user_id": "target-user-uuid",
      "list_id": "new-list-uuid",
      "title": "Todo Title",
      "description": "Todo Description",
      "completed": false,
      "priority": "medium",
      "due_date": "2024-01-10T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## Usage Examples

### Example 1: Share with Auto-Generated Name

```bash
curl -X POST http://localhost:8080/api/v1/lists/abc-123/share \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target_user_id": "def-456"
  }'
```

The list will be named: "Original List Name (from username)"

### Example 2: Share with Custom Name

```bash
curl -X POST http://localhost:8080/api/v1/lists/abc-123/share \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target_user_id": "def-456",
    "custom_name": "Weekend Project Ideas"
  }'
```

## Why Create New Copies vs. Pointing to Original?

### ✅ Advantages of Creating New Copies (Current Implementation)

1. **Data Independence & Isolation**
   - Each user gets their own copy of the list and todos
   - Users can modify, complete, or delete todos without affecting the original
   - The original owner maintains full control of their list
   - No conflicts or unexpected changes from other users

2. **Privacy & Security**
   - Users can't see each other's progress or modifications
   - No need for complex permission systems (read-only, read-write, etc.)
   - Sensitive information in todos remains private to each user
   - No accidental data exposure between users

3. **Simplicity & Maintainability**
   - Simpler data model - no need for sharing permissions, access control lists
   - No complex queries to check who can see what
   - Easier to implement and maintain
   - Fewer edge cases to handle

4. **User Experience**
   - Users can customize their copy without restrictions
   - No confusion about why a todo suddenly changed
   - Clear ownership - each user owns their list completely
   - Freedom to adapt the list to their specific needs

5. **Performance**
   - No need to check permissions on every query
   - Simpler database queries
   - Better query performance (no JOIN on permissions table)
   - Easier to cache and optimize

6. **Deletion & Cleanup**
   - If one user deletes their list, it doesn't affect others
   - No orphaned data if the original user deletes their account
   - Clear data lifecycle management

### ❌ Disadvantages of Pointing to Original (Shared Reference)

1. **Complex Permission Management**
   - Need to implement: read-only, read-write, admin permissions
   - Must check permissions on every operation
   - Complicated to handle: who can add/remove users, who can delete, etc.

2. **Data Conflicts**
   - Multiple users editing the same todos simultaneously
   - Need to handle race conditions and concurrent updates
   - Confusing when someone marks your todo as complete
   - Version conflicts and synchronization issues

3. **Privacy Concerns**
   - All users see everyone's changes in real-time
   - Can't have personal notes or modifications
   - Difficult to implement selective sharing

4. **Performance Issues**
   - Every query needs permission checks
   - More complex JOIN operations
   - Harder to optimize and cache
   - Increased database load

5. **User Experience Issues**
   - Confusion about who changed what
   - Notifications needed for all changes
   - Conflicts over todo priority, due dates, etc.
   - Loss of personal control

## When Would Shared Reference Be Better?

The shared reference approach (pointing to original) would be more appropriate for:

1. **Collaborative Real-Time Projects**
   - Teams working together on the same task list
   - Shared household chores where completion status matters to everyone
   - Project management where everyone needs to see current status

2. **Real-Time Collaboration Requirements**
   - When users explicitly want to see each other's changes
   - When synchronization is a key feature
   - When the list represents shared resources or responsibilities

3. **Central Authority Use Cases**
   - Manager assigning tasks to team members
   - Teacher creating assignments for students
   - When one person controls the list content

## Implementation Notes

### Current Implementation Details

- **Service Layer**: `Share()` method in `TodoListService`
- **Handler**: `Share()` in `TodoListHandler`
- **Route**: `POST /api/v1/lists/:id/share`
- **Authentication**: Required - uses JWT middleware
- **Authorization**: Only list owner can share their lists

### Migration Tracking (Optional)

A `list_shares` table is available to track sharing history:

```sql
CREATE TABLE list_shares (
    id CHAR(36) PRIMARY KEY,
    source_list_id CHAR(36) NOT NULL,
    source_user_id CHAR(36) NOT NULL,
    target_list_id CHAR(36) NOT NULL,
    target_user_id CHAR(36) NOT NULL,
    shared_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ...
);
```

This table is optional and can be used for:
- Audit trails
- Analytics (who shares what)
- Recommendation systems
- Tracking list popularity

### Error Cases

1. **400 Bad Request**
   - Invalid list ID format
   - Trying to share with yourself
   - Missing target_user_id

2. **403 Forbidden**
   - List doesn't belong to the requesting user
   - Unauthorized access attempt

3. **404 Not Found**
   - List doesn't exist
   - Target user doesn't exist

4. **500 Internal Server Error**
   - Database errors
   - Failed to create list or todos

## Architecture Layers

The share functionality follows clean architecture principles:

```
Presentation Layer (api/handler)
    ↓
Application Layer (domain/service)
    ↓
Domain Layer (domain/entity, domain/repository)
    ↓
Infrastructure Layer (internal/repository/sqlc_impl)
```

Each layer has clear responsibilities and dependencies flow inward.

## Testing Considerations

When testing the share functionality:

1. Test user authorization (only owner can share)
2. Test target user validation (must exist)
3. Test todo duplication (all todos copied correctly)
4. Test custom name vs. auto-generated name
5. Test that modifications to shared list don't affect original
6. Test sharing empty lists
7. Test sharing lists with many todos (performance)
8. Test `keep_completed=true` preserves completed status and CompletedAt on duplicated/imported todos
9. Test `keep_completed=false` (default) resets all copied todos to incomplete

## Future Enhancements

If collaborative features are needed in the future, consider:

1. Add a separate "collaborative lists" feature alongside sharing
2. Implement real-time collaboration with WebSockets
3. Add granular permissions (viewer, editor, admin)
4. Implement change tracking and activity logs
5. Add comments and discussions on todos
6. Implement conflict resolution strategies
