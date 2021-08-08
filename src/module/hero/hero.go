package hero

import (
	"fmt"
	heromodel "srv/src/module/hero/model"

	"github.com/gin-gonic/gin"
)

// SrvListHero ..
func SrvListHero(c *gin.Context) {
	heroArr := FetchHeroFromDatabase()

	c.JSON(200, gin.H{
		"heroArr": heroArr,
	})
}

func FetchHeroFromDatabase() []heromodel.Hero {
	heroArr := make([]heromodel.Hero, 0)

	for i := 0; i < 10; i++ {
		heroArr = append(heroArr, heromodel.Hero{
			ID:   i,
			Name: fmt.Sprintf("hero %d", i),
		})
	}

	return heroArr
}
