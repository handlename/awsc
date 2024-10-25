package aws

import (
	"fmt"
	"regexp"

	"github.com/handlename/awsc/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

var readonlyActionRx = regexp.MustCompile(`^((get|list|describe|select)-.*)|(ls)$`)

// IsReadonly determines whether the command will only read resources on AWS.
func (s Service) IsReadonly(argv []string) (bool, error) {
	if (len(argv) == 1 && argv[0] == "help") ||
		(len(argv) == 2 && argv[1] == "help") ||
		(len(argv) == 3 && argv[2] == "help") {
		return true, nil
	}

	if len(argv) < 2 {
		return false, failure.New(
			errorcode.ErrInvalidArgument,
			failure.Message("argv is too short"),
			failure.Context{"argv": fmt.Sprintf("%+v", argv)},
		)
	}

	service, action := argv[0], argv[1]
	if service == "" || action == "" {
		return false, failure.New(
			errorcode.ErrInvalidArgument,
			failure.Message("service or action is empty"),
			failure.Context{"service": service, "action": action},
		)
	}

	if readonlyActionRx.MatchString(action) {
		return true, nil
	}

	return false, nil
}
