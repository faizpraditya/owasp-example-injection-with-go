package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// go get github.com/gin-gonic/gin
// go get github.com/jackc/pgx
// go get github.com/jmoiron/sqlx

// create database db_sample_injection

// create table user_credentials (id serial primary key, user_name varchar(15), user_password varchar(100), is_blocked int);

type UserCredential struct {
	Id           uint   `db:"id"`
	UserName     string `db:"user_name"`
	UserPassword string `db:"user_password"`
	IsBlocked    int    `db:"is_blocked"`
}

type Login struct {
	User     string `json:"user_name" binding:"required"`
	Password string `json:"user_password" binding:"required"`
}

func main() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "postgres", "12345678", "localhost", "5432", "db_sample_injection")

	conn, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer func(conn *sqlx.DB) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	r := gin.Default()
	route := r.Group("/enigma")
	route.POST("/auth", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var userCredential = UserCredential{}
		// sql := fmt.Sprintf("SELECT * FROM user_credentials WHERE user_name = '%s' and user_password = '%s'", login.User, login.Password)
		// Handling sql injection
		sql := "SELECT * FROM user_credentials WHERE user_name = $1 and user_password = $2"
		fmt.Println(sql)

		err := conn.Get(&userCredential, sql, login.User, login.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	listenAddress := fmt.Sprintf("%s:%s", "localhost", "8888")
	err = r.Run(listenAddress)
	if err != nil {
		panic(err)
	}
}
