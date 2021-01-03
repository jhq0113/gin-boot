package controllers

import (
	"gin-boot/mysql"
	"gin-boot/utils"
	"math"
	"math/rand"
	"net/http"
	"strconv"
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

func (m *Mysql) Delete(ctx *gin.Context) {
	id := QueryInt64(ctx, "id", 0)
	if id < 1 {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"msg": "参数id必须大于0",
		})
		return
	}

	pool := global.BootMysqlGroup.SelectPool(true)
	condition := map[string]interface{}{
		"id": id,
	}

	exists, err := pool.Exists("`user`", condition)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"msg": err.Error(),
		})
		return
	}

	if exists {
		result, err := pool.DeleteAll("user", map[string]interface{}{
			"id": id,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"msg": err.Error(),
			})
			return
		}

		if rows, _ := result.RowsAffected(); rows != 1 {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"msg": "del failed",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"msg": "ok",
	})
}

func (m *Mysql) One(ctx *gin.Context) {
	id := QueryInt64(ctx, "id", 0)
	if id < 1 {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"msg": "参数id必须大于0",
		})
		return
	}

	pool := global.BootMysqlGroup.SelectPool(false)
	query := pool.Find().
		From("`user`").
		Select("`id`", "`user_name`", "`add_time`").
		Where(map[string]interface{}{
			"id": id,
		})
	row := pool.One(query)

	user := User{}
	err := row.Scan(&user.Id, &user.UserName, &user.AddTime)
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"msg": "id为" + strconv.FormatInt(id, 10) + "的记录不存在",
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (m *Mysql) Some(ctx *gin.Context) {
	page := QueryInt64(ctx, "page", 1)
	if page < 1 {
		page = 1
	}

	pageSize := QueryInt64(ctx, "pageSize", 15)
	if pageSize < 2 {
		pageSize = 2
	}

	if pageSize > 100 {
		pageSize = 100
	}

	pool := global.BootMysqlGroup.SelectPool(false)
	query := pool.Find().
		From("`user`").
		Select("`id`", "`user_name`", "`add_time`").
		Order("`id` DESC").
		Limit((page-1)*pageSize, pageSize)

	rows, _ := pool.All(query)

	userList := []User{}
	mysql.FormatRows(rows, func(fieldValue map[string][]byte) {
		userList = append(userList, User{
			Id:       utils.Bytes2Int64(fieldValue["id"]),
			UserName: string(fieldValue["user_name"]),
			AddTime:  utils.Bytes2Int64(fieldValue["add_time"]),
		})
	})

	ctx.JSON(http.StatusOK, userList)
}
