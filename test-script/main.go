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
