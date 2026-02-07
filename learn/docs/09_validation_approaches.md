# Validation Approaches in Go

Understanding different ways to validate data and why we use struct tags.

---

## The Problem

When users send data to your API, you need to validate it:
- Is the username present?
- Is the email valid?
- Is the password strong enough?
- Does the username have no spaces?

---

## Approach 1: Manual Validation (Direct Calls)

You **could** write validation logic directly in the handler:

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err.Error())
        return
    }

    // MANUAL VALIDATION - checking everything manually
    if req.Username == "" {
        c.JSON(400, gin.H{"error": "username is required"})
        return
    }
    if len(req.Username) < 3 || len(req.Username) > 30 {
        c.JSON(400, gin.H{"error": "username must be 3-30 characters"})
        return
    }
    if strings.Contains(req.Username, " ") {
        c.JSON(400, gin.H{"error": "username cannot contain spaces"})
        return
    }
    if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`).MatchString(req.Username) {
        c.JSON(400, gin.H{"error": "username must start with letter and contain only letters, numbers, underscores"})
        return
    }

    if req.Email == "" {
        c.JSON(400, gin.H{"error": "email is required"})
        return
    }
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(req.Email) {
        c.JSON(400, gin.H{"error": "invalid email format"})
        return
    }

    if req.Password == "" {
        c.JSON(400, gin.H{"error": "password is required"})
        return
    }
    if len(req.Password) < 8 {
        c.JSON(400, gin.H{"error": "password must be at least 8 characters"})
        return
    }
    // Check for uppercase, lowercase, number, special char...
    hasUpper := false
    hasLower := false
    hasNumber := false
    hasSpecial := false
    for _, char := range req.Password {
        if unicode.IsUpper(char) { hasUpper = true }
        if unicode.IsLower(char) { hasLower = true }
        if unicode.IsNumber(char) { hasNumber = true }
        if unicode.IsPunct(char) { hasSpecial = true }
    }
    if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
        c.JSON(400, gin.H{"error": "password must contain uppercase, lowercase, number, and special character"})
        return
    }

    // FINALLY, after 50+ lines of validation, do the actual work
    response, err := h.userService.Register(c.Request.Context(), req)
    if err != nil {
        utils.BadRequest(c, err.Error())
        return
    }

    utils.Created(c, response)
}
```

**Problems:**
- Handler is 50+ lines, mostly validation
- Same validation code repeated in tests
- Hard to maintain
- Business logic is buried
- Every endpoint needs this

---

## Approach 2: Declarative Validation (Struct Tags) ✅ BETTER

Use struct tags to declare validation rules:

```go
// DTO - DECLARES validation rules
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=30,alphanumunder,nospaces"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8,strongpassword"`
    FullName string `json:"full_name" binding:"required,min=2,max=100"`
}

// HANDLER - Clean and focused
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest

    // Gin automatically validates using the struct tags!
    if err := c.ShouldBindJSON(&req); err != nil {
        validationErrors := validator.GetValidationErrors(err)
        c.JSON(400, gin.H{"error": "Validation failed", "fields": validationErrors})
        return
    }

    // All validation passed! Do the actual work
    response, err := h.userService.Register(c.Request.Context(), req)
    if err != nil {
        utils.BadRequest(c, err.Error())
        return
    }

    utils.Created(c, response)
}
```

**Benefits:**
- Handler is ~15 lines instead of 50+
- Validation rules are visible at a glance
- Reusable across endpoints
- Easy to test
- Standard Go pattern

---

## How It Works Internally

### 1. Gin Uses go-playground/validator

When you call `c.ShouldBindJSON(&req)`, Gin internally:

```go
// Inside Gin's code (simplified)
func (c *Context) ShouldBindJSON(obj interface{}) error {
    // 1. Parse JSON into struct
    if err := json.Unmarshal(bodyBytes, obj); err != nil {
        return err
    }

    // 2. Validate using go-playground/validator
    validate := validator.New()
    if err := validate.Struct(obj); err != nil {
        return err  // Returns validation errors
    }

    return nil
}
```

### 2. Built-in vs Custom Validators

**Built-in validators** (already in go-playground/validator):
- `required` - field must be present
- `min=3` - minimum length
- `max=30` - maximum length
- `email` - valid email format

**Custom validators** (we need to add):
- `alphanumunder` - our custom rule
- `nospaces` - our custom rule
- `strongpassword` - our custom rule

### 3. Registration Process

```
                    main.go
                       ↓
        Get Gin's validator engine
                       ↓
        Register custom validators
        (alphanumunder, nospaces, etc.)
                       ↓
        Now Gin knows what to do when
        it sees these tags in structs
