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
package logrushooksentry

import (
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Hook struct {
	options   Options
	logLevels []log.Level
}

type Options struct {
	SentryDSN string
	Release   string
	LogLevels []log.Level
	Tags      map[string]string
	// default: 1s
	FlushDuration time.Duration
}

const RequestKey = "request"

// create new Hook.
func NewHook(options Options) (*Hook, error) {
	hook := Hook{
		options: options,
	}

	if hook.options.FlushDuration == 0 {
		hook.options.FlushDuration = time.Second
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:     options.SentryDSN,
		Release: options.Release,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Sentry init failed")
	}

	if options.Tags != nil {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTags(options.Tags)
		})
	}

	hook.logLevels = options.LogLevels

	// use errors levels for default
	if hook.logLevels == nil {
		hook.logLevels = []log.Level{
			log.ErrorLevel,
			log.FatalLevel,
			log.WarnLevel,
			log.PanicLevel,
		}
	}

	return &hook, nil
}

// Graceful stop sentry.
func (hook *Hook) Stop() {
	sentry.Flush(hook.options.FlushDuration)
	sentry.Recover()
	time.Sleep(hook.options.FlushDuration)
}

// func to interface log.Hook.Levels.
func (hook *Hook) Levels() []log.Level {
	return hook.logLevels
}

// func to interface log.Hook.Fire.
func (hook *Hook) Fire(entry *log.Entry) error {
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

		if entry.HasCaller() {
			scope.SetExtra("Caller", entry.Caller)
		}

		for key, value := range entry.Data {
			switch key {
			case log.ErrorKey:
				// localHub.CaptureException don't save in sentry message
				if len(entry.Message) > 0 {
					scope.SetExtra("Message", entry.Message)
				}
			case RequestKey:
				scope.SetRequest(value.(*http.Request))
			default:
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
		hook.Stop()
	}

	return nil
}
