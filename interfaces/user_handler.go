package interfaces

import (
	"etomne/application"
	"etomne/infrastructure/auth"
	"github.com/gin-gonic/gin"
)

type Users struct {
	us application.UserApiInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewUsers(us application.UserApiInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Users {
	return &Users{
		us: us,
		rd: rd,
		tk: tk,
	}
}

func (s *Users) SaveUser(c *gin.Context) {

}
