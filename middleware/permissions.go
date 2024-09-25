package middleware

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func PermissionMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		err := GetUserIDFromToken(c)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			c.Abort()
// 			return
// 		}
		
		// if !isAdmin(userID) && c.Param("userID") != userID {
		// 	c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		// 	c.Abort()
		// 	return
		// }

		// c.Next()
// 	}
// }

// // Example isAdmin function
// func isAdmin(userID string) bool {
// 	// Replace with your actual admin check logic
// 	return userID == "adminUserID"
// }
