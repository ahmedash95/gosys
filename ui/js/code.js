maxChartLables = 10
var config = {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            label: "Requets",
            backgroundColor: "#375eab",
            borderColor: "#375eab",
            data: [],
            fill: false,
        }]
    },
    options: {
        responsive: true,

        scales: {
            xAxes: [{
                display: true,
                scaleLabel: {
                    display: true,
                    labelString: 'Time'
                }
            }],
            yAxes: [{
                display: true,
                scaleLabel: {
                    display: true,
                    labelString: 'Requests'
                },
                ticks: {
                    min: 0,
                    max: 0,
                    stepSize: 10
                }
            }]
        }
    }
};

window.onload = function() {
    var ctx = document.getElementById("canvas").getContext("2d");
    window.myLine = new Chart(ctx, config);
};

var pushLogUpdate = function(time,count){
    var index = window.myLine.config.data.labels.indexOf(time)
    if(index == -1){
        window.myLine.config.data.labels.push(time)
        index = window.myLine.config.data.labels.length - 1
    }
    if(window.myLine.config.data.datasets[0].data[index] === undefined){
        window.myLine.config.data.datasets[0].data[index] = 0
    }
    window.myLine.config.data.datasets[0].data[index] += count
    // take last 10 elements
    if(window.myLine.config.data.datasets[0].data.length > maxChartLables){
        window.myLine.config.data.datasets[0].data = window.myLine.config.data.datasets[0].data.slice(Math.max(window.myLine.config.data.datasets[0].data.length - maxChartLables, 1))
        window.myLine.config.data.labels = window.myLine.config.data.labels.slice(Math.max(window.myLine.config.data.labels.length - maxChartLables, 1))
    }

    maxValue = Math.max.apply(Math,window.myLine.data.datasets[0].data)
    window.myLine.config.options.scales.yAxes[0].ticks.max = maxValue + ((25 * maxValue) / 100)

    window.myLine.update()
}

setTimeout(function(){
    // load old logs
    $.ajax({
        url : "http://127.0.0.1:3000/logs",
        success : function(response){
            var messages = JSON.parse(response)
            for(var i = 0; i < messages.length; i++){
                pushLogUpdate(messages[i].time,messages[i].hits)
            }
            startWS()
        },
        error : function(){
            startWS()
        }
    })
},300)

var startWS = function(){
    // WebScoket Messages
    var sock = null;
    var wsuri = "ws://127.0.0.1:3000/ws";

    sock = new WebSocket(wsuri);

    sock.onopen = function() {
        console.log("connected to " + wsuri);
    }

    sock.onclose = function(e) {
        console.log("connection closed (" + e.code + ")");
    }

    sock.onmessage = function(e) {
        console.log("message received: " + e.data);
        var msg = JSON.parse(e.data);
        setTimeout(function(){
            pushLogUpdate(msg.time,msg.hits)
        },300)
    }
}