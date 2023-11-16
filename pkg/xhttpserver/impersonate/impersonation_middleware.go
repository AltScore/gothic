package impersonate

import (
	"context"
	"github.com/AltScore/gothic/v2/pkg/ids"
	"github.com/AltScore/gothic/v2/pkg/xapi"
	"github.com/AltScore/gothic/v2/pkg/xcontext"
	"github.com/AltScore/gothic/v2/pkg/xerrors"
	"github.com/AltScore/gothic/v2/pkg/xuser"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strings"
)

const (
	UserKeyName              = xcontext.UserCtxKey
	ImpersonatingUserKeyName = "x_impersonating_user"
)

// UserReader is a collaborator that is used to find users by their Id
type UserReader interface {
	FindById(ctx context.Context, id ids.Id) (xuser.User, error)
}

// ManagedUser is a user that is managed by another user
type ManagedUser interface {
	IsManagedBy(id ids.Id) bool
}

// impersonatePartnerUserMiddleware allows a user to authenticate with its owns credentials and impersonate another user
// This is used by Partner Aggregators to impersonate their managed users
// The middleware takes the current user from the context and check if the impersonate-to header exists. If exists
// and belongs to a partner managed by the current user, the current user is replaced by the impersonated user.
// Aggregator partner user is kept also in the context.
type impersonatePartnerUserMiddleware struct {
	logger                     *zap.Logger
	headerName                 string
	permissionToImpersonate    string
	permissionToImpersonateAll string
	users                      UserReader
}

// NewImpersonateUserMiddleware returns a middleware to impersonate a user
func NewImpersonateUserMiddleware(
	users UserReader,
	impersonatedHeaderName string,
	permissionToImpersonate string,
	permissionToImpersonateAll string,
	logger *zap.Logger,
) echo.MiddlewareFunc {
	xerrors.EnsureNotEmpty(users, "users")
	xerrors.EnsureNotEmpty(impersonatedHeaderName, "impersonatedHeaderName")
	xerrors.EnsureNotEmpty(permissionToImpersonate, "permissionToImpersonate")
	xerrors.EnsureNotEmpty(permissionToImpersonateAll, "permissionToImpersonateAll")
	xerrors.EnsureNotEmpty(logger, "logger")

	logger.Info("Starting impersonating middleware", zap.String("header", impersonatedHeaderName), zap.String("permission", permissionToImpersonate))

	m := &impersonatePartnerUserMiddleware{
		logger:                     logger,
		headerName:                 impersonatedHeaderName,
		permissionToImpersonate:    permissionToImpersonate,
		permissionToImpersonateAll: permissionToImpersonateAll,
		users:                      users,
	}
	return m.makeMiddlewareFunc
}

func (m *impersonatePartnerUserMiddleware) makeMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		impersonatedUserId, found := m.findImpersonatedUserId(c)

		if !found {
			// No impersonation header, ignore
			return next(c)
		}

		user, found := xapi.UserFromContext(c)

		if !found {
			// No user in context, ignore
			return next(c)
		}

		if user.Id() == impersonatedUserId {
			// Ignore if the user is trying to impersonate itself
			m.logger.Debug("ignoring impersonation to self", zap.String("id", impersonatedUserId.String()))
			return next(c)
		}

		impersonatedPartner, err := m.findImpersonatedUser(c.Request().Context(), user, impersonatedUserId)

		if err != nil {
			m.logger.Debug("error impersonating", zap.String("id", impersonatedUserId.String()), zap.String("error", err.Error()))
			return err
		}

		m.logger.Debug(
			"impersonating",
			zap.String("id", impersonatedUserId.String()),
			zap.String("name", impersonatedPartner.Name()),
			zap.String("impersonator", user.Id().String()),
		)

		m.updateUsersInContext(c, user, impersonatedPartner)

		return next(c)
	}
}

func (m *impersonatePartnerUserMiddleware) updateUsersInContext(c echo.Context, user xuser.User, impersonatedPartner xuser.User) {
	// Saves the user in the context to get extracted from there
	impersonated := &impersonatedUserType{
		User:         impersonatedPartner,
		impersonator: user,
	}

	ctx := xcontext.WithUser(c.Request().Context(), impersonated)

	c.SetRequest(c.Request().WithContext(ctx))

	c.Set(ImpersonatingUserKeyName, user)
	c.Set(UserKeyName, impersonated)
}

func (m *impersonatePartnerUserMiddleware) findImpersonatedUserId(c echo.Context) (ids.Id, bool) {
	headerValues := c.Request().Header[m.headerName]

	if len(headerValues) == 0 {
		return ids.Empty(), false
	}

	id, err := ids.Parse(headerValues[0])

	return id, err == nil
}

func (m *impersonatePartnerUserMiddleware) findImpersonatedUser(ctx context.Context, user xuser.User, id ids.Id) (xuser.User, error) {
	cannotImpersonateAll := !user.HasPermission(m.permissionToImpersonateAll)
	if cannotImpersonateAll && !user.HasPermission(m.permissionToImpersonate) {
		return nil, NewImpersonationError("User %s does not have permission to impersonate", user.Id())
	}

	impersonated, err := m.users.FindById(ctx, id)
	if err != nil {
		return nil, NewImpersonationError("%s cannot impersonate %s, does not exists", user.Id(), id)
	}

	if user.Tenant() != impersonated.Tenant() {
		m.logger.Warn("impersonation error, different tenant", zap.String("impersonator", user.Id().String()), zap.String("impersonated", id.String()))
		return nil, NewImpersonationError("%s cannot impersonate %s, it is not in the same tenant", user.Id(), id)
	}

	if cannotImpersonateAll {
		// User cannot impersonate all, check if it can impersonate this partner
		if managedUser, ok := impersonated.(ManagedUser); !ok {
			return nil, NewImpersonationError("%s cannot impersonate %s, it is not a managed user", user.Id(), id)
		} else if !managedUser.IsManagedBy(user.Id()) {
			return nil, NewImpersonationError("%s cannot impersonate %s, it is not managed by %s", user.Id(), id, user.Id())
		}
	}

	return impersonated, nil
}

type impersonatedUserType struct {
	xuser.User
	impersonator xuser.User
}

var _ xuser.ImpersonatedUser = (*impersonatedUserType)(nil)
var _ xuser.User = (*impersonatedUserType)(nil)

func (i *impersonatedUserType) RealUserId() ids.Id {
	return i.impersonator.Id()
}

func (i *impersonatedUserType) Tenant() string {
	return i.impersonator.Tenant()
}

func (i *impersonatedUserType) HasPermission(permission string) bool {
	if strings.HasSuffix(permission, ".all") {
		// Ignore "all" permissions because it is impersonating a specific user
		return false
	}
	return i.impersonator.HasPermission(permission)
}
