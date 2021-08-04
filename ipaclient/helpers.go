package ipaclient

import (
	"auri/config"
	"auri/logger"

	"github.com/pkg/errors"
)

// runIPATransaction executes the given transaction code if IPA is enabled
// otherwise a a log message is emitted
func runIPATransaction(name string, f func(IPA) error) error {
	// we log failures and give anonymous error message back in order to protect
	// possible sensitive details here
	hiddenErrorLog := func(err error) error {
		logger.IPALogger.Error(err)
		return errors.New("IPA operation failed")
	}

	logger.IPALogger.Infof("IPA: %v", name)
	if !config.GetInstance().IPAEnable {
		logger.BuffaloLogger.Warnf("IPA: IPA functionality is disabled. Skipping the operation '%v'", name)
		return nil
	}

	//protect our errors and avoid public information about our technical details here
	//as this might leak secrets
	client, err := GetClient()
	if err != nil {
		return hiddenErrorLog(err)
	}

	if err = f(client); err != nil {
		return hiddenErrorLog(err)
	}

	return nil
}
