package test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/matobi/mam-go-lib/pkg/conf"
)

const (
	confPort = "port"
)

func TestEmpty(t *testing.T) {
	c := conf.NewConfig("")
	_, err := c.LogAndValidate()
	if err != nil {
		t.Errorf("empty cfg not validated")
	}
}

func TestAddGet(t *testing.T) {
	data := []struct {
		t     conf.ValueType
		name  string
		value string
	}{
		{conf.VtStr, "str1", "str1"},
		{conf.VtStr, "str2", "str2"},
		{conf.VtInt, "int1", "1"},
		{conf.VtInt, "int2", "123456789"},
		{conf.VtInt, "int3", "-12345"},
	}
	c := conf.NewConfig("")
	for _, d := range data {
		c.Add(d.t, d.name, d.value)
	}

	_, err := c.LogAndValidate()
	if err != nil {
		t.Errorf("unexpected failed validate")
	}

	for _, d := range data {
		value := c.Str(d.name)
		if d.value != value {
			t.Errorf("unexpected value; n=%s; v=%s; got=%s", d.name, d.value, value)
		}
		if d.t == conf.VtInt && toInt(d.value) != toInt(value) {
			t.Errorf("unexpected int value; n=%s; v=%d; got=%d", d.name, toInt(d.value), toInt(value))
		}
	}
}

func toInt(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Errorf("failed parse int; %s", s)
		return -1000000 - rand.Int63n(1000000)
	}
	return n
}

func TestOverwrite(t *testing.T) {
	name := "str1"
	value1 := "value1"
	value2 := "value2"
	c := conf.NewConfig("")

	c.Add(conf.VtStr, name, value1)
	v := c.Str(name)
	if v != value1 {
		t.Errorf("unexpected value; n=%s; v=%s; got=%s", name, value1, v)
	}

	c.Add(conf.VtStr, name, value2)
	v = c.Str(name)
	if v != value2 {
		t.Errorf("unexpected value; n=%s; v=%s; got=%s", name, value2, v)
	}
}

func TestProfile(t *testing.T) {
	c := conf.NewConfig("profActive")
	c.Add(conf.VtStr, "nameAll", "valueAll")
	c.AddProfile(conf.VtStr, "profActive", "name", "valueActive")
	c.AddProfile(conf.VtStr, "profInctive", "name", "valueInactive")

	if "valueAll" != c.Str("nameAll") {
		t.Errorf("value for global profile missing; exp=%s; got=%s", "valueAll", c.Str("nameAll"))
	}
	if "valueActive" != c.Str("name") {
		t.Errorf("wrong value for active profile; exp=%s; got=%s", "valueActive", c.Str("name"))
	}
}

func TestBadInt(t *testing.T) {
	c := conf.NewConfig("")
	c.Add(conf.VtInt, "int", "hello")
	_, err := c.LogAndValidate()
	if err == nil {
		t.Errorf("missing validate error")
	}
}

// todo:
//func TestDir(t *testing.T) {
//}
// todo:
//func TestFile(t *testing.T) {
//}
