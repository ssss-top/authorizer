package test

import (
	"fmt"
	"testing"

	"github.com/authorizerdev/authorizer/server/constants"
	"github.com/authorizerdev/authorizer/server/db"
	"github.com/authorizerdev/authorizer/server/envstore"
	"github.com/authorizerdev/authorizer/server/graph/model"
	"github.com/authorizerdev/authorizer/server/resolvers"
	"github.com/authorizerdev/authorizer/server/sessionstore"
	"github.com/stretchr/testify/assert"
)

func logoutTests(t *testing.T, s TestSetup) {
	t.Helper()
	t.Run(`should logout user`, func(t *testing.T) {
		req, ctx := createContext(s)
		email := "logout." + s.TestInfo.Email

		_, err := resolvers.MagicLinkLoginResolver(ctx, model.MagicLinkLoginInput{
			Email: email,
		})

		verificationRequest, err := db.Provider.GetVerificationRequestByEmail(email, constants.VerificationTypeMagicLinkLogin)
		verifyRes, err := resolvers.VerifyEmailResolver(ctx, model.VerifyEmailInput{
			Token: verificationRequest.Token,
		})

		token := *verifyRes.AccessToken
		sessions := sessionstore.GetUserSessions(verifyRes.User.ID)
		cookie := ""
		// set all they keys in cookie one of them should be session cookie
		for key := range sessions {
			if key != token {
				cookie += fmt.Sprintf("%s=%s;", envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyCookieName)+"_session", key)
			}
		}

		req.Header.Set("Cookie", cookie)
		_, err = resolvers.LogoutResolver(ctx)
		assert.Nil(t, err)
		cleanData(email)
	})
}
