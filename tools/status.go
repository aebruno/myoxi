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

package tools

import (
	"fmt"

	"github.com/aebruno/myoxi/device"
)

func Status(device device.Device) error {
	err := device.ResetDevice()
	if err != nil {
		return fmt.Errorf("Failed to reset device: %s", err)
	}

	user, err := device.GetUser()
	if err != nil {
		return fmt.Errorf("Failed to get user: %s", err)
	}

	model, err := device.GetModel()
	if err != nil {
		return fmt.Errorf("Failed to get device model: %s", err)
	}

	count, err := device.GetSessionCount()
	if err != nil {
		return fmt.Errorf("Failed to get session count: %s", err)
	}

	fmt.Printf("Device model: %s\n", model)
	fmt.Printf("Userinfo: %s\n", user)
	fmt.Printf("Session count: %d\n", count)
	fmt.Printf("------------------------------\n")

	for i := uint8(0); i < count; i++ {
		duration, err := device.GetSessionDuration(i)
		if err != nil {
			return fmt.Errorf("Failed to fetch session duration: %s", err)
		}

		startTime, err := device.GetSessionTime(i)
		if err != nil {
			return fmt.Errorf("Failed to fetch session time: %s", err)
		}

		fmt.Printf("- %s (%s)\n", startTime, duration)
	}

	return nil
}
