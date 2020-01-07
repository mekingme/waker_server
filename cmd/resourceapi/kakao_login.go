package resourceapi

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/waker_server/pkg/database"
	"io/ioutil"
	"net/http"
)

func SetKakaoLoginAPI(router *gin.Engine, db *gorm.DB) {

	router.POST("/login/kakao", func(context *gin.Context) {

		userAccount := database.UserAccount{}

		if err := context.ShouldBindJSON(&userAccount); err != nil{
			context.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
			return
		}

		profileURL := "https://kapi.kakao.com/v1/api/talk/profile?access_token="+ userAccount.UserToken
		resp, err := http.Get(profileURL)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("Get kakao profile error : " + err.Error())
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Read kakao profile error : " + err.Error())
		}

		err = json.Unmarshal(body, &userAccount.KakaoProfile)
		if err != nil {
			fmt.Println("Unmarshal kakao profile error : " + err.Error())
		}

		if isNew := db.NewRecord(&userAccount); isNew == true {
			fmt.Println("New kakao userAccount created!")
			db.Create(&userAccount)
		}

		context.JSON(http.StatusOK,gin.H{"status" : "success"})
	})

}
