package watchdog

import (
	"time"

	"github.com/gin-gonic/gin"
)

func CheckImageGen(r *gin.Engine) {
	go func() {
		for {
			now := time.Now()
			//jede Minute
			next := now.Truncate(time.Minute).Add(time.Minute)
			time.Sleep(time.Until(next))

			go PingImageGen(r)
		}
	}()
}

func CheckUntis() {
	go func() {
		for {
			now := time.Now()
			//jede Minute
			next := now.Truncate(time.Minute).Add(time.Minute)
			time.Sleep(time.Until(next))

			go PingUntis()
		}
	}()
}
