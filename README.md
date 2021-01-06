## installation
```
go get github.com/maksim-paskal/logrus-hook-sentry
```
## envieronment
```bash
export SENTRY_DSN=https://123@sentry.com/345
```

## usage

```go
package main

import (
	"errors"

	sentrylogrushook "github.com/maksim-paskal/logrus-hook-sentry"
	log "github.com/sirupsen/logrus"
)

var ErrTest error = errors.New("test error")

func main() {
	hook, err := sentrylogrushook.NewHook(sentrylogrushook.Options{
		Release: "test",
	})
	if err != nil {
		log.WithError(err).Fatal()
	}

	log.AddHook(hook)

	log.Info("test info")
	log.Warn(ErrTest)
	log.WithError(ErrTest).Error("some message")

	defer hook.Stop()
}
```