package gollab_test

import (
	"github.com/danielslee/gollab"
	"github.com/danielslee/gollab/runetoken"
	"io/ioutil"
	"strings"
	"testing"
)

func TestInsertApply(t *testing.T) {
	i := gollab.Insert{Tokens: runetoken.Array("hello")}
	var w runetoken.StringWriter
	if err := i.Apply(nil, &w); err != nil {
		t.Error(err)
		return
	}

	if w.String() != "hello" {
		t.Errorf("expected 'hello', got '%s'", w.String())
	}
}


func TestInsertApplyUnicode(t *testing.T) {
	i := gollab.Insert{Tokens: runetoken.Array("안녕하세요")}
	var w runetoken.StringWriter
	if err := i.Apply(nil, &w); err != nil {
		t.Error(err)
		return
	}

	if w.String() != "안녕하세요" {
		t.Errorf("expected '안녕하세요', got '%s'", w.String())
	}
}


func TestDeleteApply(t *testing.T) {
	d := gollab.Delete{Count: 2}
	r := runetoken.StringReader{Reader: strings.NewReader("hello")}
	if err := d.Apply(r, nil); err != nil {
		t.Error(err)
		return
	}

	if all, err := ioutil.ReadAll(r.Reader); err != nil {
		t.Error(err)
		return
	} else if string(all) != "llo" {
		t.Errorf("expected text left in reader to be 'llo', got '%s'",
			string(all))
	}
}


func TestDeleteApplyUnicode(t *testing.T) {
	d := gollab.Delete{Count: 2}
	r := runetoken.StringReader{Reader: strings.NewReader("안녕하세요")}
	if err := d.Apply(r, nil); err != nil {
		t.Error(err)
		return
	}

	if all, err := ioutil.ReadAll(r.Reader); err != nil {
		t.Error(err)
		return
	} else if string(all) != "하세요" {
		t.Errorf("expected text left in reader to be '하세요', got '%s'",
			string(all))
	}
}


func TestRetainApply(t *testing.T) {
	retain := gollab.Retain{Count: 2}
	r := runetoken.StringReader{Reader: strings.NewReader("hello")}
	var w runetoken.StringWriter
	if err := retain.Apply(r, &w); err != nil {
		t.Error(err)
		return
	}

	if w.String() != "he" {
		t.Errorf("expected 'he', got '%s'", w.String())
	}

	if all, err := ioutil.ReadAll(r.Reader); err != nil {
		t.Error(err)
		return
	} else if string(all) != "llo" {
		t.Errorf("expected text left in reader to be 'llo' got '%s'",
			string(all))
	}
}


func TestRetainApplyUnicode(t *testing.T) {
	retain := gollab.Retain{Count: 2}
	r := runetoken.StringReader{Reader: strings.NewReader("안녕하세요")}
	var w runetoken.StringWriter
	if err := retain.Apply(r, &w); err != nil {
		t.Error(err)
		return
	}

	if w.String() != "안녕" {
		t.Errorf("expected '안녕', got '%s'", w.String())
	}

	if all, err := ioutil.ReadAll(r.Reader); err != nil {
		t.Error(err)
		return
	} else if string(all) != "하세요" {
		t.Errorf("expected text left in reader to be '하세요' got '%s'",
			string(all))
	}
}
