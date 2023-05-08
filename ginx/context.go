package ginx

import (
	"strconv"
	"strings"

	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/x"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func QueryBool(c *gin.Context, key string) bool {
	q := c.Query(key)
	if q == "" {
		return false
	}
	result, err := strconv.ParseBool(q)
	if err != nil {
		errx.PanicValidatition("解析参数: " + key + " 失败: " + err.Error())
	}
	return result
}

func QueryInt(ctx *gin.Context, key string, defaultValue int) int {
	qStr := ctx.Query(key)
	if qStr == "" {
		return defaultValue
	} else {
		v, err := strconv.Atoi(qStr)
		if err != nil {
			errx.PanicValidatition(key + " is not integer")
		}
		return v
	}
}
func QueryID(ctx *gin.Context, key string) *primitive.ObjectID {
	qStr := ctx.Query(key)
	if qStr == "" {
		return nil
	}
	result, err := primitive.ObjectIDFromHex(qStr)
	if err != nil {
		errx.PanicValidatition(err.Error())
	}
	return &result
}

func QueryIDRequired(ctx *gin.Context, key string) primitive.ObjectID {
	qStr := ctx.Query(key)
	if qStr == "" {
		errx.PanicValidatition(key + " is required")
	}
	result, err := primitive.ObjectIDFromHex(qStr)
	if err != nil {
		errx.PanicValidatition(err.Error())
	}
	return result
}

func QueryPage(ctx *gin.Context) (int, int) {
	page := QueryInt(ctx, "page", 1)
	if page < 1 {
		errx.PanicValidatition("page must be greater than 0")
	}
	pageSize := QueryInt(ctx, "pageSize", 10000)
	if pageSize < 1 {
		errx.PanicValidatition("pageSize must be greater than 0")
	}
	return page, pageSize
}

func QueryPageParams(ctx *gin.Context) *x.PageParam {
	page, pageSize := QueryPage(ctx)
	return &x.PageParam{
		Page:     int64(page),
		PageSize: int64(pageSize),
	}
}

func QueryPageAndSort(ctx *gin.Context) *x.PageAndSort {
	page, pageSize := QueryPage(ctx)
	sort := ctx.Query("sort")
	includeTotal := QueryBool(ctx, "includeTotal")
	p := x.NewPageAndSort(int64(page), int64(pageSize), sort)
	p.IncludeTotal = includeTotal
	return p
}

func QueryStrings(ctx *gin.Context, key string) []string {
	v := ctx.Query(key)
	if v == "" {
		return make([]string, 0)
	}
	return strings.Split(v, ",")
}

func Path(ctx *gin.Context, key string) string {
	v := ctx.Param(key)
	if v == "" {
		errx.PanicValidatition(key + " must not be empty")
	}
	return v
}
