package logging_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/RichardKnop/logging"
	"github.com/stretchr/testify/assert"
)

func TestColouredFormatter(t *testing.T) {
	var (
		out, errOut = bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
		logger      = logging.New(out, errOut, new(logging.ColouredFormatter))
		now         time.Time
		actual      []byte
		expected    string
		err         error
	)

	// Test logger.Info
	now = time.Now()
	logger.Info("Test logger.Info")
	actual, err = ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"\x1b[0;94mINFO: %s coloured_formatter_test.go:27 Test logger.Info\x1b[0m\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Infof
	now = time.Now()
	logger.Infof("Test %s.%s", "logger", "Infof")
	actual, err = ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"\x1b[0;94mINFO: %s coloured_formatter_test.go:40 Test logger.Infof\x1b[0m\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Warning
	now = time.Now()
	logger.Warning("Test logger.Warning")
	actual, err = ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"\x1b[0;95mWARNING: %s coloured_formatter_test.go:53 Test logger.Warning\x1b[0m\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Warningf
	now = time.Now()
	logger.Warningf("Test %s.%s", "logger", "Warningf")
	actual, err = ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"\x1b[0;95mWARNING: %s coloured_formatter_test.go:66 Test logger.Warningf\x1b[0m\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Error
	now = time.Now()
	logger.Error("Test logger.Error")
	actual, err = ioutil.ReadAll(errOut)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"\x1b[0;91mERROR: %s coloured_formatter_test.go:79 Test logger.Error\x1b[0m\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Errorf
	now = time.Now()
	logger.Errorf("Test %s.%s", "logger", "Errorf")
	actual, err = ioutil.ReadAll(errOut)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"\x1b[0;91mERROR: %s coloured_formatter_test.go:92 Test logger.Errorf\x1b[0m\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Fatal
	now = time.Now()
	logger.Fatal("Test logger.Fatal")
	actual, err = ioutil.ReadAll(errOut)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"\x1b[0;91mFATAL: %s coloured_formatter_test.go:105 Test logger.Fatal\x1b[0m\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))

	// Test logger.Fatalf
	now = time.Now()
	logger.Fatalf("Test %s.%s", "logger", "Fatalf")
	actual, err = ioutil.ReadAll(errOut)
	if err != nil {
		log.Fatal(err)
	}
	expected = fmt.Sprintf(
		"\x1b[0;91mFATAL: %s coloured_formatter_test.go:118 Test logger.Fatalf\x1b[0m\n",
		now.Format("2006/01/02 15:04:05"),
	)
	assert.Equal(t, expected, string(actual))
}
