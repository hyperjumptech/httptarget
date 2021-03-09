package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type EndPoint struct {
	ID            int
	BasePath      string
	DelayMinMs    int
	DelayMaxMs    int
	ReturnCode    int
	ReturnHeaders map[string][]string
	ReturnBody    string
}

type EndPoints struct {
	List []*EndPoint
}

func (e *EndPoints) Add(ep *EndPoint) error {
	ep.ID = e.LastID() + 1
	if len(ep.BasePath) == 0 {
		return fmt.Errorf("basepath must not empty")
	}
	if !strings.HasPrefix(ep.BasePath, "/") {
		ep.BasePath = "/" + ep.BasePath
	}
	if ep.DelayMinMs > ep.DelayMaxMs {
		t := ep.DelayMinMs
		ep.DelayMinMs = ep.DelayMaxMs
		ep.DelayMaxMs = t
	}
	e.List = append(e.List, ep)
	logrus.Infof("Added test endpoint on %s , code %d, minDelay %d, maxDelay %d", ep.BasePath, ep.ReturnCode, ep.DelayMinMs, ep.DelayMaxMs)
	return nil
}

func (e *EndPoints) GetByPath(path string) *EndPoint {
	if e.List == nil {
		e.List = make([]*EndPoint, 0)
	}
	for _, ep := range e.List {
		if strings.HasPrefix(path, ep.BasePath) {
			return ep
		}
	}
	return nil
}

func (e *EndPoints) LastID() int {
	if e.List == nil {
		e.List = make([]*EndPoint, 0)
	}
	ret := -1
	for _, ep := range e.List {
		if ep.ID > ret {
			ret = ep.ID
		}
	}
	return ret
}
