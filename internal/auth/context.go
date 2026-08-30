package auth

import "github.com/gin-gonic/gin"

const UserContextKey = "authenticated_user"

func SetUser(c *gin.Context, user *UserContext) {
	c.Set(UserContextKey, user)
}

func GetUser(c *gin.Context) (*UserContext, bool) {
	value, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}

	user, ok := value.(*UserContext)
	return user, ok
}