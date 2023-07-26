package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID string `json:"id"`
    Firstname  string `json:"firstname"` 
    Lastname string `json:"lastname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
    db, err := sql.Open("mysql", "root:zadicus_123@tcp(127.0.0.1:3306)/authgo")
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
 
    r := gin.Default()
	authGroup := r.Group("/")
	userGroup := r.Group("/user")
    authGroup.POST("/signup", signupHandler(db))
    authGroup.POST("/login", loginHandler(db))

	userGroup.Use(AuthMiddleware())
	userGroup.GET("/profile", profileHandler(db))
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}

func signupHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user User
		err := c.ShouldBindJSON(&user)
		if err != nil {
            c.JSON(400, gin.H{
                "error": "Bad Request",
            })
            return
        }
        rows, err := db.Query("SELECT username FROM users where username = ?", user.Username)
        if err != nil {
            c.JSON(500, gin.H{
                "error": "Server Error",
            })
            return
        }
        defer rows.Close()
		if rows.Next() {
			c.JSON(400, gin.H{
				"error": "username already exists",
			})
			return
		}
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Server Error",
			})
			return
		}
		log.Println(hashedPassword)
		user.Password = hashedPassword
		rows, err = db.Query("INSERT INTO users (firstname, lastname, username, password) VALUES (?, ?, ?, ?)", user.Firstname, user.Lastname, user.Username, hashedPassword)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Database error",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.Password)
		}
		c.JSON(200, gin.H{
			"message": "User created successfully",
		})
    }
}

func loginHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
		}
		rows, err := db.Query("SELECT password FROM users WHERE username = ?", user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
			return
		}
		defer rows.Close()
		if rows.Next(){
			hashedPassword := ""
			err := rows.Scan(&hashedPassword)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Server error",
				})
				return
			}
			if !VerifyPassword(user.Password, hashedPassword) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid username or password",
				})
				return
			}
			tokenString, err := CreateToken(user.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error creating token",
				})
				return
			}
			rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.Password)
			for rows.Next() {
				rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.Password)
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful",
				"token": tokenString,
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
    }
}

func profileHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		rows, err := db.Query("SELECT id, firstname, lastname, username FROM users WHERE id = ?", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
			return
		}
		defer rows.Close()
		var user User
		if rows.Next() {
			rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Username)
			c.JSON(http.StatusOK, gin.H{
				"user": user,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error",
		})
	}}