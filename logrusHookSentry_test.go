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
package logrushooksentry_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	logrushooksentry "github.com/maksim-paskal/logrus-hook-sentry"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

var ErrTest = errors.New("test error")

func TestHook(t *testing.T) {
	t.Parallel()

	mapTags := make(map[string]string)

	mapTags["test"] = "value"

	hook, err := logrushooksentry.NewHook(ctx, logrushooksentry.Options{
		Release: "test",
		Tags:    mapTags,
	})
	if err != nil {
		t.Fatal(err)
	}

	log.AddHook(hook)

	log.Info("test info")
	log.Warn(ErrTest)
	log.WithError(ErrTest).Error("some message")
}

func TestRequestJson(t *testing.T) {
	t.Parallel()

	log.SetFormatter(&log.JSONFormatter{})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://127.0.0.1?test=value", nil)
	req.Header.Add("key", "value")

	logData := logrushooksentry.AddRequest(req)

	if logData[logrushooksentry.RequestMethod] != http.MethodPost {
		t.Fatal("requestType has wrong attributes")
	}

	log.WithFields(logData).Info("test")
}
