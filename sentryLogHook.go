/*
Copyright paskal.maksim@gmail.com
Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sentrylogrushook

import (
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

type SentryLogHook struct {
	// logLevels to fire message to sentry
	logLevels []log.Level
}

func NewHook(sentryDSN string, release string, logLevels []log.Level) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:     sentryDSN,
		Release: release,
	})
	if err != nil {
		log.WithError(err).Fatal()
	}

	if logLevels != nil {
		logLevels = []log.Level{
			log.ErrorLevel,
			log.FatalLevel,
			log.WarnLevel,
			log.PanicLevel,
		}
	}

	log.AddHook(&SentryLogHook{
		logLevels: logLevels,
	})
}

func (slh *SentryLogHook) Levels() []log.Level {
	return slh.logLevels
}

func (slh *SentryLogHook) Fire(entry *log.Entry) error {
	sentryLevel := sentry.LevelInfo

	switch entry.Level {
	case log.PanicLevel:
		sentryLevel = sentry.LevelFatal
	case log.FatalLevel:
		sentryLevel = sentry.LevelFatal
	case log.ErrorLevel:
		sentryLevel = sentry.LevelError
	case log.WarnLevel:
		sentryLevel = sentry.LevelWarning
	case log.InfoLevel:
		sentryLevel = sentry.LevelInfo
	case log.DebugLevel:
		sentryLevel = sentry.LevelDebug
	case log.TraceLevel:
		sentryLevel = sentry.LevelDebug
	}

	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentryLevel)
		for key, value := range entry.Data {
			if key == log.ErrorKey {
				// localHub.CaptureException don't save in sentry message
				scope.SetExtra("Message", entry.Message)
			} else {
				scope.SetExtra(key, value)
			}
		}
	})

	if err, ok := entry.Data[log.ErrorKey].(error); ok && err != nil {
		localHub.CaptureException(err)
	} else {
		localHub.CaptureMessage(entry.Message)
	}

	if entry.Level <= log.FatalLevel {
		sentry.Flush(time.Second)
	}

	return nil
}
