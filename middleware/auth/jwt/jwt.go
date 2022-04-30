package jwt

import (
	"strings"

	sysctx "context"

	"github.com/zhiyunliu/gel/context"

	"github.com/golang-jwt/jwt/v4"

	"github.com/zhiyunliu/gel/errors"
	"github.com/zhiyunliu/gel/middleware"
)

type authKey struct{}

const (

	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// bearerFormat authorization token format
	bearerFormat string = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey string = "Authorization"

	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized("JWT token is missing")
	ErrMissingKeyFunc         = errors.Unauthorized("keyFunc is missing")
	ErrTokenInvalid           = errors.Unauthorized("Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized("Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("Wrong signing method")
	ErrWrongContext           = errors.Unauthorized("Wrong context for middleware")
	ErrNeedTokenProvider      = errors.Unauthorized("Token provider is missing")
	ErrSignToken              = errors.Unauthorized("Can not sign token.Is the key correct?")
	ErrGetKey                 = errors.Unauthorized("Can not get key while signing token")
)

// Option is jwt option.
type Option func(*options)

// Parser is a jwt parser
type options struct {
	signingMethod jwt.SigningMethod
}

// WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// Server is a server auth middleware. Check the token and extract the info from token.
func Server(keyFunc jwt.Keyfunc, opts ...Option) middleware.Middleware {
	o := &options{
		signingMethod: jwt.SigningMethodHS256,
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context) (reply interface{}) {

			if keyFunc == nil {
				return ErrMissingKeyFunc
			}
			authVal := ctx.Header(authorizationKey)

			auths := strings.SplitN(authVal, " ", 2)
			if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
				return ErrMissingJwtToken
			}
			jwtToken := auths[1]
			var (
				tokenInfo *jwt.Token
				err       error
			)

			tokenInfo, err = jwt.Parse(jwtToken, keyFunc)

			if err != nil {
				if ve, ok := err.(*jwt.ValidationError); ok {
					if ve.Errors&jwt.ValidationErrorMalformed != 0 {
						return ErrTokenInvalid
					} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
						return ErrTokenExpired
					} else {
						return ErrTokenParseFail
					}
				}
				return errors.Unauthorized(err.Error())
			} else if !tokenInfo.Valid {
				return ErrTokenInvalid
			} else if tokenInfo.Method != o.signingMethod {
				return ErrUnSupportSigningMethod
			}
			ctx = NewContext(ctx, tokenInfo.Claims)
			return handler(ctx)
		}

	}
}

// NewContext put auth info into context
func NewContext(ctx context.Context, info jwt.Claims) context.Context {
	nctx := sysctx.WithValue(ctx.Context(), authKey{}, info)
	ctx.ResetContext(nctx)
	return ctx
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (token jwt.Claims, ok bool) {
	token, ok = ctx.Context().Value(authKey{}).(jwt.Claims)
	return
}
