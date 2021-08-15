import React, { useState, useRef, useMemo, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import { Tooltip, Badge, FormGroup, Label, Input } from 'reactstrap';
import Select from 'react-select';

const statusCodes = [
  { "value": 100, "label": "100 Continue" },
  { "value": 101, "label": "101 Switching Protocols" },
  { "value": 102, "label": "102 Processing" },
  { "value": 103, "label": "103 Early Hints" },
  { "value": 200, "label": "200 OK" },
  { "value": 201, "label": "201 Created" },
  { "value": 202, "label": "202 Accepted" },
  { "value": 203, "label": "203 Non-Authoritative Information" },
  { "value": 204, "label": "204 No Content" },
  { "value": 205, "label": "205 Reset Content" },
  { "value": 206, "label": "206 Partial Content" },
  { "value": 207, "label": "207 Multi-Status" },
  { "value": 208, "label": "208 Already Reported" },
  { "value": 226, "label": "226 IM Used" },
  { "value": 300, "label": "300 Multiple Choices" },
  { "value": 301, "label": "301 Moved Permanently" },
  { "value": 302, "label": "302 Found" },
  { "value": 303, "label": "303 See Other" },
  { "value": 304, "label": "304 Not Modified" },
  { "value": 305, "label": "305 Use Proxy" },
  { "value": 307, "label": "307 Temporary Redirect" },
  { "value": 308, "label": "308 Permanent Redirect" },
  { "value": 400, "label": "400 Bad Request" },
  { "value": 401, "label": "401 Unauthorized" },
  { "value": 402, "label": "402 Payment Required" },
  { "value": 403, "label": "403 Forbidden" },
  { "value": 404, "label": "404 Not Found" },
  { "value": 405, "label": "405 Method Not Allowed" },
  { "value": 406, "label": "406 Not Acceptable" },
  { "value": 407, "label": "407 Proxy Authentication Required" },
  { "value": 408, "label": "408 Request Timeout" },
  { "value": 409, "label": "409 Conflict" },
  { "value": 410, "label": "410 Gone" },
  { "value": 411, "label": "411 Length Required" },
  { "value": 412, "label": "412 Precondition Failed" },
  { "value": 413, "label": "413 Request Entity Too Large" },
  { "value": 414, "label": "414 Request URI Too Long" },
  { "value": 415, "label": "415 Unsupported Media Type" },
  { "value": 416, "label": "416 Requested Range Not Satisfiable" },
  { "value": 417, "label": "417 Expectation Failed" },
  { "value": 418, "label": "418 I'm a teapot" },
  { "value": 421, "label": "421 Misdirected Request" },
  { "value": 422, "label": "422 Unprocessable Entity" },
  { "value": 423, "label": "423 Locked" },
  { "value": 424, "label": "424 Failed Dependency" },
  { "value": 425, "label": "425 Too Early" },
  { "value": 426, "label": "426 Upgrade Required" },
  { "value": 428, "label": "428 Precondition Required" },
  { "value": 429, "label": "429 Too Many Requests" },
  { "value": 431, "label": "431 Request Header Fields Too Large" },
  { "value": 451, "label": "451 Unavailable For Legal Reasons" },
  { "value": 500, "label": "500 Internal Server Error" },
  { "value": 501, "label": "501 Not Implemented" },
  { "value": 502, "label": "502 Bad Gateway" },
  { "value": 503, "label": "503 Service Unavailable" },
  { "value": 504, "label": "504 Gateway Timeout" },
  { "value": 505, "label": "505 HTTP Version Not Supported" },
  { "value": 506, "label": "506 Variant Also Negotiates" },
  { "value": 507, "label": "507 Insufficient Storage" },
  { "value": 508, "label": "508 Loop Detected" },
  { "value": 510, "label": "510 Not Extended" },
  { "value": 511, "label": "511 Network Authentication Required" },
]

const Message = (props) => {
  let {message} = props;

  const [msgRead, setMsgRead] = useState(false);
  const [showBody, setShowBody] = useState(false);
  const onClick = () => { setShowBody(!showBody); setMsgRead(true); };

  let methodMapping = {
    "GET":    "success",
    "POST":   "warning",
    "DELETE": "danger",
    "PUT":    "secondary",
    "PATCH":  "info",
  }

  return (
    <div className="row message-item over-actions-trigger hover-shadow py-2 px-1 mx-0">
      <div className="col col-md-3">
        {message.timestamp}
      </div>
      <div className="col col-md-9 text-center">
        <a className={`ctaSubject `+ (msgRead ? 'ctaSubjectRead' : 'ctaSubjectUnread')} onClick={onClick}>
          <Badge color={methodMapping[message.request.Method]}>{message.request.Method}</Badge> https://{message.request.Host}{message.url} {message.request.Proto}
        </a>
      </div>
      <div className="row message-detail" style={{ display: showBody ? "block" : "none" }}>
        <div className="msg-detail-box">
        {message.dump.split("\n").map(function(item, idx) {
          return (
            <span key={idx}>
              {item}<br/>
            </span>
          )
        })}
        </div>
      </div>
    </div>
  )
}

const Messages = (props) => {
  return(
    <ul>
        {props.messages.slice().reverse().filter(msg => msg !== null).map((message, idx) => {
            if( message != null ){
              return (
                <div className="msg-panel">
                  <div className="row messages-header">
                    <div className="col col-md-3">Time</div>
                    <div className="col col-md-9 text-center">Request</div>
                  </div>
                  <Message key={idx} message={message} />
                </div>
              )
            }
        })}
    </ul>
  )
}

const Listener = (props) => {
  const baseWsURL = process.env.NODE_ENV === 'production' ?  'wss://whdbg.dev' : 'ws://localhost:8080';
  // const baseWsURL = 'ws://localhost:8080'

    //Public API that will echo messages sent to it back to the client
  const [socketUrl, setSocketUrl] = useState(baseWsURL+'/ws/'+props.match.params.listener);
  const [listener, setListener] = useState('');
  const [selectedStatus, setSelectedStatus] = useState("");
  const [responseBody, setResponseBody] = useState({});

  const [copyTooltipOpen, setCopyTooltipOpen] = useState(false);
  const copyToggle = () => setCopyTooltipOpen(!copyTooltipOpen);

  const [responseOpen, setResponseOpen] = useState(false);
  const responseToggle = () => setResponseOpen(!responseOpen);

  const messageHistory = useRef([]);

  const {
    lastJsonMessage,
    readyState,
  } = useWebSocket(socketUrl);

  const connectionStatus = {
    [ReadyState.CONNECTING]: 'Connecting',
    [ReadyState.OPEN]: 'Open',
    [ReadyState.CLOSING]: 'Closing',
    [ReadyState.CLOSED]: 'Closed',
    [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
  }[readyState];

  messageHistory.current = useMemo(() =>
      messageHistory.current.concat(lastJsonMessage), [lastJsonMessage])

  useEffect(() => {
    if ( connectionStatus == "Open" ){
      var payload = {responseBody: "", statusCode: 200}

      if ( 'target' in responseBody ) {
        if ( 'value' in responseBody.target ) {
          payload.responseBody = responseBody.target.value
        }
      }

      if (selectedStatus !== null) {
        payload.statusCode = selectedStatus.value
      }

      const requestOptions = {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
      };
      fetch('/api/'+props.match.params.listener+'/update', requestOptions)
        .then(response => response.json());
    }
  }, [selectedStatus, responseBody]);

  useEffect(() => {
    console.log("Updating listener:", props.match.params.listener);
    setSocketUrl(baseWsURL+'/ws/'+props.match.params.listener);
    setListener(`https://${props.match.params.listener}.whdbg.dev`);
  }, []);

  return (
      <>
      <div className="listening">
        <div className="instructions">
          Send your HTTP requests to <code id="copy" className="codeToCopy" onClick={() =>
                {navigator.clipboard.writeText(listener)}}>
            {listener}
          </code>
          <Tooltip placement="top" trigger="click" isOpen={copyTooltipOpen} target="copy" toggle={copyToggle}>
            Copied
          </Tooltip>

          <FormGroup className="msg-panel">
            <div className="row">
              <h3><a className="toggleResponseBody" href="#" onClick={responseToggle}>Response</a></h3>
            </div>
            <div style={{ display: responseOpen ? "block" : "none" }}>
            <div className="row">
              <div className="col-sm-4">
                <Label for="responseStatusCode">Status Code</Label>
                <Select
                  id="responseStatusCode"
                  className="statusSelect"
                  defaultValue={statusCodes[4]}
                  onChange={setSelectedStatus}
                  noOptionsMessage={"Select a Status Code"}
                  options={statusCodes}/>
              </div>
              <div className="col-sm-8">
                <Label for="responseBody">Response Body</Label>
                <Input className="statusSelect" type="textarea" name="text" id="responseBody" onChange={setResponseBody} placeholder="No response body will reflect the body. If you want to respond with a specific payload, put it here!" />
              </div>
              </div>
              </div>
          </FormGroup>

        </div>
        <div>
          <h5 className="fs-0 px-3 pt-3 pb-2 mb-0 ">REQUESTS</h5>
          <span>The WebSocket is currently {connectionStatus}</span>
          <Messages messages={messageHistory.current} />
          </div>
      </div>
      </>
  );
}

export default Listener;
