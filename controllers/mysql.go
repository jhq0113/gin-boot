package controllers

import (
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gin-boot/global"
)

type Mysql struct {
	Controller
}

type User struct {
	Id       int64  `json:"id"`
	UserName string `json:"user_name"`
	AddTime  int64  `json:"add_time"`
}

func (m *Mysql) Insert(ctx *gin.Context) {
	pool := global.BootMysqlGroup.SelectPool(true)
	result, err := pool.Insert("user", map[string]interface{}{
		"user_name": rand.Intn(math.MaxInt32),
		"add_time":  time.Now().Unix(),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
		return
	}

	userId, _ := result.LastInsertId()
	ctx.JSON(http.StatusOK, map[string]int64{
		"userId": userId,
	})
}

func (m *Mysql) BatchInsert(ctx *gin.Context) {
	pool := global.BootMysqlGroup.SelectPool(true)
	result, err := pool.BatchInsert("user", []map[string]interface{}{
		{
			"user_name": rand.Intn(math.MaxInt32),
			"add_time":  time.Now().Unix(),
		},
		{
			"user_name": rand.Intn(math.MaxInt32),
			"add_time":  time.Now().Unix(),
		},
		{
			"user_name": rand.Intn(math.MaxInt32),
			"add_time":  time.Now().Unix(),
		},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
		return
	}

	rows, _ := result.RowsAffected()
	ctx.JSON(http.StatusOK, map[string]int64{
		"rows": rows,
	})
}

func (m *Mysql) One(ctx *gin.Context) {
	pool := global.BootMysqlGroup.SelectPool(false)
	query := pool.Find().
		From("`user`").
		Select("`id`", "`user_name`", "`add_time`").
		Order("`id` DESC")
	row := pool.One(query)

	user := User{}
	err := row.Scan(&user.Id, &user.UserName, &user.AddTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (m *Mysql) Some(ctx *gin.Context) {
	pool := global.BootMysqlGroup.SelectPool(false)
	query := pool.Find().
		From("`user`").
		Select("`id`", "`user_name`", "`add_time`").
		Order("`id` DESC").
		Limit(0, 20)

	rows, _ := pool.All(query)
	userList := make([]User, 0, 20)

	user := User{}
	for rows.Next() {
		_ = rows.Scan(&user.Id, &user.UserName, &user.AddTime)
		userList = append(userList, user)
	}

	ctx.JSON(http.StatusOK, userList)
}
