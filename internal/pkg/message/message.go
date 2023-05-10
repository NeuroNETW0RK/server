package message

import (
	"net/http"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	ErrCode int         `json:"errcode"`
	Data    interface{} `json:"data"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	response := new(Response)
	response.Data = data
	response.Msg = "ok"
	response.Code = 200
	response.ErrCode = 100000

	if err != nil {
		log.C(c).Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		response.Msg = coder.String()
		response.Code = coder.HTTPStatus()
		response.ErrCode = coder.Code()

		c.JSON(
			coder.HTTPStatus(),
			response,
		)
		return
	}

	c.JSON(http.StatusOK, response)
}

func success(c *gin.Context, msg string, data interface{}) {
	response := new(Response)
	response.Data = data
	response.Msg = msg
	response.Code = 200
	response.ErrCode = 100000
	c.JSON(http.StatusOK, response)
}

func failed(c *gin.Context, msg string, err error) {
	log.C(c).Errorf("%#+v", err)
	response := new(Response)
	response.Data = nil
	coder := errors.ParseCoder(err)
	if msg != "" {
		response.Msg = msg
	} else {
		response.Msg = coder.String()
	}
	response.Code = coder.HTTPStatus()
	response.ErrCode = coder.Code()

	c.JSON(
		coder.HTTPStatus(),
		response,
	)
}

func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	success(c, msg, data)
}

func Success(c *gin.Context, data interface{}) {
	success(c, "ok", data)
}

func FailedWithMsg(c *gin.Context, msg string, err error) {
	failed(c, msg, err)
}

func Failed(c *gin.Context, err error) {
	failed(c, "", err)
}
