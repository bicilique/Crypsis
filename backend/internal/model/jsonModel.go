package model

import "github.com/gin-gonic/gin"

// JSONSuccessResponse sends a JSON response with a success status.
// The response will have three fields: "success", "message" and "data".
// The "success" field will always be true.
// The "message" field will contain the message passed to the function.
// The "data" field will contain the data passed to the function, it can be any type of data.
func JSONSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// JSONSuccessResponseWithCount sends a response with a count field which is the number of items in the response.
// The data can be any type of data, but it must be serializable to JSON.
func JSONSuccessResponseWithCount(c *gin.Context, statusCode int, message string, count int64, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"count":   count,
		"data":    data,
	})
}

// JSONErrorResponse sends a JSON response with a success status set to false.
// The response will have three fields: "success", "message" and "error".
// The "success" field will always be false.
// The "message" field will contain the message passed to the function.
// The "error" field will contain the data passed to the function, it can be any type of data.
func JSONErrorResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"error":   data,
	})
	c.Abort()
}
