package interfaces

import (
	"fmt"
	"github.com/AhEhIOhYou/etomne/backend/application"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Comment struct {
	modelApp application.ModelAppInterface
	userApp  application.UserAppInterface
	comApp   application.CommentAppInterface
	rd       auth.AuthInterface
	tk       auth.TokenInterface
}

func NewComment(mApp application.ModelAppInterface, uApp application.UserAppInterface, comApp application.CommentAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Comment {
	return &Comment{
		modelApp: mApp,
		userApp:  uApp,
		comApp:   comApp,
		rd:       rd,
		tk:       tk,
	}
}

// SaveComment doc
// @Summary		Save new comment
// @Tags		Comment
// @Accept		mpfd
// @Produce		json
// @Param		model_id  formData  string  true  "Model ID"
// @Param		message   formData  string  true  "Message"
// @Success		201  {object}  entities.Comment
// @Failure     401  string  unauthorized
// @Failure     400  string  user not found, unauthorized
// @Failure     422  string  error
// @Failure     500  string  error
// @Router		/comment [post]
// @Security	bearerAuth
func (com *Comment) SaveComment(c *gin.Context) {
	metadata, err := com.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, err := com.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = com.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	modelId, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	message := c.PostForm("message")
	if fmt.Sprintf("%T", message) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "invalid_json")
		return
	}

	emptyComment := entities.Comment{
		Message: message,
	}

	saveCommentError := emptyComment.Validate("")
	if len(saveCommentError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveCommentError)
		return
	}

	var Comment = entities.Comment{
		AuthorId: userId,
		ModelId:  modelId,
		Message:  message,
	}

	Comment.Prepare()

	saveComment, saveErr := com.comApp.SaveComment(&Comment)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}

	c.JSON(http.StatusCreated, saveComment)
}

// GetComments doc
// @Summary		Get all comments by Model ID
// @Tags		Comment
// @Accept 		json
// @Produce		json
// @Param		model_id  query  string  true   "Model ID"
// @Param		count  	  query  string  false  "Count"
// @Success		201  {object}  []entities.Comment
// @Failure     401  string  unauthorized
// @Failure     400  string  user not found, unauthorized
// @Failure     422  string  error
// @Failure     500  string  error
// @Router		/comment [get]
func (com *Comment) GetComments(c *gin.Context) {

	modelId, err := strconv.ParseUint(c.Query("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "invalid query")
		return
	}

	replyId, _ := strconv.ParseUint(c.Query("reply_id"), 10, 64)

	count, _ := strconv.ParseUint(c.Query("count"), 10, 64)
	if count == 0 {
		count = 2
	}

	var comments []entities.Comment

	if replyId != 0 {
		comments, err = com.comApp.GetReplies(replyId, count)
	} else {
		comments, err = com.comApp.GetCommentsByModel(modelId, count)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, comments)
}
