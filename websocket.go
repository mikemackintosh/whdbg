package main

import (
	"html/template"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	home     = template.Must(template.New("").Parse(homeHTML))
	sock     = template.Must(template.New("").Parse(sockHTML))
)

const homeHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<title>Whdbg</title>
<style type="text/css">
</style>
</head>
<body>
<ul>
{{range $sub, $client := .subs}}
<li><a href="/_/{{$sub}}">{{$sub}} -- {{ $client.lastAccess }}</a></li>
{{end}}
</ul>
</body>
</html>`

const sockHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<title>Chat Example</title>
<script type="text/javascript">
window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    console.log("Socket", "{{.ws}}");
    if (window["WebSocket"]) {
        conn = new WebSocket("{{.ws}}");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};
</script>
<style type="text/css">
html {
    overflow: hidden;
}

body {
    overflow: hidden;
    padding: 0;
    margin: 0;
    width: 100%;
    height: 100%;
    background: gray;
}

#log {
    background: white;
    margin: 0;
    padding: 0.5em 0.5em 0.5em 0.5em;
    position: absolute;
    top: 30px;
    left: 0.5em;
    right: 0.5em;
    bottom: 3em;
    overflow: auto;
}

</style>
</head>
<body>
<div>To see your requests, send them to: <code>https://{{.sub}}.whdbg.dev</code>. </div>
<div id="log"></div>
</body>
</html>`
