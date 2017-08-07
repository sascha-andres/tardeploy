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
