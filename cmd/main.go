/*-----------------------------------------------------------------------------------
  --  HttpTarget                                                                       --
  --  Copyright (C) 2021  HttpTarget's Contributors                                    --
  --                                                                               --
  --  This program is free software: you can redistribute it and/or modify         --
  --  it under the terms of the GNU Affero General Public License as published     --
  --  by the Free Software Foundation, either version 3 of the License, or         --
  --  (at your option) any later version.                                          --
  --                                                                               --
  --  This program is distributed in the hope that it will be useful,              --
  --  but WITHOUT ANY WARRANTY; without even the implied warranty of               --
  --  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the                --
  --  GNU Affero General Public License for more details.                          --
  --                                                                               --
  --  You should have received a copy of the GNU Affero General Public License     --
  --  along with this program.  If not, see <https:   -- www.gnu.org/licenses/>.   --
  -----------------------------------------------------------------------------------*/

package main

import (
	"flag"
	"github.com/hyperjumptech/httptarget/model"
	"github.com/hyperjumptech/httptarget/server"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	portFlag     = flag.Int("p", 51423, "Listen port")
	hostFlag     = flag.String("h", "0.0.0.0", "Bind host")
	bodyFlag     = flag.String("body", "OK", "HTTP response body")
	pathFlag     = flag.String("path", "/", "Base path")
	codeFlag     = flag.Int("code", 200, "Response code")
	minDelayFlag = flag.Int("minDelay", 0, "Minimum Delay Millisecond")
	maxDelayFlag = flag.Int("maxDelay", 200, "Maximum Delay Millisecond")
	helpFlag     = flag.Bool("help", false, "Display this usage message")
)

func main() {
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	logrus.SetLevel(logrus.TraceLevel)

	initEp := &model.EndPoint{
		ID:            0,
		BasePath:      *pathFlag,
		DelayMinMs:    *minDelayFlag,
		DelayMaxMs:    *maxDelayFlag,
		ReturnCode:    *codeFlag,
		ReturnHeaders: map[string][]string{"Content-Type": {"text/plain"}},
		ReturnBody:    *bodyFlag,
	}

	err := server.Start(*hostFlag, *portFlag, initEp)
	if err != nil {
		logrus.Error(err)
	}
}
