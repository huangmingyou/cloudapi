package aliyun

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/aWildProgrammer/fconf"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"

)

//   信道

var cmsg chan string = make(chan string)
var cmsg2 chan string = make(chan string)

// 账号信息
type ak struct {
	id string
	key string
	region string
}
var AK ak

func init(){
		c, err := fconf.NewFileConf("./conf.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	AK = ak{c.String("aliyun.id"),c.String("aliyun.key"),c.String("aliyun.region")}
}

// CalcExample godoc
// @Summary 查询阿里云余额
// @Description 查询余额
// @Tags balance
// @Accept json
// @Produce json
// @Success 200 {integer} string "answer"
// @Router /aliyun/balance [get]
func Qbalance(c *gin.Context) {
	client, err := bssopenapi.NewClientWithAccessKey(AK.region, AK.id, AK.key)

	request := bssopenapi.CreateQueryAccountBalanceRequest()
	request.Scheme = "https"


	response, err := client.QueryAccountBalance(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	//fmt.Printf("response is %s\n", response.GetHttpContentString())
	c.String(http.StatusOK,"%s\n",response.Data.AvailableAmount)
}
// 查询vg 的 ecs

func slbvgecs(vgid string) []string{
	client, err := slb.NewClientWithAccessKey(AK.region, AK.id, AK.key)

	request := slb.CreateDescribeVServerGroupAttributeRequest()
	request.Scheme = "https"

	request.VServerGroupId = vgid

	response, err := client.DescribeVServerGroupAttribute(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	//fmt.Printf("response is %#v\n", response)
	var ecslist []string
	count :=len(response.BackendServers.BackendServer)
	for i := 0; i<count;i++ {
		ecslist = append(ecslist,response.BackendServers.BackendServer[i].ServerId)
	}
	fmt.Println(ecslist)
	return ecslist
}


// 查询slb的vg


func slbvg(slbid string) []string{
	client, err := slb.NewClientWithAccessKey(AK.region, AK.id, AK.key)

	request := slb.CreateDescribeVServerGroupsRequest()
	request.Scheme = "https"

	request.LoadBalancerId = slbid

	response, err := client.DescribeVServerGroups(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	//fmt.Printf("response is %#v\n", response)
	var vglist []string
	count :=len(response.VServerGroups.VServerGroup)
	for i := 0; i<count;i++ {
		vglist = append(vglist,response.VServerGroups.VServerGroup[i].VServerGroupId)
	}
	return vglist
}








// 根据ecs id 查询 ip 
func ecsid2ip(ecsid string) string {
	//return "hello"

	client2, err := ecs.NewClientWithAccessKey(AK.region, AK.id, AK.key)

	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.Scheme = "https"

	request.InstanceId = ecsid

	response, err := client2.DescribeInstanceAttribute(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	//fmt.Printf("response is %#v\n", response)
//	fmt.Printf("%s\n",response.VpcAttributes.PrivateIpAddress.IpAddress[0])
	cmsg <- response.VpcAttributes.PrivateIpAddress.IpAddress[0]
	return response.VpcAttributes.PrivateIpAddress.IpAddress[0]
}

// CalcExample godoc
// @Summary 列出阿里云slb对应的后端服务器ip
// @Description plus
// @Tags slb
// @Param  ip path string true "slb ip"
// @Accept json
// @Produce json
// @Success 200 {integer} string "answer"
// @Router /aliyun/slb/{ip} [get]
func Listslbip(c *gin.Context) {
	ip := c.Param("ip")
	client, err := slb.NewClientWithAccessKey(AK.region, AK.id, AK.key)

	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	//request.Address = "47.99.3.94"
	request.Address = ip

	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		fmt.Print(err.Error())
	}

	r2 := slb.CreateDescribeLoadBalancerAttributeRequest()
	r2.Scheme = "https"
	r2.LoadBalancerId=response.LoadBalancers.LoadBalancer[0].LoadBalancerId

	respons2,err2 := client.DescribeLoadBalancerAttribute(r2)
	if err2 != nil {
		fmt.Print(err.Error())
	}
	msg := ""
	ecsid := ""
	i := 0
	count := len(respons2.BackendServers.BackendServer)

	if count > 0 {
	for i = 0 ; i<count; i++ {
		ecsid = respons2.BackendServers.BackendServer[i].ServerId
		go ecsid2ip(ecsid)
	}
	for i = 0 ; i<count; i++ {
		msg += <-cmsg
		msg += "\n"
	}
	}else{
		slbid := response.LoadBalancers.LoadBalancer[0].LoadBalancerId
		vglist := slbvg(slbid)
		ecslist := slbvgecs(vglist[0])

		count := len(ecslist)
		for i = 0 ;i<count ;i++ {
			go ecsid2ip(ecslist[i])
		}
		for i = 0 ;i<count ;i++ {
			msg += <-cmsg
			msg += "\n"
		}
	}
	c.String(http.StatusOK, "%s", msg)


}

// CalcExample godoc
// @Summary 列出所有阿里云slb
// @Description plus
// @Tags slb
// @Accept json
// @Produce json
// @Success 200 {integer} string "answer"
// @Router /aliyun/slb [get]
func Listallslb(c *gin.Context) {
	//nettype := c.Param("nettype")
	client, err := slb.NewClientWithAccessKey(AK.region, AK.id, AK.key)

	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	//request.AddressType= nettype

	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	c.String(http.StatusOK,"%s\n",response.GetHttpContentString())
}
