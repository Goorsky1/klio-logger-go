package logger_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/g2a-com/klio-logger-go"
)

func TestParseLevel(t *testing.T) {
	var l log.Level
	var ok bool

	l, ok = log.ParseLevel("spam")
	assert.Equal(t, log.SpamLevel, l)
	assert.Equal(t, true, ok)

	l, ok = log.ParseLevel("debug")
	assert.Equal(t, log.DebugLevel, l)
	assert.Equal(t, true, ok)

	l, ok = log.ParseLevel("verbose")
	assert.Equal(t, log.VerboseLevel, l)
	assert.Equal(t, true, ok)

	l, ok = log.ParseLevel("info")
	assert.Equal(t, log.InfoLevel, l)
	assert.Equal(t, true, ok)

	l, ok = log.ParseLevel("warn")
	assert.Equal(t, log.WarnLevel, l)
	assert.Equal(t, true, ok)

	l, ok = log.ParseLevel("error")
	assert.Equal(t, log.ErrorLevel, l)
	assert.Equal(t, true, ok)

	l, ok = log.ParseLevel("fatal")
	assert.Equal(t, log.FatalLevel, l)
	assert.Equal(t, true, ok)

	l, ok = log.ParseLevel("SPAM")
	assert.Equal(t, log.SpamLevel, l)
	assert.Equal(t, true, ok)

	l, ok = log.ParseLevel("unknown")
	assert.Equal(t, log.DefaultLevel, l)
	assert.Equal(t, false, ok)

	l, ok = log.ParseLevel("  spam")
	assert.Equal(t, log.DefaultLevel, l)
	assert.Equal(t, false, ok)
}

func TestNew(t *testing.T) {
	var b bytes.Buffer
	l := log.New(&b)
	assert.IsType(t, &log.Logger{}, l)
}

func TestWithLevel(t *testing.T) {
	var b bytes.Buffer

	l1 := log.New(&b)
	l2 := l1.WithLevel(log.SpamLevel)
	l3 := l2.WithLevel(log.WarnLevel)

	assert.Equal(t, log.DefaultLevel, l1.Level())
	assert.Equal(t, log.SpamLevel, l2.Level())
	assert.Equal(t, log.WarnLevel, l3.Level())
}

func TestWithTags(t *testing.T) {
	var b bytes.Buffer

	l1 := log.New(&b)
	l2 := l1.WithTags("a", "b")
	l3 := l2.WithTags()

	l2.Tags()[0] = "xyz" // shouldn't affect logger tags

	assert.Equal(t, []string{}, l1.Tags())
	assert.Equal(t, []string{"a", "b"}, l2.Tags())
	assert.Equal(t, []string{}, l3.Tags())
}

func TestPrint(t *testing.T) {
	var b bytes.Buffer

	b.Reset()
	log.New(&b).Print("foo")
	assert.Equal(t, "\033_klio_log_level \"info\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())

	b.Reset()
	log.New(&b).WithTags("a", "b", "c").WithLevel(log.SpamLevel).Print("foo")
	assert.Equal(t, "\033_klio_log_level \"spam\"\033\\\033_klio_tags [\"a\",\"b\",\"c\"]\033\\foo\033_klio_reset\033\\\n", b.String())

	b.Reset()
	log.New(&b).WithTags("a", "b", "c").WithLevel(log.SpamLevel).WithLevel(log.DefaultLevel).WithTags().Print("foo")
	assert.Equal(t, "\033_klio_log_level \"info\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())

	b.Reset()
	log.New(&b).WithTags("\033\\").WithLevel(log.Level("\"")).Print("foo")
	assert.Equal(t, "\033_klio_log_level \"\\\"\"\033\\\033_klio_tags [\"\\u001b\\\\\"]\033\\foo\033_klio_reset\033\\\n", b.String())
}

func TestPrintf(t *testing.T) {
	var b bytes.Buffer

	b.Reset()
	log.New(&b).Printf("%s", "foo")
	assert.Equal(t, "\033_klio_log_level \"info\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
}

func TestSetOutput(t *testing.T) {
	var b1 bytes.Buffer
	var b2 bytes.Buffer

	l := log.New(&b1)
	l.SetOutput(&b2)
	l.Print("foo")

	assert.Equal(t, "", b1.String())
	assert.Equal(t, "\033_klio_log_level \"info\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b2.String())
}

func TestStandardLogger(t *testing.T) {
	l := log.StandardLogger()

	assert.Equal(t, os.Stdout, l.Output())
	assert.Equal(t, log.Level("info"), l.Level())
	assert.Equal(t, []string{}, l.Tags())
}

func TestErrorLogger(t *testing.T) {
	l := log.ErrorLogger()

	assert.Equal(t, os.Stderr, l.Output())
	assert.Equal(t, log.Level("error"), l.Level())
	assert.Equal(t, []string{}, l.Tags())
}

func TestConvenienceFunctions(t *testing.T) {
	var b bytes.Buffer

	log.StandardLogger().SetOutput(&b)
	defer log.StandardLogger().SetOutput(os.Stdout)

	b.Reset()
	log.Spam("foo")
	assert.Equal(t, "\033_klio_log_level \"spam\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Debug("foo")
	assert.Equal(t, "\033_klio_log_level \"debug\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Verbose("foo")
	assert.Equal(t, "\033_klio_log_level \"verbose\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Info("foo")
	assert.Equal(t, "\033_klio_log_level \"info\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Warn("foo")
	assert.Equal(t, "\033_klio_log_level \"warn\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Error("foo")
	assert.Equal(t, "\033_klio_log_level \"error\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Fatal("foo")
	assert.Equal(t, "\033_klio_log_level \"fatal\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Spamf("%s", "foo")
	assert.Equal(t, "\033_klio_log_level \"spam\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Debugf("%s", "foo")
	assert.Equal(t, "\033_klio_log_level \"debug\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Verbosef("%s", "foo")
	assert.Equal(t, "\033_klio_log_level \"verbose\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Infof("%s", "foo")
	assert.Equal(t, "\033_klio_log_level \"info\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Warnf("%s", "foo")
	assert.Equal(t, "\033_klio_log_level \"warn\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Errorf("%s", "foo")
	assert.Equal(t, "\033_klio_log_level \"error\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
	b.Reset()
	log.Fatalf("%s", "foo")
	assert.Equal(t, "\033_klio_log_level \"fatal\"\033\\\033_klio_tags []\033\\foo\033_klio_reset\033\\\n", b.String())
}