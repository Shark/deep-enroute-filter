(function() {
  function wsUri() {
    var loc = window.location, new_uri;
    if (loc.protocol === "https:") {
        new_uri = "wss:";
    } else {
        new_uri = "ws:";
    }
    new_uri += "//" + loc.host;
    new_uri += loc.pathname + "ws";
    return new_uri;
  }

  function updateTimestamps() {
    for(var timestamp of document.querySelectorAll('.timestamp')) {
      let secondsAgo = Math.floor((Date.now() - Date.parse(timestamp.dataset.value))/1000);
      if(secondsAgo >= 0) {
        if(timestamp.dataset.hideIfPast) {
          timestamp.innerHTML = '–';
        } else {
          timestamp.innerHTML = `${secondsAgo} seconds ago`;
        }
      } else {
        timestamp.innerHTML = `in ${secondsAgo*(-1)} seconds`;
      }
    }
  }
  setInterval(updateTimestamps, 5000);

  var reconnectTimer;

  function connectWebSocket() {
    ws = new WebSocket(wsUri());
    ws.onopen = function(evt) {
      console.log("WebSocket open");

      if(reconnectTimer) {
        clearInterval(reconnectTimer);
        reconnectTimer = 0;
      }

      document.querySelector('#connection-status').innerHTML = 'connected';
    }
    ws.onclose = function(evt) {
      console.log("WebSocket closed");

      if(!reconnectTimer) {
        reconnectTimer = setInterval(connectWebSocket, 2000);
      }

      document.querySelector('#connection-status').innerHTML = 'disconnected';
    }
    ws.onmessage = function(evt) {
      console.log("WebSocket response: " + evt.data);
      let data = JSON.parse(evt.data);

      if(data.Type === 'ProcessedMessages') {
        let container = document.querySelector('#processed-messages--container');
        container.innerHTML = data.Payload.map(message => {
          let ruleProcessingResultMarkup = message.RuleProcessingResults.map(result => {
            var resultStr = result.RuleName;
            if(result.RuleMessage !== null) {
              resultStr += `: ${result.RuleMessage}`;
            }
            return `
              <span class="tag ${result.Allowed ? 'is-success' : 'is-warning'}">
                ${resultStr}
              </span>
            `;
          }).join('');

          return `
            <tr>
              <td class="timestamp" data-value="${message.Timestamp}"></td>
              <td><pre>${message.Destination}</pre></td>
              <td><pre>${message.Method}</pre></td>
              <td><pre>${message.Path}</pre></td>
              <td>${ruleProcessingResultMarkup}</td>
            </tr>
          `;
        }).join('\n');

        updateTimestamps();
      } else if(data.Type === 'CoreEndpointsEvent') {
        let container = document.querySelector('#core-endpoints--container')
        var rowMarkup = []
        for(var key of Object.keys(data.Payload)) {
          rowMarkup.push(`
            <tr>
              <td><pre>${key}</pre></td>
              <td>${data.Payload[key]}</td>
            </tr>
          `)
        }
        container.innerHTML = rowMarkup.join('\n');
      } else if(data.Type === 'RateLimitStateEvent') {
        let container = document.querySelector('#rate-limit-state--container')
        var rowMarkup = []
        for(var bucket of data.Payload) {
          rowMarkup.push(`
            <tr>
              <td><pre>${bucket.SrcIP}</pre></td>
              <td><pre>${bucket.DstIP}</pre></td>
              <td>${bucket.BucketRemaining}/${bucket.BucketCapacity}</td>
              <td class="timestamp" data-value="${bucket.BucketReset}" data-hide-if-past="true"></td>
            </tr>
          `)
        }
        container.innerHTML = rowMarkup.join('\n');
        updateTimestamps();
      }
       // {"Type":"ProcessedMessages","Payload":[{"Timestamp":"2018-07-07T15:30:07.11355415+02:00","Method":"GET","Path":"/basic","RuleProcessingResults":[{"Allowed":true,"RuleName":"MethodRule","RuleMessage":null}]}]}

    }
    ws.onerror = function(evt) {
      console.log(`WebSocket error: ${evt.data}"`);
    }
  }

  connectWebSocket();
})()
