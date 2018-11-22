// Copyright 2018 Andrew E. Bruno
//
// This file is part of myoxi.
//
// myoxi is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// myoxi is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with myoxi.  If not, see <http://www.gnu.org/licenses/>.

package device

import (
	"io"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

type MockCMS50 struct {
	command uint8
	counter int
}

func (c *MockCMS50) Read(p []byte) (int, error) {
	var res []byte
	switch c.command {
	case CommandHello1:
		res = []byte{0xc, 0x80}
	case CommandHello2:
		res = []byte{0xc, 0x80}
	case CommandGetUserInfo:
		res = []byte{0x05, 0x80, 0x80, 0xf5, 0xf3, 0xe5, 0xf2, 0x80, 0x80}
	case CommandGetOximeterModel:
		res = []byte{0x02, 0x80, 0x80, 0xb5, 0xb0, 0xc6, 0xa0, 0xa0, 0xa0, 0x02, 0x81, 0xff, 0xa0, 0x80, 0x80, 0x80, 0x80, 0x80}
	case CommandGetSessionCount:
		res = []byte{0x0a, 0x80, 0x80, 0x81}
	case CommandGetSessionDuration:
		res = []byte{0x08, 0x88, 0x80, 0x80, 0xfc, 0xca, 0x80, 0x80}
	case CommandGetSessionTime:
		if c.counter == 0 {
			res = []byte{0x07, 0x80, 0x80, 0x80, 0x94, 0x92, 0x8b, 0x92}
		} else {
			res = []byte{0x12, 0x0, 0x0, 0x0, 0x0, 0xb, 0x27, 0x0}
		}
	case CommandGetSessionData:
		res = []byte{
			0x0f, 0x80, 0xe2, 0xcd, 0xe2, 0xcc, 0xe1, 0xcd,
			0x0f, 0x00, 0x62, 0x4c, 0x62, 0x4a, 0x62, 0x4a,
			0x0f, 0x00, 0x62, 0x48, 0x62, 0x47, 0x62, 0x46,
			0x0f, 0x80, 0x92, 0x80, 0x8b, 0xa7, 0xe1, 0xc4,
			0x0f, 0x80, 0xe1, 0xc5, 0xe0, 0xc5, 0xe0, 0xc5,
		}
	}

	copy(p, res)
	c.counter++

	return len(res), io.EOF
}

func (c *MockCMS50) Write(p []byte) (int, error) {
	c.counter = 0
	c.command = p[2]
	return len(p), nil
}

func newTestDevice() *CMS50 {
	cms := &CMS50{}
	cms.device = &MockCMS50{}

	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	return cms
}

func TestReset(t *testing.T) {
	cms := newTestDevice()
	err := cms.ResetDevice()
	if err != nil {
		t.Error(err)
	}
}

func TestGetUser(t *testing.T) {
	cms := newTestDevice()
	user, err := cms.GetUser()
	if err != nil {
		t.Error(err)
	}

	if user != "user" {
		t.Errorf("Invalid user: got '% x' should be '% x'", user, "user")
	}
}

func TestGetModel(t *testing.T) {
	cms := newTestDevice()
	model, err := cms.GetModel()
	if err != nil {
		t.Error(err)
	}

	if model != "50F" {
		t.Errorf("Invalid model: got '%s' should be '%s'", model, "50F")
	}
}

func TestGetSessionCount(t *testing.T) {
	cms := newTestDevice()
	count, err := cms.GetSessionCount()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Errorf("Invalid session count: got '%d' should be '%d'", count, 1)
	}
}

func TestGetSessionDuration(t *testing.T) {
	cms := newTestDevice()
	duration, err := cms.GetSessionDuration(1)
	if err != nil {
		t.Error(err)
	}

	validDuration := time.Duration(25918) * time.Second

	if duration != validDuration {
		t.Errorf("Invalid session duration: got '%s' should be '%s'", duration, validDuration)
	}
}

func TestGetSessionTime(t *testing.T) {
	cms := newTestDevice()
	dateTime, err := cms.GetSessionTime(1)
	if err != nil {
		t.Error(err)
	}

	validDateTime := time.Date(2018, time.November, 18, 0, 11, 39, 0, time.Local)

	if dateTime != validDateTime {
		t.Errorf("Invalid session time: got '%s' should be '%s'", dateTime, validDateTime)
	}
}

func TestGetSessionData(t *testing.T) {
	cms := newTestDevice()
	data, err := cms.GetSessionData(0)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 15 {
		t.Errorf("Invalid session data: got '%d' should be '%d'", len(data), 3)
	}
}
