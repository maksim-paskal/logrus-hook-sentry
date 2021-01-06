## installation
```
go get github.com/maksim-paskal/sentry-logrus-hook
```

## usage

```
package main

import (
	"errors"

	sentrylogrushook "github.com/maksim-paskal/sentry-logrus-hook"
	log "github.com/sirupsen/logrus"
)

var ErrTest error = errors.New("test error")

func main() {
	hook, err := sentrylogrushook.NewHook(sentrylogrushook.SentryLogHookOptions{
		Release: "test",
	})
	if err != nil {
		log.WithError(err).Fatal()
	}

	log.AddHook(hook)

	log.Info("test info")
	log.WithError(ErrTest).Warn("test warn")
	log.WithError(ErrTest).Error("test error")

	defer hook.Stop()
}
```