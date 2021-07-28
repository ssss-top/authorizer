package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/authorizerdev/authorizer/server/db"
	"github.com/authorizerdev/authorizer/server/enum"
	"github.com/authorizerdev/authorizer/server/graph/model"
	"github.com/authorizerdev/authorizer/server/session"
	"github.com/authorizerdev/authorizer/server/utils"
)

func Token(ctx context.Context) (*model.LoginResponse, error) {
	gc, err := utils.GinContextFromContext(ctx)
	var res *model.LoginResponse
	if err != nil {
		return res, err
	}
	token, err := utils.GetAuthToken(gc)
	if err != nil {
		return res, err
	}

	claim, accessTokenErr := utils.VerifyAuthToken(token)
	expiresAt := claim.ExpiresAt

	user, err := db.Mgr.GetUserByEmail(claim.Email)
	if err != nil {
		return res, err
	}

	userIdStr := fmt.Sprintf("%d", user.ID)

	sessionToken := session.GetToken(userIdStr)

	if sessionToken == "" {
		return res, fmt.Errorf(`unauthorized`)
	}
	// TODO check if refresh/session token has expired

	expiresTimeObj := time.Unix(expiresAt, 0)
	currentTimeObj := time.Now()
	if accessTokenErr != nil || expiresTimeObj.Sub(currentTimeObj).Minutes() <= 5 {
		// if access token has expired and refresh/session token is valid
		// generate new accessToken
		token, expiresAt, _ = utils.CreateAuthToken(utils.UserAuthInfo{
			ID:    userIdStr,
			Email: user.Email,
		}, enum.AccessToken)
	}
	utils.SetCookie(gc, token)
	res = &model.LoginResponse{
		Message:              `Email verified successfully.`,
		AccessToken:          &token,
		AccessTokenExpiresAt: &expiresAt,
		User: &model.User{
			ID:        userIdStr,
			Email:     user.Email,
			Image:     &user.Image,
			FirstName: &user.FirstName,
			LastName:  &user.LastName,
			CreatedAt: &user.CreatedAt,
			UpdatedAt: &user.UpdatedAt,
		},
	}
	return res, nil
}