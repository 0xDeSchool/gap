package ginx

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/multi_tenancy"
	"github.com/gin-gonic/gin"
)

func MultiTenancy() gin.HandlerFunc {
	return func(c *gin.Context) {
		ts := app.Get[multi_tenancy.TenantService]()
		result, err := ts.ResolveTenant(c)
		errx.CheckError(err)
		ctx := c.Request.Context()
		ctx = multi_tenancy.WithTenant(ctx, result)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
