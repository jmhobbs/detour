package hosts

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMappingAdd(t *testing.T) {
	hm := HostMapping{}

	hm.Add("127.0.0.1", "example.com")
	if !reflect.DeepEqual(hm["127.0.0.1"], []string{"example.com"}) {
		t.Error("First Add Failed")
	}

	hm.Add("127.0.0.1", "www.example.com")
	if !reflect.DeepEqual(hm["127.0.0.1"], []string{"example.com", "www.example.com"}) {
		t.Error("Second Add Failed")
	}
}

func TestMappingRemove(t *testing.T) {
	hm := HostMapping{}

	hm.Add("127.0.0.1", "example.com")
	hm.Add("127.0.0.1", "www.example.com")
	hm.Add("192.168.1.1", "local.example.com")

	if !reflect.DeepEqual(hm["127.0.0.1"], []string{"example.com", "www.example.com"}) || !reflect.DeepEqual(hm["192.168.1.1"], []string{"local.example.com"}) {
		t.Error("Setup Failed")
	}

	hm.Remove("local.example.com")
	if !reflect.DeepEqual(hm["192.168.1.1"], []string{}) {
		t.Error("Failed to remove local.example.com")
	}

	hm.Remove("example.com")
	if !reflect.DeepEqual(hm["127.0.0.1"], []string{"www.example.com"}) {
		t.Error("Failed to remove example.com")
	}
}

func TestMappingWrite(t *testing.T) {
	hm := HostMapping{}

	hm.Add("127.0.0.1", "example.com")
	hm.Add("127.0.0.1", "www.example.com")
	hm.Add("192.168.1.1", "local.example.com")

	if !reflect.DeepEqual(hm["127.0.0.1"], []string{"example.com", "www.example.com"}) || !reflect.DeepEqual(hm["192.168.1.1"], []string{"local.example.com"}) {
		t.Error("Setup Failed")
	}

	buf := bytes.NewBuffer([]byte(""))
	err := hm.Write(buf)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := "127.0.0.1\texample.com www.example.com\n192.168.1.1\tlocal.example.com\n"

	if buf.String() != expected {
		t.Errorf("Write() failed.\nGot:\n%v\nExpected:\n%v\n", buf.String(), expected)
	}
}
