package pkg

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	TrafficKey = "X-Request-Id"
	LoggerKey  = "_go-admin-logger-request"
)

func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}

// Assert 条件断言
// 当断言条件为 假 时触发 panic
// 对于当前请求不会再执行接下来的代码，并且返回指定格式的错误信息和错误码
func Assert(condition bool, msg string, code ...int) {
	if !condition {
		statusCode := 200
		if len(code) > 0 {
			statusCode = code[0]
		}
		panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)
	}
}

// HasError 错误断言
// 当 error 不为 nil 时触发 panic
// 对于当前请求不会再执行接下来的代码，并且返回指定格式的错误信息和错误码
// 若 msg 为空，则默认为 error 中的内容
func HasError(err error, msg string, code ...int) {
	if err != nil {
		statusCode := 200
		if len(code) > 0 {
			statusCode = code[0]
		}
		if msg == "" {
			msg = err.Error()
		}
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%v error: %#v", file, line, err)
		panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)
	}
}

// GenerateMsgIDFromContext 生成msgID
func GenerateMsgIDFromContext(c *gin.Context) string {
	requestId := c.GetHeader(TrafficKey)
	if requestId == "" {
		requestId = uuid.New().String()
		c.Header(TrafficKey, requestId)
	}
	return requestId
}

func GetOrmOrMongo(c *gin.Context) (interface{}, error) {
	idb, exist := c.Get("db")
	if !exist {
		return nil, errors.New("db connect not exist")
	}
	switch idb.(type) {
	case *gorm.DB:
		return idb.(*gorm.DB), nil
	case *mongo.Database:
		return idb.(*mongo.Database), nil
	default:
		return nil, errors.New("db connect not exist")
	}
}

// GetOrm 获取orm连接
func GetOrm(c *gin.Context) (*gorm.DB, error) {
	idb, exist := c.Get("db")
	if !exist {
		return nil, errors.New("db connect not exist")
	}
	switch idb.(type) {
	case *gorm.DB:
		//新增操作
		return idb.(*gorm.DB), nil
	default:
		return nil, errors.New("db connect not exist")
	}
}

// GetMongo 获取mongo连接
func GetMongo(c *gin.Context) (*mongo.Database, error) {
	idb, exist := c.Get("db")
	if !exist {
		return nil, errors.New("db connect not exist")
	}
	switch idb.(type) {
	case *mongo.Database:
		//新增操作
		return idb.(*mongo.Database), nil
	default:
		return nil, errors.New("db connect not exist")
	}
}

func GetMongoClient(c *gin.Context) (*mongo.Client, error) {
	client, exist := c.Get("client")
	if !exist {
		return nil, errors.New("db connect not exist")
	}
	switch client.(type) {
	case *mongo.Client:
		//新增操作
		return client.(*mongo.Client), nil
	default:
		return nil, errors.New("db connect not exist")
	}
}
