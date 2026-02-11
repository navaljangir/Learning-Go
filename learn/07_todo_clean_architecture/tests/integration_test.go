package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==========================================================================
// Auth Flow Tests
// ==========================================================================

func TestRegisterAndUseToken(t *testing.T) {
	app := NewTestApp()

	// Register a new user
	token := app.RegisterUser(t, "alice", "alice@example.com", "Passw0rd!", "Alice Smith")
	assert.NotEmpty(t, token)

	// Use the token on a protected endpoint (get profile)
	rec := app.DoRequest("GET", "/api/v1/users/profile", token, nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	resp := ParseResponse(t, rec)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "alice", data["username"])
	assert.Equal(t, "Alice Smith", data["full_name"])
}

func TestRegisterDuplicateUsername(t *testing.T) {
	app := NewTestApp()

	app.RegisterUser(t, "bob", "bob@example.com", "Passw0rd!", "Bob Jones")

	// Try registering with the same username
	rec := app.DoRequest("POST", "/api/v1/auth/register", "", map[string]string{
		"username":  "bob",
		"email":     "bob2@example.com",
		"password":  "Passw0rd!",
		"full_name": "Bob Other",
	})
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterDuplicateEmail(t *testing.T) {
	app := NewTestApp()

	app.RegisterUser(t, "carol", "carol@example.com", "Passw0rd!", "Carol Lee")

	rec := app.DoRequest("POST", "/api/v1/auth/register", "", map[string]string{
		"username":  "carol2",
		"email":     "carol@example.com",
		"password":  "Passw0rd!",
		"full_name": "Carol Other",
	})
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginWrongPassword(t *testing.T) {
	app := NewTestApp()

	app.RegisterUser(t, "dave", "dave@example.com", "Passw0rd!", "Dave King")

	rec := app.DoRequest("POST", "/api/v1/auth/login", "", map[string]string{
		"username": "dave",
		"password": "WrongPass1!",
	})
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginNonexistentUser(t *testing.T) {
	app := NewTestApp()

	rec := app.DoRequest("POST", "/api/v1/auth/login", "", map[string]string{
		"username": "ghost",
		"password": "Passw0rd!",
	})
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestProtectedRouteWithoutToken(t *testing.T) {
	app := NewTestApp()

	rec := app.DoRequest("GET", "/api/v1/users/profile", "", nil)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestProtectedRouteWithInvalidToken(t *testing.T) {
	app := NewTestApp()

	rec := app.DoRequest("GET", "/api/v1/users/profile", "invalid-token-here", nil)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginAfterRegister(t *testing.T) {
	app := NewTestApp()

	app.RegisterUser(t, "eve", "eve@example.com", "Passw0rd!", "Eve White")

	// Login should return a valid token
	token := app.LoginUser(t, "eve", "Passw0rd!")
	assert.NotEmpty(t, token)

	// The login token should also work on protected routes
	rec := app.DoRequest("GET", "/api/v1/users/profile", token, nil)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRegisterValidation(t *testing.T) {
	app := NewTestApp()

	// Missing required fields
	rec := app.DoRequest("POST", "/api/v1/auth/register", "", map[string]string{})
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Weak password
	rec = app.DoRequest("POST", "/api/v1/auth/register", "", map[string]string{
		"username":  "weakuser",
		"email":     "weak@example.com",
		"password":  "weak",
		"full_name": "Weak User",
	})
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ==========================================================================
// User Profile Tests
// ==========================================================================

func TestUpdateProfile(t *testing.T) {
	app := NewTestApp()

	token := app.RegisterUser(t, "frank", "frank@example.com", "Passw0rd!", "Frank Old")

	rec := app.DoRequest("PUT", "/api/v1/users/profile", token, map[string]string{
		"full_name": "Frank Updated",
	})
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the update
	rec = app.DoRequest("GET", "/api/v1/users/profile", token, nil)
	resp := ParseResponse(t, rec)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "Frank Updated", data["full_name"])
}

// ==========================================================================
// Todo CRUD Tests
// ==========================================================================

func TestTodoCRUDFlow(t *testing.T) {
	app := NewTestApp()
	token := app.RegisterUser(t, "todouser", "todo@example.com", "Passw0rd!", "Todo User")

	// CREATE a todo
	rec := app.DoRequest("POST", "/api/v1/todos", token, map[string]any{
		"title":    "Buy groceries",
		"priority": "medium",
	})
	require.Equal(t, http.StatusCreated, rec.Code)

	resp := ParseResponse(t, rec)
	todoData := resp["data"].(map[string]any)
	todoID := todoData["id"].(string)
	assert.Equal(t, "Buy groceries", todoData["title"])
	assert.Equal(t, false, todoData["completed"])

	// LIST todos
	rec = app.DoRequest("GET", "/api/v1/todos", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	listData := resp["data"].(map[string]any)
	todos := listData["todos"].([]any)
	assert.Len(t, todos, 1)
	assert.Equal(t, float64(1), listData["total"])

	// GET by ID
	rec = app.DoRequest("GET", "/api/v1/todos/"+todoID, token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	todoData = resp["data"].(map[string]any)
	assert.Equal(t, "Buy groceries", todoData["title"])

	// UPDATE
	newTitle := "Buy organic groceries"
	rec = app.DoRequest("PUT", "/api/v1/todos/"+todoID, token, map[string]any{
		"title": newTitle,
	})
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	todoData = resp["data"].(map[string]any)
	assert.Equal(t, newTitle, todoData["title"])

	// TOGGLE complete
	rec = app.DoRequest("PATCH", "/api/v1/todos/"+todoID+"/toggle", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	todoData = resp["data"].(map[string]any)
	assert.Equal(t, true, todoData["completed"])

	// TOGGLE back to incomplete
	rec = app.DoRequest("PATCH", "/api/v1/todos/"+todoID+"/toggle", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	todoData = resp["data"].(map[string]any)
	assert.Equal(t, false, todoData["completed"])

	// DELETE
	rec = app.DoRequest("DELETE", "/api/v1/todos/"+todoID, token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	// Verify deleted — should be 404
	rec = app.DoRequest("GET", "/api/v1/todos/"+todoID, token, nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTodoInvalidID(t *testing.T) {
	app := NewTestApp()
	token := app.RegisterUser(t, "badid", "badid@example.com", "Passw0rd!", "Bad ID User")

	rec := app.DoRequest("GET", "/api/v1/todos/not-a-uuid", token, nil)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTodoAccessControl(t *testing.T) {
	app := NewTestApp()

	tokenA := app.RegisterUser(t, "userA", "usera@example.com", "Passw0rd!", "User A")
	tokenB := app.RegisterUser(t, "userB", "userb@example.com", "Passw0rd!", "User B")

	// User A creates a todo
	rec := app.DoRequest("POST", "/api/v1/todos", tokenA, map[string]any{
		"title":    "A's private todo",
		"priority": "high",
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	resp := ParseResponse(t, rec)
	todoID := resp["data"].(map[string]any)["id"].(string)

	// User B tries to access User A's todo → should be 403
	rec = app.DoRequest("GET", "/api/v1/todos/"+todoID, tokenB, nil)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	// User B tries to delete User A's todo → should be 403
	rec = app.DoRequest("DELETE", "/api/v1/todos/"+todoID, tokenB, nil)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	// User B's own todo list should be empty
	rec = app.DoRequest("GET", "/api/v1/todos", tokenB, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	resp = ParseResponse(t, rec)
	listData := resp["data"].(map[string]any)
	assert.Equal(t, float64(0), listData["total"])
}

func TestCreateTodoWithDescription(t *testing.T) {
	app := NewTestApp()
	token := app.RegisterUser(t, "descuser", "desc@example.com", "Passw0rd!", "Desc User")

	rec := app.DoRequest("POST", "/api/v1/todos", token, map[string]any{
		"title":       "Learn Go",
		"description": "Study concurrency patterns",
		"priority":    "high",
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	resp := ParseResponse(t, rec)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "Study concurrency patterns", data["description"])
}

// ==========================================================================
// List Operations Tests
// ==========================================================================

func TestListCRUDFlow(t *testing.T) {
	app := NewTestApp()
	token := app.RegisterUser(t, "listuser", "list@example.com", "Passw0rd!", "List User")

	// CREATE list
	rec := app.DoRequest("POST", "/api/v1/lists", token, map[string]string{
		"name": "Work Tasks",
	})
	require.Equal(t, http.StatusCreated, rec.Code)

	resp := ParseResponse(t, rec)
	listData := resp["data"].(map[string]any)
	listID := listData["id"].(string)
	assert.Equal(t, "Work Tasks", listData["name"])

	// LIST all lists
	rec = app.DoRequest("GET", "/api/v1/lists", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	listsData := resp["data"].(map[string]any)
	lists := listsData["lists"].([]any)
	assert.Len(t, lists, 1)

	// GET by ID (should include todos array)
	rec = app.DoRequest("GET", "/api/v1/lists/"+listID, token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	listDetail := resp["data"].(map[string]any)
	assert.Equal(t, "Work Tasks", listDetail["name"])
	assert.NotNil(t, listDetail["todos"])

	// RENAME
	rec = app.DoRequest("PUT", "/api/v1/lists/"+listID, token, map[string]string{
		"name": "Personal Tasks",
	})
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	listData = resp["data"].(map[string]any)
	assert.Equal(t, "Personal Tasks", listData["name"])

	// DELETE
	rec = app.DoRequest("DELETE", "/api/v1/lists/"+listID, token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	// Verify deleted — should be 404
	rec = app.DoRequest("GET", "/api/v1/lists/"+listID, token, nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListWithTodos(t *testing.T) {
	app := NewTestApp()
	token := app.RegisterUser(t, "listtodos", "lt@example.com", "Passw0rd!", "List Todos User")

	// Create a list
	rec := app.DoRequest("POST", "/api/v1/lists", token, map[string]string{
		"name": "Shopping",
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	resp := ParseResponse(t, rec)
	listID := resp["data"].(map[string]any)["id"].(string)

	// Create todos in that list
	for _, title := range []string{"Milk", "Bread", "Eggs"} {
		rec = app.DoRequest("POST", "/api/v1/todos", token, map[string]any{
			"title":    title,
			"priority": "low",
			"list_id":  listID,
		})
		require.Equal(t, http.StatusCreated, rec.Code)
	}

	// GET list by ID — should include 3 todos
	rec = app.DoRequest("GET", "/api/v1/lists/"+listID, token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	listDetail := resp["data"].(map[string]any)
	todos := listDetail["todos"].([]any)
	assert.Len(t, todos, 3)
}

func TestListAccessControl(t *testing.T) {
	app := NewTestApp()

	tokenA := app.RegisterUser(t, "listOwner", "lo@example.com", "Passw0rd!", "List Owner")
	tokenB := app.RegisterUser(t, "listIntruder", "li@example.com", "Passw0rd!", "Intruder")

	// User A creates a list
	rec := app.DoRequest("POST", "/api/v1/lists", tokenA, map[string]string{
		"name": "Secret List",
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	resp := ParseResponse(t, rec)
	listID := resp["data"].(map[string]any)["id"].(string)

	// User B tries to access User A's list → 403
	rec = app.DoRequest("GET", "/api/v1/lists/"+listID, tokenB, nil)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	// User B tries to delete → 403
	rec = app.DoRequest("DELETE", "/api/v1/lists/"+listID, tokenB, nil)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteListCascadesTodos(t *testing.T) {
	app := NewTestApp()
	token := app.RegisterUser(t, "cascadeuser", "cascade@example.com", "Passw0rd!", "Cascade User")

	// Create list
	rec := app.DoRequest("POST", "/api/v1/lists", token, map[string]string{
		"name": "Temp List",
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	resp := ParseResponse(t, rec)
	listID := resp["data"].(map[string]any)["id"].(string)

	// Create a todo in the list
	rec = app.DoRequest("POST", "/api/v1/todos", token, map[string]any{
		"title":    "Temp Todo",
		"priority": "low",
		"list_id":  listID,
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	resp = ParseResponse(t, rec)
	todoID := resp["data"].(map[string]any)["id"].(string)

	// Delete the list
	rec = app.DoRequest("DELETE", "/api/v1/lists/"+listID, token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	// The todo should also be gone (cascade delete)
	rec = app.DoRequest("GET", "/api/v1/todos/"+todoID, token, nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ==========================================================================
// Share Flow Tests
// ==========================================================================

func TestShareAndImportFlow(t *testing.T) {
	app := NewTestApp()

	tokenA := app.RegisterUser(t, "sharer", "sharer@example.com", "Passw0rd!", "Sharer")
	tokenB := app.RegisterUser(t, "importer", "importer@example.com", "Passw0rd!", "Importer")

	// User A creates a list with todos
	rec := app.DoRequest("POST", "/api/v1/lists", tokenA, map[string]string{
		"name": "Recipes",
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	resp := ParseResponse(t, rec)
	listID := resp["data"].(map[string]any)["id"].(string)

	for _, title := range []string{"Pasta", "Salad", "Soup"} {
		rec = app.DoRequest("POST", "/api/v1/todos", tokenA, map[string]any{
			"title":    title,
			"priority": "medium",
			"list_id":  listID,
		})
		require.Equal(t, http.StatusCreated, rec.Code)
	}

	// User A generates a share link
	rec = app.DoRequest("POST", "/api/v1/lists/"+listID+"/share", tokenA, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	shareData := resp["data"].(map[string]any)
	shareToken := shareData["share_token"].(string)
	assert.Len(t, shareToken, 64, "share token should be 64 chars")

	// User B imports the shared list
	rec = app.DoRequest("POST", "/api/v1/lists/import/"+shareToken, tokenB, nil)
	require.Equal(t, http.StatusCreated, rec.Code)

	resp = ParseResponse(t, rec)
	importedList := resp["data"].(map[string]any)
	assert.Contains(t, importedList["name"], "shared")

	importedTodos := importedList["todos"].([]any)
	assert.Len(t, importedTodos, 3, "imported list should have 3 todos")

	// Verify the imported list appears in User B's lists
	rec = app.DoRequest("GET", "/api/v1/lists", tokenB, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	resp = ParseResponse(t, rec)
	listsData := resp["data"].(map[string]any)
	assert.Equal(t, float64(1), listsData["total"])
}

func TestCannotImportOwnList(t *testing.T) {
	app := NewTestApp()

	token := app.RegisterUser(t, "selfshare", "self@example.com", "Passw0rd!", "Self Sharer")

	// Create list
	rec := app.DoRequest("POST", "/api/v1/lists", token, map[string]string{
		"name": "My List",
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	resp := ParseResponse(t, rec)
	listID := resp["data"].(map[string]any)["id"].(string)

	// Generate share link
	rec = app.DoRequest("POST", "/api/v1/lists/"+listID+"/share", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	resp = ParseResponse(t, rec)
	shareToken := resp["data"].(map[string]any)["share_token"].(string)

	// Try to import own list → should fail with 400
	rec = app.DoRequest("POST", "/api/v1/lists/import/"+shareToken, token, nil)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestImportWithInvalidToken(t *testing.T) {
	app := NewTestApp()

	token := app.RegisterUser(t, "badtoken", "bt@example.com", "Passw0rd!", "Bad Token")

	rec := app.DoRequest("POST", "/api/v1/lists/import/invalidtoken", token, nil)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ==========================================================================
// Duplicate List Tests
// ==========================================================================

func TestDuplicateList(t *testing.T) {
	app := NewTestApp()
	token := app.RegisterUser(t, "dupuser", "dup@example.com", "Passw0rd!", "Dup User")

	// Create list with todos
	rec := app.DoRequest("POST", "/api/v1/lists", token, map[string]string{
		"name": "Original",
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	resp := ParseResponse(t, rec)
	listID := resp["data"].(map[string]any)["id"].(string)

	rec = app.DoRequest("POST", "/api/v1/todos", token, map[string]any{
		"title":    "Task 1",
		"priority": "high",
		"list_id":  listID,
	})
	require.Equal(t, http.StatusCreated, rec.Code)

	// Duplicate the list
	rec = app.DoRequest("POST", "/api/v1/lists/"+listID+"/duplicate", token, nil)
	require.Equal(t, http.StatusCreated, rec.Code)

	resp = ParseResponse(t, rec)
	dupList := resp["data"].(map[string]any)
	assert.Contains(t, dupList["name"], "Copy")
	dupTodos := dupList["todos"].([]any)
	assert.Len(t, dupTodos, 1)

	// Both original and duplicate should exist
	rec = app.DoRequest("GET", "/api/v1/lists", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	resp = ParseResponse(t, rec)
	assert.Equal(t, float64(2), resp["data"].(map[string]any)["total"])
}

// ==========================================================================
// Pagination Tests
// ==========================================================================

func TestTodoPagination(t *testing.T) {
	app := NewTestApp()
	token := app.RegisterUser(t, "pageuser", "page@example.com", "Passw0rd!", "Page User")

	// Create 15 todos
	for i := 1; i <= 15; i++ {
		rec := app.DoRequest("POST", "/api/v1/todos", token, map[string]any{
			"title":    fmt.Sprintf("Todo %d", i),
			"priority": "low",
		})
		require.Equal(t, http.StatusCreated, rec.Code)
	}

	// Page 1 (default page_size=10)
	rec := app.DoRequest("GET", "/api/v1/todos?page=1&page_size=10", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp := ParseResponse(t, rec)
	data := resp["data"].(map[string]any)
	assert.Len(t, data["todos"].([]any), 10)
	assert.Equal(t, float64(15), data["total"])
	assert.Equal(t, float64(1), data["page"])
	assert.Equal(t, float64(10), data["page_size"])
	assert.Equal(t, float64(2), data["total_pages"])

	// Page 2 (should have 5)
	rec = app.DoRequest("GET", "/api/v1/todos?page=2&page_size=10", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	resp = ParseResponse(t, rec)
	data = resp["data"].(map[string]any)
	assert.Len(t, data["todos"].([]any), 5)
	assert.Equal(t, float64(2), data["page"])
}

// ==========================================================================
// Move Todos Tests
// ==========================================================================

func TestMoveTodosBetweenLists(t *testing.T) {
	app := NewTestApp()
	token := app.RegisterUser(t, "moveuser", "move@example.com", "Passw0rd!", "Move User")

	// Create two lists
	rec := app.DoRequest("POST", "/api/v1/lists", token, map[string]string{"name": "List A"})
	require.Equal(t, http.StatusCreated, rec.Code)
	listAID := ParseResponse(t, rec)["data"].(map[string]any)["id"].(string)

	rec = app.DoRequest("POST", "/api/v1/lists", token, map[string]string{"name": "List B"})
	require.Equal(t, http.StatusCreated, rec.Code)
	listBID := ParseResponse(t, rec)["data"].(map[string]any)["id"].(string)

	// Create todos in List A
	var todoIDs []string
	for _, title := range []string{"Move 1", "Move 2"} {
		rec = app.DoRequest("POST", "/api/v1/todos", token, map[string]any{
			"title":    title,
			"priority": "medium",
			"list_id":  listAID,
		})
		require.Equal(t, http.StatusCreated, rec.Code)
		todoIDs = append(todoIDs, ParseResponse(t, rec)["data"].(map[string]any)["id"].(string))
	}

	// Move todos from List A to List B
	rec = app.DoRequest("PATCH", "/api/v1/todos/move", token, map[string]any{
		"todo_ids": todoIDs,
		"list_id":  listBID,
	})
	require.Equal(t, http.StatusOK, rec.Code)

	// List A should have 0 todos
	rec = app.DoRequest("GET", "/api/v1/lists/"+listAID, token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	resp := ParseResponse(t, rec)
	todosA := resp["data"].(map[string]any)["todos"].([]any)
	assert.Len(t, todosA, 0)

	// List B should have 2 todos
	rec = app.DoRequest("GET", "/api/v1/lists/"+listBID, token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	resp = ParseResponse(t, rec)
	todosB := resp["data"].(map[string]any)["todos"].([]any)
	assert.Len(t, todosB, 2)
}

// ==========================================================================
// Health Check
// ==========================================================================

func TestHealthEndpoint(t *testing.T) {
	app := NewTestApp()

	rec := app.DoRequest("GET", "/health", "", nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	resp := ParseResponse(t, rec)
	assert.Equal(t, "ok", resp["status"])
}
