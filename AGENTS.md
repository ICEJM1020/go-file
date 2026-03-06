# AGENTS.md

## Environment Setup

**Important**: The local Windows environment does not have Golang installed. All Go-related commands (build, test, run) must be executed in WSL (Windows Subsystem for Linux).

- **WSL Distribution**: Ubuntu-24.04
- **Go Version**: golang-1.26 (pre-installed in WSL)

### Running Commands in WSL

```bash
# Run commands directly from Windows
wsl -d Ubuntu-24.04 -- go build -o go-file
wsl -d Ubuntu-24.04 -- go test ./...
wsl -d Ubuntu-24.04 -- go run main.go
```

### Path Mapping

- Windows path: `E:\FunnySpace\go-file`
- WSL path: `/mnt/e/FunnySpace/go-file`

## Build Commands

```bash
# Build the executable
go build -o go-file

# Build with version info (from Git tags)
go build -ldflags "-s -w -X 'go-file/common.Version=$(git describe --tags)' -extldflags '-static'" -o go-file

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o go-file-linux
GOOS=darwin GOARCH=amd64 go build -o go-file-mac
GOOS=windows GOARCH=amd64 go build -o go-file.exe
# GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o go-file.exe

# Run directly
go run main.go

# Run with options
./go-file --port 3000 --path ./share
./go-file --host 192.168.1.100 --video ./videos
./go-file --no-browser
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run tests in specific package
go test ./model
go test ./controller

# Run a single test
go test -run TestFunctionName ./model

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

## Code Style Guidelines

### Import Ordering

Imports are organized in three groups separated by blank lines:
1. Standard library imports
2. Third-party imports (github.com/*, golang.org/*, etc.)
3. Internal module imports (go-file/*)

```go
import (
    "fmt"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"

    "go-file/common"
    "go-file/model"
)
```

### Naming Conventions

- **Packages**: lowercase, single word (`common`, `model`, `controller`, `middleware`, `router`)
- **Structs**: PascalCase (`File`, `User`, `FileDeleteRequest`)
- **Methods**: PascalCase for exported, camelCase for internal
- **Variables**: camelCase
- **Constants**: PascalCase (`RoleAdminUser`, `UserStatusEnabled`)
- **Functions**: PascalCase for exported, camelCase for internal

### Error Handling

```go
// Check errors immediately and handle
db, err := model.InitDB()
if err != nil {
    common.FatalLog(err)
}

// Return errors from functions
func (file *File) Insert() error {
    var err error
    err = DB.Create(file).Error
    return err
}

// Handle errors in HTTP handlers
form, err := c.MultipartForm()
if err != nil {
    c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
    return
}
```

### Logging

Use the custom logging functions from `common` package:
- `common.SysLog(message)` - System information
- `common.SysError(message)` - System errors
- `common.FatalLog(args...)` - Fatal errors (exits)
- `common.P2pLog(message)` - P2P logging

```go
common.SysLog("Server started")
common.SysError("failed to get IP: " + err.Error())
common.FatalLog("failed to initialize database")
```

### Struct Definitions with GORM Tags

```go
type File struct {
    Id              int    `json:"id"`
    Filename        string `json:"filename"`
    Description     string `json:"description"`
    Uploader        string `json:"uploader"`
    Link            string `json:"link" gorm:"unique"`
    Time            string `json:"time"`
    DownloadCounter int    `json:"download_counter"`
}
```

### HTTP Handler Pattern

```go
func HandlerName(c *gin.Context) {
    // Extract data from context
    username := c.GetString("username")

    // Process request
    data, err := someOperation()
    if err != nil {
        common.SysError(err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "操作失败",
        })
        return
    }

    // Return response
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    data,
    })
}
```

### Middleware Pattern

```go
func MiddlewareName() func(c *gin.Context) {
    return func(c *gin.Context) {
        // Pre-processing
        session := sessions.Default(c)
        username := session.Get("username")

        if username == nil {
            c.JSON(http.StatusForbidden, gin.H{
                "success": false,
                "message": "无权进行此操作",
            })
            c.Abort()
            return
        }

        // Set context values
        c.Set("username", username)
        c.Next()
    }
}
```

### Database Query Pattern

```go
// Query all
var files []*File
err = DB.Find(&files).Error

// Query with conditions
err = DB.Limit(common.ItemsPerPage).Offset(startIdx).
    Where("filename LIKE ?", "%"+query+"%").
    Order("id desc").
    Find(&files).Error

// Update
DB.Model(&File{}).Where("link = ?", link).
    UpdateColumn("download_counter", gorm.Expr("download_counter + 1"))
```

### Environment Variables

- `REDIS_CONN_STRING` - Redis connection string for rate limiting
- `SQL_DSN` - MySQL connection string (DSN format)
- `SQLITE_PATH` - Custom SQLite database path (default: `go-file.db`)
- `SESSION_SECRET` - Session encryption key
- `UPLOAD_PATH` - File upload directory (default: `./upload`)
- `GIN_MODE` - Gin mode (`debug` or `release`)

### Project Structure

```
go-file/
├── main.go              # Application entry point
├── common/              # Shared utilities and constants
│   ├── constants.go     # Global variables and flags
│   ├── logger.go        # Logging functions
│   ├── redis.go         # Redis client
│   ├── utils.go         # Helper functions
│   └── public/          # Embedded static files
├── controller/          # HTTP request handlers
├── middleware/          # Gin middleware (auth, rate-limit, etc.)
├── model/               # Database models and queries
├── router/              # Route definitions
└── upload/              # Default upload directory (ignored in git)
```

### File Operations

Always validate file paths to prevent directory traversal:

```go
fullPath := filepath.Join(common.UploadPath, subfolder, filename)
if !strings.HasPrefix(fullPath, common.UploadPath) {
    // We may being attacked!
    c.Status(403)
    return
}
```

### Constants

Define role and permission constants using `const` blocks:

```go
const (
    RoleGuestUser  = 0
    RoleCommonUser = 1
    RoleAdminUser  = 10
)

const (
    UserStatusEnabled  = 1
    UserStatusDisabled = 2 // don't use 0
)
```
