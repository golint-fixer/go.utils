package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	cases := []struct {
		src   string
		dst   string
		isnil bool
		rm    bool
		indir bool
	}{
		{os.Args[0], os.Args[0], true, false, false},
		{os.Args[0], "temp_file", true, true, false},
		{"not_exit", "temp_file", false, false, false},
		{os.Args[0], "", false, false, false},
		{os.Args[0], ".", true, true, true},
	}
	for i, cas := range cases {
		var srcnt []byte
		var err error
		if cas.isnil {
			if srcnt, err = ioutil.ReadFile(cas.src); err != nil {
				t.Errorf("want %v==nil (%d)", err, i)
			}
		}
		var srcfi os.FileInfo
		if cas.isnil {
			if srcfi, err = os.Stat(cas.src); err != nil {
				t.Errorf("want %v==nil (%d)", err, i)
			}
		}
		err = CopyFile(cas.dst, cas.src)
		if (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
		if err != nil {
			continue
		}
		dst := cas.dst
		if cas.indir {
			dst = filepath.Join(cas.dst, filepath.Base(cas.src))
		}
		dstfi, err := os.Stat(dst)
		if err != nil {
			t.Errorf("want %v==nil (%d)", err, i)
		}
		if d, s := dstfi.Mode(), srcfi.Mode(); d != s {
			t.Errorf("want %d=%d (%d)", d, s, i)
		}
		dstcnt, err := ioutil.ReadFile(dst)
		if err != nil {
			t.Errorf("want %v==nil (%d)", err, i)
		}
		if string(srcnt) != string(dstcnt) {
			t.Errorf("want %q=%q (%d)", string(srcnt), string(dstcnt), i)
		}
		if cas.rm {
			os.Remove(dst)
		}
	}
}

func TestExists(t *testing.T) {
	cases := []struct {
		path  string
		ok    bool
		isnil bool
	}{
		{os.Args[0], true, true}, {"not_existing23", false, true},
	}
	for i, cas := range cases {
		ok, err := Exists(cas.path)
		if (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
		if err != nil {
			continue
		}
		if ok != cas.ok {
			t.Errorf("want %t==%t (%d)", ok, cas.ok, i)
		}
	}
}

func TestIsDir(t *testing.T) {
	cases := []struct {
		path  string
		ok    bool
		isnil bool
	}{
		{os.Args[0], false, true}, {filepath.Dir(os.Args[0]), true, true},
		{"not_exist", false, false},
	}
	for i, cas := range cases {
		ok, err := IsDir(cas.path)
		if (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
		if err != nil {
			continue
		}
		if ok != cas.ok {
			t.Errorf("want %t==%t (%d)", ok, cas.ok, i)
		}
	}
}

func TestIsTheSame(t *testing.T) {
	cases := []struct {
		lf    string
		rf    string
		ok    bool
		isnil bool
	}{
		{os.Args[0], os.Args[0], true, true},
		{"", ".", true, true},
		{"dd", "u", false, true},
		{"dir", filepath.Join("dir", "..", "dir"), true, true},
	}
	for i, cas := range cases {
		ok, err := IsTheSame(cas.lf, cas.rf)
		if (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
		if err != nil {
			continue
		}
		if ok != cas.ok {
			t.Errorf("want %t==%t for %q and %q (%d)", ok, cas.ok, cas.lf, cas.rf, i)
		}
	}
}
