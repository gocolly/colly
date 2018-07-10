// Copyright 2018 Adam Tauber
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package debug

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// WebDebugger is a web based debuging frontend for colly
type WebDebugger struct {
	// Address is the address of the web server. It is 127.0.0.1:7676 by default.
	Address         string
	initialized     bool
	CurrentRequests map[uint32]requestInfo
	RequestLog      []requestInfo
}

type requestInfo struct {
	URL            string
	Started        time.Time
	Duration       time.Duration
	ResponseStatus string
	ID             uint32
	CollectorID    uint32
}

// Init initializes the WebDebugger
func (w *WebDebugger) Init() error {
	if w.initialized {
		return nil
	}
	defer func() {
		w.initialized = true
	}()
	if w.Address == "" {
		w.Address = "127.0.0.1:7676"
	}
	w.RequestLog = make([]requestInfo, 0)
	w.CurrentRequests = make(map[uint32]requestInfo)
	http.HandleFunc("/", w.indexHandler)
	http.HandleFunc("/status", w.statusHandler)
	log.Println("Starting debug webserver on", w.Address)
	go http.ListenAndServe(w.Address, nil)
	return nil
}

// Event updates the debugger's status
func (w *WebDebugger) Event(e *Event) {
	switch e.Type {
	case "request":
		w.CurrentRequests[e.RequestID] = requestInfo{
			URL:         e.Values["url"],
			Started:     time.Now(),
			ID:          e.RequestID,
			CollectorID: e.CollectorID,
		}
	case "response", "error":
		r := w.CurrentRequests[e.RequestID]
		r.Duration = time.Since(r.Started)
		r.ResponseStatus = e.Values["status"]
		w.RequestLog = append(w.RequestLog, r)
		delete(w.CurrentRequests, e.RequestID)
	}
}

func (w *WebDebugger) indexHandler(wr http.ResponseWriter, r *http.Request) {
	wr.Write([]byte(`<!DOCTYPE html>
<html>
<head>
 <title>Colly Debugger WebUI</title>
 <script src="https://code.jquery.com/jquery-latest.min.js" type="text/javascript"></script>
 <link rel="stylesheet" type="text/css" href="https://semantic-ui.com/dist/semantic.min.css">
</head>
<body>
<div class="ui inverted vertical masthead center aligned segment" id="menu">
 <div class="ui tiny secondary inverted menu">
   <a class="item" href="/"><b>Colly WebDebugger</b></a>
 </div>
</div>
<div class="ui grid container">
 <div class="row">
  <div class="eight wide column">
   <h1>Current Requests <span id="current_request_count"></span></h1>
   <div id="current_requests" class="ui small feed"></div>
  </div>
  <div class="eight wide column">
   <h1>Finished Requests <span id="request_log_count"></span></h1>
   <div id="request_log" class="ui small feed"></div>
  </div>
 </div>
</div>
<script>
function curRequestTpl(url, started, collectorId) {
  return '<div class="event"><div class="content"><div class="summary">' + url + '</div><div class="meta">Collector #' + collectorId + ' - ' + started + "</div></div></div>";
}
function requestLogTpl(url, duration, collectorId) {
  return '<div class="event"><div class="content"><div class="summary">' + url + '</div><div class="meta">Collector #' + collectorId + ' - ' + (duration/1000000000) + "s</div></div></div>";
}
function fetchStatus() {
  $.getJSON("/status", function(data) {
    $("#current_requests").html("");
    $("#request_log").html("");
    $("#current_request_count").text('(' + Object.keys(data.CurrentRequests).length + ')');
    $("#request_log_count").text('(' + data.RequestLog.length + ')');
    for(var i in data.CurrentRequests) {
      var r = data.CurrentRequests[i];
      $("#current_requests").append(curRequestTpl(r.Url, r.Started, r.CollectorId));
    }
    for(var i in data.RequestLog.reverse()) {
      var r = data.RequestLog[i];
      $("#request_log").append(requestLogTpl(r.Url, r.Duration, r.CollectorId));
    }
    setTimeout(fetchStatus, 1000);
  });
}
$(document).ready(function() {
    fetchStatus();
});
</script>
</body>
</html>
`))
}

func (w *WebDebugger) statusHandler(wr http.ResponseWriter, r *http.Request) {
	jsonData, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		panic(err)
	}
	wr.Write(jsonData)
}
