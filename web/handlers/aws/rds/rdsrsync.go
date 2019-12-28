package rds

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/weiqiang333/devops/internal/awscli"
	"github.com/weiqiang333/devops/web/handlers/auth"
)


//PostRdsRsyncWorkorder 提交同步数据库工单
func PostRsyncWorkorder(c *gin.Context) {
	username := fmt.Sprint(auth.Me(c))
	databaseName, dnOk := c.GetPostForm("databaseNmae")
	qrCode, qrOk := c.GetPostForm("qrcode")
	if ! dnOk {
		c.JSON(http.StatusConflict, gin.H{"response": "请正确提交数据库名称"})
		return
	}
	if ! qrOk {
		c.JSON(http.StatusConflict, gin.H{"response": "请正确提交二次验证码"})
		return
	}
	ok, err := auth.VerifyCode(username, qrCode)
	if ! ok {
		c.JSON(http.StatusConflict, gin.H{"response": "二次验证失败"})
		return
	}

	err = awscli.CreateWorkorder(databaseName, username)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"response": fmt.Sprintf("创建失败: %v", err)})
		return
	}
	log.Printf("PostRdsRsyncWorkorder Success: %s %s", username, databaseName)
	c.JSON(http.StatusOK, gin.H{"response": "Create Success"})
	return
}


//GetRsyncWorkorder
func GetSyncWorkorder(c *gin.Context)  {
	username := fmt.Sprint(auth.Me(c))
	rdsWorkorder, err := awscli.SearchWorkorder(0)
	if err != nil {
		c.HTML(http.StatusPaymentRequired, "awsrds/rdsrsync.tmpl", gin.H{
			"awsAdmin": "active",
			"user": username,
			"rdsWorkorder": rdsWorkorder,
			"lastCreate": rdsWorkorder[0],
		})
		return
	}
	c.HTML(http.StatusOK, "awsrds/rdsrsync.tmpl", gin.H{
		"awsAdmin": "active",
		"user": username,
		"rdsWorkorder": rdsWorkorder,
		"lastCreate": rdsWorkorder[0],
	})
	return
}


//GetOrderId
func GetOrderId(c *gin.Context)  {
	username := fmt.Sprint(auth.Me(c))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	rdsWorkorder, err := awscli.SearchWorkorder(id)
	awscli.IfRsyncStatus(id)
	if err != nil {
		log.Printf("GetOrderId SearchOrder fail: %v", err)
	}
	rdsOrder, err := awscli.SearchOrder(id, rdsWorkorder[0].Database)
	if err != nil {
		log.Printf("GetOrderId SearchOrder fail: %v", err)
	}
	orderLogs,err := awscli.SearchOrderLogs(id)
	if err != nil {
		log.Printf("GetOrderId SearchOrderLogs fail: %v", err)
	}
	c.HTML(http.StatusOK, "awsrds/rdsrsyncorder.tmpl", gin.H{
		"awsAdmin": "active",
		"user": username,
		"rdsWorkorder": rdsWorkorder[0],
		"rdsOrder": rdsOrder,
		"orderLogs": orderLogs,
	})
}


//PostOrder 提交审批
func PostApproval(c *gin.Context)  {
	workorderId, errW := strconv.Atoi(c.Param("id"))
	orderId, errO  := strconv.Atoi(c.Query("orderId"))
	approvalStatus := c.Query("approvalStatus")
	if errO != nil || errW != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("./%v", workorderId))
		return
	}

	var status bool
	if approvalStatus == "agree" {
		status = true
	} else if approvalStatus == "reject" {
		status = false
	}
	err := awscli.InsertOrderLog(workorderId, orderId, status)
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("./%v", workorderId))
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("./%v", workorderId))
	return
}
