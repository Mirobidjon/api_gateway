package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bitbucket.org/udevs/api_gateway/api/models"
	"bitbucket.org/udevs/api_gateway/config"
	"bitbucket.org/udevs/api_gateway/pkg/logger"
	services "bitbucket.org/udevs/api_gateway/services"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrAlreadyExists       = "ALREADY_EXISTS"
	ErrNotFound            = "NOT_FOUND"
	ErrInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavailable  = "SERVICE_UNAVAILABLE"
	SigningKey             = []byte("FfLbN7pIEYe8@!EqrttOLiwa(H8)7Ddo")
	SuperAdminUserType     = "superadmin"
	SystemUserType         = "admin"
)

type handlerV1 struct {
	log      logger.Logger
	cfg      config.Config
	services services.ServiceManager
}

type HandlerV1Options struct {
	Log      logger.Logger
	Cfg      config.Config
	Services services.ServiceManager
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		log:      options.Log,
		cfg:      options.Cfg,
		services: options.Services,
	}
}

func (h *handlerV1) handleErrorResponse(c *gin.Context, code int, message string, err interface{}) {
	h.log.Error(message, logger.Int("code", code), logger.Any("error", err))
	c.JSON(code, models.ResponseModel{
		Code:    code,
		Message: message,
		Error:   err,
	})
}

func (h *handlerV1) handleSuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	h.log.Info(message, logger.Any("response", data))
	c.JSON(code, models.ResponseModel{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func (h *handlerV1) ParseQueryParam(c *gin.Context, key string, defaultValue string) (int, error) {
	valueStr := c.DefaultQuery(key, defaultValue)

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		h.log.Error("error while parsing query param"+", canceled ", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return 0, err
	}

	return value, nil
}

func (h *handlerV1) BadRequestResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}

func ProtoToStruct(data interface{}, m protoreflect.ProtoMessage) error {
	jsonMarshaller := protojson.MarshalOptions{
		AllowPartial:    true,
		EmitUnpopulated: true,
	}

	js, err := jsonMarshaller.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(js, data)
	return err
}
