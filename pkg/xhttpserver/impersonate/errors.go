package impersonate

import "github.com/AltScore/gothic/v2/pkg/xerrors"

func NewImpersonationError(msg string, args ...any) error {
	return xerrors.NewUnauthorized(msg, args...)
}