```

**In main.go:**
```go
func registerCustomValidators() {
    // Get Gin's validator instance
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        // Teach it our custom validators
        v.RegisterValidation("alphanumunder", alphaNumericUnderscore)
        v.RegisterValidation("nospaces", noSpaces)
        v.RegisterValidation("strongpassword", strongPassword)
    }
}
```

Now when Gin sees `binding:"alphanumunder"`, it knows to call our `alphaNumericUnderscore` function!

---

## Why Not a New Library?

**We're NOT adding a new library!**

Gin already depends on go-playground/validator:

```bash
$ go mod graph | grep validator
github.com/gin-gonic/gin@v1.9.1 github.com/go-playground/validator/v10@v10.14.0
```

We're just:
1. Getting a reference to the validator Gin already uses
2. Extending it with our custom rules

**Analogy:**
- Gin is like a car with an engine
- go-playground/validator is the engine
- We're adding a turbocharger (custom validators) to the existing engine

---

## Node.js Comparison

This is similar to how you'd use validation in Node.js:

### Express + Joi/Zod (Node.js)
```javascript
// Define schema
const registerSchema = Joi.object({
  username: Joi.string().min(3).max(30).pattern(/^[a-zA-Z][a-zA-Z0-9_]*$/).required(),
  email: Joi.string().email().required(),
  password: Joi.string().min(8).pattern(/strong password regex/).required()
});

// Use in handler
app.post('/register', (req, res) => {
  // Joi validates against schema
  const { error, value } = registerSchema.validate(req.body);
  if (error) {
    return res.status(400).json({ error: error.details });
  }

  // Validation passed, do work
  const user = await userService.register(value);
  res.json(user);
});
```

### Go + Gin + validator (Our approach)
```go
// Define schema via struct tags
type RegisterRequest struct {
    Username string `binding:"required,min=3,max=30,alphanumunder"`
    Email    string `binding:"required,email"`
    Password string `binding:"required,min=8,strongpassword"`
}

// Use in handler
func Register(c *gin.Context) {
    var req RegisterRequest

    // Gin validates against tags
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": validator.GetValidationErrors(err)})
        return
    }

    // Validation passed, do work
    user, err := userService.Register(req)
    c.JSON(200, user)
}
```

---

## Summary

| Aspect | Manual Validation | Struct Tags (Our Approach) |
|--------|-------------------|---------------------------|
| Code length | 50+ lines per handler | ~15 lines per handler |
| Reusability | Copy-paste everywhere | Define once, use anywhere |
| Maintainability | Hard to change | Update struct tag |
| Readability | Validation scattered | Rules visible at top |
| Testing | Must test in handler | Can test validator separately |
| New library? | No | No (Gin already uses it!) |

---

## Why Register in main.go?

**Short answer:** So Gin knows what to do when it sees your custom tags.

**Long answer:**
1. Gin uses go-playground/validator internally
2. Validator has built-in rules (`required`, `email`, etc.)
3. Our custom rules (`alphanumunder`, `nospaces`) need to be registered
4. We register ONCE at app startup (main.go)
5. Now ALL handlers can use these tags automatically

**Without registration:**
```go
type RegisterRequest struct {
    Username string `binding:"alphanumunder"` // ❌ ERROR: unknown tag
}
```

**With registration:**
```go
type RegisterRequest struct {
    Username string `binding:"alphanumunder"` // ✅ Works! Validator knows this rule
}
```

---

## Could We Skip Registration and Call Directly?

Yes, but you'd lose ALL the benefits:

```go
// Without struct tags (manual calls)
func Register(c *gin.Context) {
    var req RegisterRequest
    c.ShouldBindJSON(&req)

    // Manual validation calls
    if !validator.IsAlphaNumUnderscore(req.Username) {
        c.JSON(400, gin.H{"error": "invalid username"})
        return
    }
    if !validator.NoSpaces(req.Username) {
        c.JSON(400, gin.H{"error": "username has spaces"})
        return
    }
    if !validator.IsStrongPassword(req.Password) {
        c.JSON(400, gin.H{"error": "weak password"})
        return
    }

    // Now do actual work...
}
```

This defeats the purpose! You're back to manual validation.

---

## Conclusion

**Registration in main.go** = One-time setup that enables **declarative validation** everywhere.

Think of it like:
- **Node.js:** `app.use(express.json())` - setup middleware once
- **Go:** `registerCustomValidators()` - setup validators once

Both are initialization steps that make your handlers cleaner!
