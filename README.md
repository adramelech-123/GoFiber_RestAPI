# GO-FIBER API

An API built using GOFiber which is an Express inspired web framework built on top of Fasthttp, the fastest HTTP engine for Go.

We also use `GORM` which is a full-featured ORM for SQL and SQLite.

## 1. Installations

```bash
# Initialize Go
go mod init github.com/adramelech-123/fiber-api

# Install GORM with the SQLite driver
go get -u "gorm.io/driver/sqlite"
go get -u gorm.io/gorm

# Install Fiber
go get github.com/gofiber/fiber/v2

# Install Air for live-reloading
go install github.com/air-verse/air@latest
air init

```

We have to do additional installations and setup in order to be able to work with SQLite. First we install `MinGW-w64` since `CGO` requires a C compiler in order to run. Golang CGO is a package in Go that enables developers to interact with C code in their Go programs. With CGO, developers can call C functions, use C variables and pointers, and pass data structures between Go and C. CGO is required to run sqlite.

1. Visit [Sourceforge](https://sourceforge.net/projects/mingw-w64/files/) to download the second file located under `MinGW-W64 GCC-8.1.0`. The file can be found [here](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z)
2. Unzip the file and move the unzipped `mingw64` folder to the C Drive.
3. Navigate to `C:\mingw64\bin` and copy this path
4. Open environment variables and in the System variables section, find the Path variable and select it.
5. Click on Edit and add the path to the MinGW-w64 bin directory (e.g., C:\mingw64\bin).
6. Click OK to close all dialogs.
7. Check if `gcc` is prperly installed by opening your vscode terminal and running `gcc --version`
8. Enable CGO by running `$env:CGO_ENABLED=1` in the terminal

Now we can install SQLite by downloading it from the SQLite website and adding the .exe file to the environment variables.

## 2. Setup basic server

Now we need to setup a basic scaffold for our server in our `main.go` file as follows:

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Ctx represents the Context which hold the HTTP request and response.
// It has methods for the request query string, parameters, body, HTTP headers and so on.
func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to this awesome Go based API")
}

func main() {
	app := fiber.New()

	app.Get("/api", welcome)

	log.Fatal(app.Listen(":3000"))
}
```

Now we can run the `air` hot reloader by simply using the command:

```bash
air
```

## 3. Connect Database

Create a new directory called `database` and add a file called `database.go`.

```go
// This package connects the application to the database and works with the ORM to setup an instance
package database

import (
	"log"
	"os"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The database instance will point to gorm.DB struct and methods
type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	// Variable db will open our sqlite database and provide the initial config
	db, err := gorm.Open(sqlite.Open("fiberapi.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(2)
	}

	log.Println("Database connection successful! 😁")

	// A logger to log into our database
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running migrations...")
	// TODO: Add migrations

	Database = DbInstance{Db: db}

}
```

## 4.Models

Now we create a directory for our `models` that will contain three types of tables/data we want to model, namely `User`, `Product` and `Order` in three separate files as follows:

```go
// User Model in user.go file
type User struct {
	ID              uint `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
}

// Product Model in product.go file
type Product struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time
	ProductName  string `json:"product_name"`
	SerialNumber string `json:"serial_number"`

}

// Order Model in order.go file
type Order struct {
	ID             uint `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time
	ProductRefer   int `json:"product_id"`
	Product        Product `gorm:"foreignKey:ProductRefer"`
	UserRefer      int `json:"user_id"`
	User           User `gorm:"foreignKey:UserRefer"`
}
```

### Note

- `json:"id"`: This specifies that when this struct is serialized to JSON, the ID field will be represented as id.
- `gorm:"primaryKey"`: This indicates to GORM that this field is the primary key of the table.
- `gorm:"foreignKey:ProductRefer"`: This indicates to GORM that this field is a foreign key relationship. The foreign key in the Order table is ProductRefer, which refers to the primary key in the Product table.
- `gorm:"foreignKey:UserRefer"`: This indicates to GORM that this field is a foreign key relationship. The foreign key in the Order table is UserRefer, which refers to the primary key in the User table.

## 5. Add Migrations to SQLite DB

Now we can add migrations to the SQLite Database in our `database.go` file.

```go
func ConnectDb() {
	// Variable db will open our sqlite database and provide the initial config
	db, err := gorm.Open(sqlite.Open("fiberapi.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(2)
	}

	log.Println("Database connection successful! 😁")

	// A logger to log into our database
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running migrations...")

	// Add migrations
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})

	Database = DbInstance{Db: db}
}
```

We should also call the `ConnectDb` in our main.go file:

```go
func main() {
	database.ConnectDb()
	app := fiber.New()

	app.Get("/api", welcome)
	log.Fatal(app.Listen(":3000"))
}
```

## 6. Creating Endpoints

### User Endpoints

**Create User**

```go
package routes

import (
	"github.com/adramelech-123/fiber-api/database"
	"github.com/adramelech-123/fiber-api/models"
	"github.com/gofiber/fiber/v2"
)


// User struct specifically for serialization (a data transfer object).
// It differs from the models.User struct, which likely represents the database model.
type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

//  Helper function to create user endpoints
func CreateResponseUser(userModel models.User) User {
	return User{
		ID: userModel.ID,
		FirstName: userModel.FirstName,
		LastName: userModel.LastName,
	}
}

// Create User Endpoint

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

// Get Users

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	// Find all records in the User model
	database.Database.Db.Find(&users)
	
	responseUsers := []User{}

	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}
```
### Key Points

- `c.BodyParser(&user)` parses the JSON body of the request into the user variable. If parsing fails, it returns a 400 status with the error message.
- `database.Database.Db.Create(&user)` adds the new user record to the database using the GORM ORM.
- `responseUser := CreateResponseUser(user)` converts the newly created user model to the serialized User struct.
- `return c.Status(200).JSON(responseUser)` sends a 200 status response with the serialized user data.

In the `main.go` file we will create a function to setup our all our endpoints as follows:

```go
func setupRoutes(app*fiber.App) {
	// Welcome endpoint
	app.Get("/api", welcome)

	// All User endpoints here
	app.Post("/api/users", routes.CreateUser)
}

```

Now the `main` function will simply look as follows:

```go
func main() {
	database.ConnectDb()
	app := fiber.New()
	setupRoutes(app) 
	log.Fatal(app.Listen(":3000"))
}
```
We no longer have to keep adding routes to the `main` function.

