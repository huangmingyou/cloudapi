package wangsu

import (
	"github.com/gin-gonic/gin"
	"github.com/aWildProgrammer/fconf"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
        "os/exec"
    "crypto/hmac"
    "crypto/sha1"
        "encoding/base64"
)
type ws struct {
	id string
	key string
}
var WS ws

func init(){
		c, err := fconf.NewFileConf("./conf.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	WS = ws{c.String("wangsu.user"),c.String("wangsu.key")}
}
func getDate() string{
  f, err := exec.Command("date", "-R","-u").Output()
  if err != nil {
    fmt.Println(err.Error())
  }
  str := string(f)
  str = strings.TrimSpace(strings.Replace(str, "+0000", "GMT", -1)) //caution:if not trim, a space character is hidden at last position
  return str
}
func encrypt(accountName string,date string,apikey string) string{
  key := []byte(apikey)
  mac := hmac.New(sha1.New, key)
  mac.Write([]byte(date))
  value :=  mac.Sum(nil)
  signed_apikey := base64.StdEncoding.EncodeToString(value)
  msg := base64.StdEncoding.EncodeToString([]byte(accountName+":"+signed_apikey))
  return msg
}
func httpDo(url string,date string,auth string,method string,accept string,requestbody string) string{
  client := &http.Client{}
  req, err := http.NewRequest(method, url, strings.NewReader(requestbody))
  if err != nil {
    // handle error
  }
  req.Header.Set("Accept", accept)
  req.Header.Set("Date", date)
  req.Header.Set("Authorization", "Basic "+ auth)
  resp, err := client.Do(req)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    // handle error
  }
  //fmt.Println(string(body))
  return string(body)
}
// CalcExample godoc
// @Summary 查询cdn回源ip
// @Description 查询回源ip
// @Param   domain path string true "域名"
// @Tags cdn
// @Accept json
// @Produce json
// @Success 200 {integer} string "answer"
// @Router /wangsu/origip/{domain} [get]
func Getorig(c *gin.Context) {
  domain := c.Param("domain")
  username := WS.id
  apikey := WS.key
  method :="GET"
  accept :="application/json"
  url := "https://open.chinanetcenter.com/api/domain/"+domain
  body :=""
  date := getDate()
  auth := encrypt(username,date,apikey)
  msg := httpDo(url,date,auth,method,accept,body)
  c.String(http.StatusOK, "%s", msg)
}

