package ginx

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/multi_tenancy"
	"github.com/gin-gonic/gin"
)

func MultiTenancy() gin.HandlerFunc {
	return func(c *gin.Context) {
		resolver := app.Get[multi_tenancy.TenantResolver]()
		result, err := resolver.ResolveTenantIdOrName(c)
		errx.CheckError(err)
		ctx := c.Request.Context()
		ctx = multi_tenancy.WithTenant(ctx, result)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
