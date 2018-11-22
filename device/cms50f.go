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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aebruno/myoxi/model"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

const (
	CommandHello1              = 0xa7
	CommandHello2              = 0xa2
	CommandGetSessionCount     = 0xa3
	CommandGetSessionTime      = 0xa5
	CommandGetSessionDuration  = 0xa4
	CommandGetUserInfo         = 0xab
	CommandGetSessionData      = 0xa6
	CommandGetOximeterDeviceid = 0xaa
	CommandGetOximeterInfo     = 0xb0
	CommandGetOximeterModel    = 0xa8
	CommandGetOximeterVendor   = 0xa9
	CommandSessionErase        = 0xae
	DurationDivisor            = 2
)

type CMS50 struct {
	device       io.ReadWriter
	model        string
	user         string
	sessionCount uint8
}

func (c *CMS50) makeCommand() []byte {
	return []byte{0x7d, 0x81, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
}

func (c *CMS50) readBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	read, err := c.device.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if read == 0 {
		return nil, fmt.Errorf("Read 0 bytes from device. Is it turned on?")
	}

	return buf[:read], nil
}

func (c *CMS50) execCommand(command uint8) error {
	cmd := c.makeCommand()
	cmd[2] |= (command & 0x7f)

	log.Debugf("Send Command: % x", cmd)

	_, err := c.device.Write(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *CMS50) execCommandWithArg(command, arg uint8) error {
	cmd := c.makeCommand()
	cmd[2] |= (command & 0x7f)
	cmd[4] |= (arg & 0x7f)

	log.Debugf("Send Command With Arg: % x", cmd)

	_, err := c.device.Write(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *CMS50) Connect(port string) error {
	conf := &serial.Config{
		Name:        port,
		Baud:        115200,
		ReadTimeout: time.Second * 5,
	}

	dev, err := serial.OpenPort(conf)
	if err != nil {
		return err
	}

	c.device = dev
	return nil
}

func (c *CMS50) ResetDevice() error {
	err := c.execCommand(CommandHello1)
	if err != nil {
		return err
	}

	res, err := c.readBytes(8)
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	if res[0] != 0xc {
		return fmt.Errorf("Unknown result for CommandHello1: % x", res)
	}

	err = c.execCommand(CommandHello2)
	if err != nil {
		return err
	}

	res, err = c.readBytes(8)
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	if res[0] != 0xc {
		return fmt.Errorf("Unknown result for CommandHello2: % x", res)
	}

	return nil
}

func (c *CMS50) GetUser() (string, error) {
	if len(c.user) > 0 {
		return c.user, nil
	}

	err := c.execCommand(CommandGetUserInfo)
	if err != nil {
		return "", err
	}

	res, err := c.readBytes(100)
	if err != nil {
		return "", err
	}

	if res[0] != 0x05 {
		return "", fmt.Errorf("Unknown result for CommandGetUserInfo: % x", res)
	}

	log.Debugf("Received %d bytes for user: % x", len(res), res)

	for i := 3; i < len(res); i++ {
		res[i] ^= 0x80
	}

	user := bytes.TrimRightFunc(res[3:], func(r rune) bool {
		if r == 0x00 {
			return true
		}
		return false
	})

	c.user = strings.TrimSpace(string(user))

	return c.user, nil
}

func (c *CMS50) GetModel() (string, error) {
	if len(c.model) > 0 {
		return c.model, nil
	}

	err := c.execCommand(CommandGetOximeterModel)
	if err != nil {
		return "", err
	}

	res, err := c.readBytes(100)
	if err != nil {
		return "", err
	}

	if res[0] != 0x02 {
		return "", fmt.Errorf("Unknown result for CommandGetOximeterModel: % x", res)
	}

	log.Debugf("Received %d bytes for model string: % x", len(res), res)

	for i := 3; i < len(res); i++ {
		res[i] ^= 0x80
	}

	c.model = strings.TrimSpace(string(res[3:8]))

	return c.model, nil
}

func (c *CMS50) GetSessionCount() (uint8, error) {
	err := c.execCommand(CommandGetSessionCount)
	if err != nil {
		return 0, err
	}

	res, err := c.readBytes(8)
	if err != nil {
		return 0, err
	}

	if res[0] != 0x0a {
		return 0, fmt.Errorf("Unknown result for CommandGetSessionCount: % x", res)
	}

	log.Debugf("Received %d bytes for session count: % x", len(res), res)

	c.sessionCount = res[3] ^ 0x80

	return c.sessionCount, nil
}

func (c *CMS50) GetSessionDuration(session uint8) (time.Duration, error) {
	err := c.execCommandWithArg(CommandGetSessionDuration, session)
	if err != nil {
		return 0, err
	}

	res, err := c.readBytes(8)
	if err != nil {
		return 0, err
	}

	if res[0] != 0x08 {
		return 0, fmt.Errorf("Unknown result for CommandGetSessionDuration: % x", res)
	}

	log.Debugf("Received %d bytes for session duration: % x", len(res), res)

	if len(res) < 7 {
		return 0, fmt.Errorf("Not enough bytes returned for CommandGetSessionDuration. Need 7 got %d", len(res))
	}

	for i := 1; i < 7; i++ {
		res[i] ^= 0x80
	}

	seconds := (uint16(res[1]) & 0x4) << 5
	seconds |= uint16(res[4])
	seconds |= (uint16(res[5]) | ((uint16(res[1]) & 0x8) << 4)) << 8
	seconds |= (uint16(res[6]) | ((uint16(res[1]) & 0x10) << 3)) << 16

	log.Debugf("Session duration is %d / %d", seconds, DurationDivisor)

	duration := time.Duration(seconds/DurationDivisor) * time.Second

	log.Debugf("Session %d has duration of %s (%.1fs)", session, duration, duration.Seconds())

	return duration, nil
}

func (c *CMS50) GetSessionTime(session uint8) (time.Time, error) {
	err := c.execCommandWithArg(CommandGetSessionTime, session)
	if err != nil {
		return time.Time{}, err
	}

	dateRes, err := c.readBytes(8)
	if err != nil {
		return time.Time{}, err
	}

	timeRes, err := c.readBytes(8)
	if err != nil {
		return time.Time{}, err
	}

	if len(dateRes) != 8 {
		return time.Time{}, fmt.Errorf("Not enough bytes returned for date CommandGetSessionTime. Need 8 got %d", len(dateRes))
	}

	if len(timeRes) != 8 {
		return time.Time{}, fmt.Errorf("Not enough bytes returned for time CommandGetSessionTime. Need 8 got %d", len(timeRes))
	}

	if dateRes[0] != 0x07 {
		return time.Time{}, fmt.Errorf("Unknown result for date CommandGetSessionTime: % x", dateRes)
	}
	if timeRes[0] != 0x12 {
		return time.Time{}, fmt.Errorf("Unknown result for time CommandGetSessionTime: % x", timeRes)
	}

	log.Debugf("Received %d bytes for session date: % x", len(dateRes), dateRes)
	log.Debugf("Received %d bytes for session time: % x", len(timeRes), timeRes)

	date := make([]int, 8)
	tim := make([]int, 8)

	for i := 4; i < 8; i++ {
		date[i] = int(dateRes[i]) & ^0x80
		tim[i] = int(timeRes[i]) & ^0x80
	}

	year, _ := strconv.Atoi(fmt.Sprintf("%d%d", date[4], date[5]))

	dateTime := time.Date(year, time.Month(date[6]), date[7], tim[4], tim[5], tim[6], tim[7], time.Local)

	log.Debugf("Session %d start time is: %s", session, dateTime)

	return dateTime, nil
}

func (c *CMS50) GetSessionData(session uint8) ([]*model.OxiRecord, error) {
	err := c.ResetDevice()
	if err != nil {
		return nil, err
	}

	err = c.execCommandWithArg(CommandGetSessionData, session)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(c.device)
	data := make([]*model.OxiRecord, 0)

	for {
		buf := make([]byte, 8)

		_, err := io.ReadFull(reader, buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if buf[0] != 0x0f {
			return nil, fmt.Errorf("Unknown result for CommandGetSessionData: % x", buf)
		}

		data = append(data, c.newOxiRecords(buf)...)
	}

	return data, nil
}

func (c *CMS50) newOxiRecord(pulse, spo2 uint8) *model.OxiRecord {
	if pulse == 0xff {
		return &model.OxiRecord{Pulse: 0, Spo2: 0}
	}

	return &model.OxiRecord{Pulse: pulse, Spo2: spo2}
}
func (c *CMS50) newOxiRecords(buf []byte) []*model.OxiRecord {
	data := make([]*model.OxiRecord, 3)
	msb := buf[1]
	for i := 2; i < len(buf); i++ {
		buf[i] &= 0x7f
		if msb&0x01 != 0 {
			buf[i] |= 0x80
		} else {
			buf[i] |= 0
		}
		msb >>= 1
	}

	data[0] = c.newOxiRecord(buf[3], buf[2])
	data[1] = c.newOxiRecord(buf[5], buf[4])
	data[2] = c.newOxiRecord(buf[7], buf[6])

	log.Debugf("Got new OxiRecords packet: %s", data)

	return data
}
