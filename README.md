## installation
```
go get github.com/maksim-paskal/sentry-logrus-hook
```

## usage

```
package main

import (
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	sentrylogrushook "github.com/maksim-paskal/sentry-logrus-hook"
	log "github.com/sirupsen/logrus"
)

var ErrTest error = errors.New("test error")

func main() {
	hook, err := sentrylogrushook.NewHook("", "test-version", nil)
	if err != nil {
		log.WithError(err).Fatal()
	}

	log.AddHook(hook)

	log.Info("test info")
	log.WithError(ErrTest).Warn("test warn")
	log.WithError(ErrTest).Error("test error")

	defer sentry.Flush(1 * time.Second)
	defer sentry.Recover()
}
```