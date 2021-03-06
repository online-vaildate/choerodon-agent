package url

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type URLTestSuite struct {
	suite.Suite
}

func (suite *URLTestSuite) TestParseURL() {
	tests := []struct {
		input string
		want  string
	}{
		{"ws://localhost:8060/agent/?key=env:choerodon-agent-test", "ws://localhost:8060/agent/log?key=env:choerodon-agent-test"},
		{"ws://localhost:8060/agent/", "ws://localhost:8060/agent/log"},
	}
	for _, test := range tests {
		base, _ := url.Parse(test.input)
		newURL, err := ParseURL(base, "log")
		suite.Nil(err, "no error parse url")
		suite.Equal(test.want, newURL.String(), "error parse url")
	}
}

func TestURLTestSuite(t *testing.T) {
	suite.Run(t, new(URLTestSuite))
}
