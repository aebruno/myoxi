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
	"time"

	"github.com/aebruno/myoxi/model"
)

type Device interface {
	Connect(port string) error
	ResetDevice() error
	GetModel() (string, error)
	GetSessionCount() (uint8, error)
	GetSessionDuration(session uint8) (time.Duration, error)
	GetSessionTime(session uint8) (time.Time, error)
	GetSessionData(session uint8) ([]*model.OxiRecord, error)
	GetUser() (string, error)
}
