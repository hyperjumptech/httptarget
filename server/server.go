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

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hyperjumptech/httptarget/model"
	"github.com/hyperjumptech/httptarget/static"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	endPoints = &model.EndPoints{}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func AddEndPoint(ep *model.EndPoint) error {
	return endPoints.Add(ep)
}

func Start(host string, port int, initEndpoint *model.EndPoint) error {
	err := endPoints.Add(initEndpoint)
	if err != nil {
		return err
	}
	listen := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{
		Addr: listen,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:      5 * time.Minute,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       5 * time.Second,
		IdleTimeout:       2 * time.Second,
		Handler:           &HttpTargetHandler{},
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Infof("Server listening at %s", listen)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	os.Exit(0)
	return nil
}

type HttpTargetHandler struct {
}

func (h *HttpTargetHandler) CreatePath(res http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
	} else {
		ep := &model.EndPoint{}
		err = json.Unmarshal(bodyBytes, ep)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(err.Error()))
		} else {
			err = endPoints.Add(ep)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				res.Write([]byte(err.Error()))
			} else {
				retbyte, err := json.Marshal(ep)
				if err != nil {
					res.WriteHeader(http.StatusInternalServerError)
					res.Write([]byte(err.Error()))
				} else {
					res.Header().Set("Content-Type", "application/json")
					res.WriteHeader(http.StatusOK)
					res.Write(retbyte)
				}
			}
		}
	}
}

func (h *HttpTargetHandler) UpdatePath(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if len(id) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("missing id in url's query"))
	} else {
		iid, err := strconv.Atoi(id)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("id is not integer"))
		} else {
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
			} else {
				ep := &model.EndPoint{}
				err = json.Unmarshal(bodyBytes, ep)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					res.Write([]byte(err.Error()))
				} else {
					err = endPoints.Update(iid, ep)
					if err != nil {
						res.WriteHeader(http.StatusBadRequest)
						res.Write([]byte(err.Error()))
					} else {
						retbyte, err := json.Marshal(ep)
						if err != nil {
							res.WriteHeader(http.StatusInternalServerError)
							res.Write([]byte(err.Error()))
						} else {
							res.Header().Set("Content-Type", "application/json")
							res.WriteHeader(http.StatusOK)
							res.Write(retbyte)
						}
					}
				}
			}
		}
	}
}

func (h *HttpTargetHandler) GetPaths(res http.ResponseWriter, req *http.Request) {
	list := endPoints.List()
	jsonBytes, err := json.Marshal(list)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
	} else {
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(jsonBytes)
	}
}

func (h *HttpTargetHandler) DeletePath(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if len(id) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("missing id in url's query"))
	} else {
		iid, err := strconv.Atoi(id)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("id is not integer"))
		} else {
			err = endPoints.Delete(iid)
			if err != nil {
				if err.Error() == "notfound" {
					res.WriteHeader(http.StatusNotFound)
					res.Write([]byte("not found"))
				} else {
					res.WriteHeader(http.StatusInternalServerError)
					res.Write([]byte(err.Error()))
				}
			} else {
				res.WriteHeader(http.StatusNoContent)
			}
		}
	}
}

