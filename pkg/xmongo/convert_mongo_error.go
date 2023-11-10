package xmongo

import (
	"fmt"
	"github.com/AltScore/gothic/v2/pkg/xerrors"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

func ConvertMongoError(err error, entity string, keyFmt string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	message := err.Error()
	if strings.Contains(message, "not found") {
		return xerrors.NewNotFoundError(entity, keyFmt, args...)
	}

	if strings.Contains(message, "no documents in result") {
		return xerrors.NewNotFoundError(entity, keyFmt, args...)
	}

	if strings.Contains(message, "duplicate key") {
		logUnexpectedError(err)
		return convertDuplicateKey(err, entity, keyFmt, args...)
	}

	if strings.Contains(message, "invalid key") {
		logUnexpectedError(err)
		return xerrors.NewNotFoundError(entity, keyFmt, args...)
	}

	if strings.Contains(message, "wants one but found many") {
		logUnexpectedError(err)
		return xerrors.NewFoundManyError(entity, keyFmt, args...)
	}

	if strings.Contains(message, "context deadline exceeded") {
		return xerrors.NewTimeoutError(entity, keyFmt, args...)
	}

	if strings.Contains(message, "context canceled") {
		return xerrors.NewCancellationError(entity, keyFmt, args...)
	}

	logUnexpectedError(err)
	return xerrors.NewUnknownError(entity, cleanMongoErrorMessage(err), keyFmt, args...)
}

func logUnexpectedError(err error) {
	zap.L().Error("unexpected mongo error", zap.Error(err))
}

func cleanMongoErrorMessage(err error) string {
	return strings.ReplaceAll(err.Error(), "mongo: ", "")
}

var duplicateMsgRegex = regexp.MustCompile(`duplicate key error collection: (\w+)\.(\w+) index: (\w+) dup key: \{ ([^}]+) }`)

func convertDuplicateKey(err error, entity string, keyFmt string, args ...interface{}) error {
	// Sample error format
	// E11000 duplicate key error collection: credits_staging.clients index: partnerId_externalId_index dup key: { partnerId: "62d6c5b7bfc4940dc844c610", externalId: "charly-id-000001" }

	matches := duplicateMsgRegex.FindStringSubmatch(err.Error())

	var details string
	if len(matches) < 5 {
		details = err.Error()
	} else {
		details = fmt.Sprintf("duplicate key %s on %s", strings.ReplaceAll(matches[4], `"`, "'"), matches[2])
	}

	return xerrors.NewDuplicateError(entity, details, keyFmt, args...)
}
