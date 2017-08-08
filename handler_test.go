// Copyright © 2017 Sascha Andres <sascha.andres@outlook.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tardeploy

import "testing"

func TestMakeApplicationValidTgz(t *testing.T) {
	app, err := makeApplication("test.tgz")
	if err != nil {
		t.Errorf("Unexpected error %#v", err)
		t.Failed()
	}
	if app != "test" {
		t.Errorf("Expected [test], got %s", app)
		t.Failed()
	}
}

func TestMakeApplicationValidTarGz(t *testing.T) {
	app, err := makeApplication("test.tar.gz")
	if err != nil {
		t.Errorf("Unexpected error %#v", err)
		t.Failed()
	}
	if app != "test" {
		t.Errorf("Expected [test], got %s", app)
		t.Failed()
	}
}

func TestMakeApplicationNoValidNameTgz(t *testing.T) {
	_, err := makeApplication(".tgz")
	if nil == err {
		t.Error("Expected error")
		t.Failed()
	}
}

func TestMakeApplicationNoValidNameTarGz(t *testing.T) {
	_, err := makeApplication(".tar.gz")
	if nil == err {
		t.Error("Expected error")
		t.Failed()
	}
}