func (h *HttpTargetHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/api/paths" {
		switch req.Method {
		case http.MethodGet:
			h.GetPaths(res, req)
		case http.MethodPost:
			h.CreatePath(res, req)
		case http.MethodPut:
			h.UpdatePath(res, req)
		case http.MethodDelete:
			h.DeletePath(res, req)
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	} else if strings.HasPrefix(req.URL.Path, "/docs") {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if req.URL.Path == "/docs" || req.URL.Path == "/docs/" {
			res.Header().Set("Location", "/docs/index.html")
			res.WriteHeader(http.StatusMovedPermanently)
			return
		} else if strings.HasSuffix(req.URL.Path, "/") {
			res.Header().Set("Location", req.URL.Path+"index.html")
			res.WriteHeader(http.StatusMovedPermanently)
			return
		}
		filePath := strings.ReplaceAll(req.URL.Path, "/docs/", "api/")
		dirFilePath := "[DIR]" + filePath
		paths := static.GetPathTree("api")
		for _, path := range paths {
			if path == dirFilePath {
				res.Header().Set("Location", req.URL.Path+"/index.html")
				res.WriteHeader(http.StatusMovedPermanently)
				return
			}
			if path == filePath {
				fdata, err := static.GetFile(filePath)
				if err != nil {
					res.WriteHeader(http.StatusInternalServerError)
					res.Write([]byte(err.Error()))
					return
				}
				res.Header().Set("Content-Type", fdata.ContentType)
				res.WriteHeader(http.StatusOK)
				res.Write(fdata.Bytes)
				return
			}
		}
		res.WriteHeader(http.StatusNotFound)
	} else if strings.HasPrefix(req.URL.Path, "/delay/") {
		var minDelay, maxDelay int
		params, err := ParsePathParams("/delay/{minDelay}", req.URL.Path)
		if err != nil {
			params, err = ParsePathParams("/delay/{minDelay}/{maxDelay}", req.URL.Path)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				res.Write([]byte("Provided path not exist. try /delay/{minDelay} or /delay/{minDelay}/{maxDelay}"))
				return
			} else {
				strMinDelay := params["minDelay"]
				strMaxDelay := params["maxDelay"]
				mid, err := strconv.Atoi(strMinDelay)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					res.Write([]byte("Provided min delay  is not integer."))
					return
				}
				minDelay = mid
				mad, err := strconv.Atoi(strMaxDelay)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					res.Write([]byte("Provided max delay  is not integer."))
					return
				}
				maxDelay = mad
			}
		} else {
			strMinDelay := params["minDelay"]
			minDelay, err := strconv.Atoi(strMinDelay)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				res.Write([]byte("Provided min delay  is not integer."))
				return
			}
			maxDelay = minDelay
		}
		randDelay := minDelay + rand.Intn(maxDelay-minDelay)
		logrus.Debugf("Path %s delay %d, min %d, max %d", req.URL.Path, randDelay, minDelay, maxDelay)
		time.Sleep(time.Duration(randDelay) * time.Millisecond)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(fmt.Sprintf("OK, delayed between %d and %d : %d", minDelay, maxDelay, randDelay)))
	} else if strings.HasPrefix(req.URL.Path, "/code/") {
		scode := "200"
		smin := "0"
		smad := "0"
		var code, minDelay, maxDelay int
		params, err := ParsePathParams("/code/{resCode}", req.URL.Path)
		if err != nil {
			params, err = ParsePathParams("/code/{resCode}/{minDelay}", req.URL.Path)
			if err != nil {
				params, err = ParsePathParams("/code/{resCode}/{minDelay}/{maxDelay}", req.URL.Path)
				if err != nil {
					res.WriteHeader(http.StatusNotFound)
					res.Write([]byte("Provided path not exist. try /code/{resCode} or /code/{resCode}/{minDelay} or /code/{resCode}/{minDelay}/{maxDelay}."))
					return
				}
				scode = params["resCode"]
				smin = params["minDelay"]
				smad = params["maxDelay"]
			} else {
				scode = params["resCode"]
				smin = params["minDelay"]
				smad = smin
			}
		} else {
			scode = params["resCode"]
		}
		codeD, codeErr := strconv.Atoi(scode)
		minD, minErr := strconv.Atoi(smin)
		maxD, maxErr := strconv.Atoi(smad)
		if codeErr != nil || minErr != nil || maxErr != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("The code, minimum delay or maximum delay is not an integer."))
			return
		}
		code = codeD
		minDelay = minD
		maxDelay = maxD

		randDelay := minDelay + rand.Intn(maxDelay-minDelay)
		logrus.Debugf("Path %s delay %d, min %d, max %d", req.URL.Path, randDelay, minDelay, maxDelay)
		time.Sleep(time.Duration(randDelay) * time.Millisecond)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(code)
		res.Write([]byte(fmt.Sprintf("OK, delayed between %d and %d : %d", minDelay, maxDelay, randDelay)))
	} else {
		ep := endPoints.GetByPath(req.URL.Path)
		if ep == nil {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte("Not Found"))
			return
		}

		randDelay := ep.DelayMinMs + rand.Intn(ep.DelayMaxMs-ep.DelayMinMs)
		logrus.Debugf("Path %s delay %d, min %d, max %d", req.URL.Path, randDelay, ep.DelayMinMs, ep.DelayMaxMs)
		time.Sleep(time.Duration(randDelay) * time.Millisecond)

		if ep.ReturnHeaders != nil {
			for k, v := range ep.ReturnHeaders {
				res.Header()[k] = v
			}
		}
		res.WriteHeader(ep.ReturnCode)
		if len(ep.ReturnBody) > 0 {
			res.Write([]byte(ep.ReturnBody))
		}
	}
}

// ParsePathParams parse request path param according to path template and extract its values.
func ParsePathParams(template, path string) (map[string]string, error) {
	var pth string
	if strings.Contains(path, "?") {
		pth = path[:strings.Index(path, "?")]
	} else {
		pth = path
	}
	templatePaths := strings.Split(template, "/")
	pathPaths := strings.Split(pth, "/")
	if len(templatePaths) != len(pathPaths) {
		return nil, fmt.Errorf("pathElement length not equals to templateElement length")
	}
	ret := make(map[string]string)
	for idx, templateElement := range templatePaths {
		pathElement := pathPaths[idx]
		if len(templateElement) > 0 && len(pathElement) > 0 {
			if templateElement[:1] == "{" && templateElement[len(templateElement)-1:] == "}" {
				tKey := templateElement[1 : len(templateElement)-1]
				ret[tKey] = pathElement
			} else if templateElement != pathElement {
				return nil, fmt.Errorf("template %s not compatible with path %s", template, path)
			}
		}
	}
	return ret, nil
}
