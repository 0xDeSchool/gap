package short_url

import (
	"github.com/0xDeSchool/gap/app"
	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/ginx"
	"github.com/gin-gonic/gin"
)

func init() {
	ginx.Configure(func(s *ginx.Server) error {
		g := s.G.Group("/api/u")
		g.GET(":key", goUrl)
		return nil
	})
}

func goUrl(ctx *gin.Context) {
	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(400, gin.H{
			"message": "key is required",
		})
		return
	}
	repo := app.GetPtr[ShortUrlRepository]()
	u, err := repo.GetUrl(ctx, key)
	errx.CheckError(err)
	if u == nil || u.Url == "" {
		ctx.JSON(404, gin.H{
			"message": "page not found",
		})
		return
	}
	ctx.JSON(200, ShortUrlInput{
		Url: u.Url,
	})
}
