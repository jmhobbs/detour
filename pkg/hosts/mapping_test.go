package hosts

import (
	"bytes"
	"testing"
)

func TestMappingAdd(t *testing.T) {
	hm := HostMapping{}

	hm.Add("127.0.0.1", "example.com")
	if hm["example.com"] != "127.0.0.1" {
		t.Error("First Add Failed")
	}

	hm.Add("127.0.1.1", "example.com")
	if hm["example.com"] != "127.0.1.1" {
		t.Error("Second Add Failed")
	}
}

func TestMappingRemove(t *testing.T) {
	hm := HostMapping{}

	hm.Add("127.0.0.1", "example.com")
	hm.Add("127.0.0.1", "www.example.com")
	hm.Add("192.168.1.1", "local.example.com")

	hm.Remove("local.example.com")
	if hm["local.example.com"] != "" {
		t.Error("Failed to remove local.example.com")
	}
}

func TestMappingWrite(t *testing.T) {
	hm := HostMapping{}

	hm.Add("127.0.0.1", "example.com")
	hm.Add("127.0.0.1", "www.example.com")
	hm.Add("192.168.1.1", "local.example.com")

	buf := bytes.NewBuffer([]byte(""))
	err := hm.Write(buf)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := "192.168.1.1\tlocal.example.com\n127.0.0.1\texample.com  www.example.com\n"

	if buf.String() != expected {
		t.Errorf("Write() failed.\nGot:\n%v\nExpected:\n%v\n", buf.String(), expected)
	}
}
