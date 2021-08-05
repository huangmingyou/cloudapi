package tencent

import (
        "fmt"
	"github.com/gin-gonic/gin"
	"github.com/aWildProgrammer/fconf"
	"net/http"

        "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
        "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
        "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
        billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
)

type tx struct {
	id string
	key string
}
var TX tx

func init(){
		c, err := fconf.NewFileConf("./conf.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	TX = tx{c.String("tencent.id"),c.String("tencent.key")}
}
// CalcExample godoc
// @Summary 查询腾讯云余额
// @Description 查询余额
// @Tags balance
// @Accept json
// @Produce json
// @Success 200 {integer} string "answer"
// @Router /tencent/balance [get]
func TencentBalance(c *gin.Context) {

        credential := common.NewCredential(
                TX.id,
                TX.key,
        )
        cpf := profile.NewClientProfile()
        cpf.HttpProfile.Endpoint = "billing.tencentcloudapi.com"
        client, _ := billing.NewClient(credential, "", cpf)

        request := billing.NewDescribeAccountBalanceRequest()


        response, err := client.DescribeAccountBalance(request)
        if _, ok := err.(*errors.TencentCloudSDKError); ok {
                fmt.Printf("An API error has returned: %s", err)
                return
        }
        if err != nil {
                panic(err)
        }
        fmt.Printf("%s", response.ToJsonString())
  	//c.String(http.StatusOK, "%s", response.ToJsonString())
  	c.String(http.StatusOK, "%d\n", *response.Response.Balance/100)
} 
