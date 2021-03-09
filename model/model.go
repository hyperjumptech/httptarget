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

package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	EndPointsCapacity = 100
)

type EndPoint struct {
	ID            int                 `json:"id"`
	BasePath      string              `json:"base_path"`
	DelayMinMs    int                 `json:"delay_min_ms"`
	DelayMaxMs    int                 `json:"delay_max_ms"`
	ReturnCode    int                 `json:"return_code"`
	ReturnHeaders map[string][]string `json:"return_headers"`
	ReturnBody    string              `json:"return_body"`
}

type EndPoints struct {
	Map map[int]*EndPoint
}

func (e *EndPoints) List() []*EndPoint {
	ret := make([]*EndPoint, 0)
	for _, v := range e.Map {
		ret = append(ret, v)
	}
	return ret
}

func (e *EndPoints) Delete(id int) error {
	if _, ok := e.Map[id]; ok {
		delete(e.Map, id)
		return nil
	}
	return fmt.Errorf("notfound")
}

func (e *EndPoints) Update(id int, ep *EndPoint) error {
	if _, ok := e.Map[id]; ok {
		ep.ID = id
		if len(ep.BasePath) == 0 {
			return fmt.Errorf("basepath must not empty")
		}
		if strings.HasPrefix(ep.BasePath, "/docs") || strings.HasPrefix(ep.BasePath, "/api") {
			return fmt.Errorf("basepath \"%s\" is prefixed with server's internal use path \"/docs\" or \"/api\"", ep.BasePath)
		}
		if isPref, by := e.IsPrefixed(id, ep.BasePath); isPref {
			return fmt.Errorf("basepath \"%s\" is prefixed with other endpoint with basepath \"%s\"", ep.BasePath, by)
		}
		if ep.ReturnCode < 200 || ep.ReturnCode >= 600 {
			return fmt.Errorf("return code out of range")
		}
		if !strings.HasPrefix(ep.BasePath, "/") {
			ep.BasePath = "/" + ep.BasePath
		}
		if ep.DelayMinMs > ep.DelayMaxMs {
			t := ep.DelayMinMs
			ep.DelayMinMs = ep.DelayMaxMs
			ep.DelayMaxMs = t
		}
		e.Map[ep.ID] = ep
		return nil
	}
	return fmt.Errorf("notfound")
}

func (e *EndPoints) Add(ep *EndPoint) error {
	ep.ID = e.LastID() + 1
	if len(ep.BasePath) == 0 {
		return fmt.Errorf("basepath must not empty")
	}
	if strings.HasPrefix(ep.BasePath, "/docs") || strings.HasPrefix(ep.BasePath, "/api") {
		return fmt.Errorf("basepath \"%s\" is prefixed with server's internal use path \"/docs\" or \"/api\"", ep.BasePath)
	}
	if isPref, by := e.IsPrefixed(-1, ep.BasePath); isPref {
		return fmt.Errorf("basepath %s is prefixed with %s", ep.BasePath, by)
	}
	if ep.ReturnCode < 200 || ep.ReturnCode >= 600 {
		return fmt.Errorf("return code out of range")
	}
	if !strings.HasPrefix(ep.BasePath, "/") {
		ep.BasePath = "/" + ep.BasePath
	}
	if ep.DelayMinMs > ep.DelayMaxMs {
		t := ep.DelayMinMs
		ep.DelayMinMs = ep.DelayMaxMs
		ep.DelayMaxMs = t
	}
	e.Map[ep.ID] = ep
	logrus.Infof("Added test endpoint on [%s], code %d, minDelay %d ms, maxDelay %d ms", ep.BasePath, ep.ReturnCode, ep.DelayMinMs, ep.DelayMaxMs)
	if len(e.Map) > EndPointsCapacity {
		_ = e.Delete(e.FirstID())
	}
	return nil
}

func (e *EndPoints) IsPrefixed(id int, path string) (bool, string) {
	for _, ep := range e.Map {
		if strings.HasPrefix(path, ep.BasePath) && ep.ID != id {
			return true, ep.BasePath
		}
	}
	return false, ""
}

func (e *EndPoints) GetByPath(path string) *EndPoint {
	if e.Map == nil {
		e.Map = make(map[int]*EndPoint)
	}
	for _, ep := range e.Map {
		if strings.HasPrefix(path, ep.BasePath) {
			return ep
		}
	}
	return nil
}

func (e *EndPoints) LastID() int {
	if e.Map == nil {
		e.Map = make(map[int]*EndPoint)
	}
	ret := -1
	for i, _ := range e.Map {
		if i > ret {
			ret = i
		}
	}
	return ret
}

func (e *EndPoints) FirstID() int {
	if e.Map == nil {
		e.Map = make(map[int]*EndPoint)
	}
	ret := 1000000000
	for i, _ := range e.Map {
		if i < ret {
			ret = i
		}
	}
	return ret
}
