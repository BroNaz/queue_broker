package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s *Server) PushMessagesInAQueue() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		s.logger.Info("Accepted the request: ", r.URL, " Method: ", r.Method)
		s.logger.Debug("Query: ", r.URL.Query())

		vars := mux.Vars(r)
		key := vars["key"]

		vals := r.URL.Query()
		val, ok := vals["v"]
		var value string
		if ok {
			if len(val) >= 1 {
				value = val[0] // The first `?v=...`
			}
		}
		if ok && value == "" {
			s.logger.Debug("The value exists, but it is null")
			rw.WriteHeader(415) // Unsupported Media Type
			return
		}

		s.logger.Debug("Store key: ", key, " value: ", value)
		s.store.Store(key, value)

		rw.WriteHeader(201) // Created
	}
}

func (s *Server) PopFromTheQueue() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		s.logger.Info("Accepted the request: ", r.URL, " Method: ", r.Method)

		vars := mux.Vars(r)
		key := vars["key"]

		if res, ok := s.store.Pop(key); ok == true {
			s.logger.Debug("Answer is ready immediately! Answer: ", res)
			rw.Header().Set("Content-Type", "charset=UTF-8")
			rw.Write([]byte(res))
			return
		}

		vals := r.URL.Query()
		val, ok := vals["timeout"]
		if ok == false {
			s.logger.Debug("Did not give a timeout")
			errorJson(rw, "error", 404)
			return
		}

		var value string
		if ok {
			if len(val) >= 1 {
				value = val[0] // The first `?v=...`
			}
		}
		timeout, err := strconv.ParseInt(value, 0, 64)
		if ok && value == "" && err != nil {
			s.logger.Debug("The value exists, but it is null")
			rw.WriteHeader(415) // Unsupported Media Type
			return
		}

		changed := longPoll(timeout, func(c chan<- bool) {
			resLongPoll := false
			for i := 0; i < int(timeout+1); i++ {

				if ok := s.store.Exist(key); ok == true {
					resLongPoll = true
					break
				}

				time.Sleep(time.Second)
			}
			c <- resLongPoll
		})

		if changed == true {
			if res, ok := s.store.Pop(key); ok == true {
				s.logger.Debug("longPoll is true! Answer: ", res)
				rw.Header().Set("Content-Type", "charset=UTF-8")
				rw.Write([]byte(res))
				return
			}

		}
		s.logger.Debug("longPoll is false!")
		errorJson(rw, "error", 404)
	}
}

func longPoll(timeSecond int64, checkUpdates func(chan<- bool)) bool {
	chanRes := make(chan bool)

	go checkUpdates(chanRes)

	select {
	case res := <-chanRes:
		return res
	case <-time.After(time.Duration(timeSecond) * time.Second):
		return false
	}
}

func errorJson(rw http.ResponseWriter, errorText string, errorCode int) {
	type errorJson struct {
		Text   string
		Status int
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(errorCode)

	data := errorJson{
		Text:   errorText,
		Status: errorCode,
	}
	json.NewEncoder(rw).Encode(data)
}
