// Legacy function for parsing time
// function parseSentryTime(t) {
//     let parts = t.split(' ')
//     let utc_t = `${parts[0]}T${parts[1]}:00.000Z`;
//     let newDate = new Date(utc_t);
//     return newDate.getTime() / 1000;
// }

// Create base chart
var chart = LightweightCharts.createChart(document.getElementById('chart-v1'), {
    width: 1280,
    height: 720,
    layout: {
        backgroundColor: '#FFF',
        textColor: '#000',
        fontSize: 16,
    },
    timeScale: {
        timeVisible: true,
        borderColor: '#000',
    },
    rightPriceScale: {
        borderColor: '#000',
    },
    grid: {
        horzLines: {
            color: '#F0F3FA',
        },
        vertLines: {
            color: '#F0F3FA',
        },
    },
    crosshair: {
        mode: LightweightCharts.CrosshairMode.Normal,
    },
    watermark: {
        color: 'rgba(11, 94, 29, 0.4)',
        visible: true,
        text: 'Kerrigan Chart v0.1',
        fontSize: 32,
        horzAlign: 'left',
        vertAlign: 'bottom',
    },
});

// Candlesticks data
let chartData = [];
let candleSeries = chart.addCandlestickSeries();
fetch('https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=5m&limit=1000')
    .then(function (response) {
        for (var pair of response.headers.entries()) {
            console.log(pair[0] + ': ' + pair[1]);
        }
        return response.json();
    })
    .then(function (data) {
        for (let e of data) {
            chartData.push({
                close: parseFloat(e[4]),
                high: parseFloat(e[2]),
                low: parseFloat(e[3]),
                open: parseFloat(e[1]),
                time: e[0] / 1000
            })
        };
        candleSeries.setData(chartData);
    })
    .catch(function (error) {
        console.log('Request failed', error)
    });

let sentrySeries = chart.addLineSeries();
let historySentryData = []
fetch('/history')
    .then(function (response) {
        return response.json();
    })
    .then(function (data) {
        for (let e of data) {
            historySentryData.push({
                time: e.time,
                value: e.value
            })
        };
        sentrySeries.setData(historySentryData);
    })
    .catch(function (error) {
        console.log('Request failed', error)
    });

// Websocket portion
const baseURL = 'wss://stream.binance.com:9443';
// const tradeStream = baseURL + '/ws/btcusdt@trade';
// const tradeSocket = new WebSocket(tradeStream);
const k5mStreamURI = baseURL + '/ws/btcusdt@kline_5m';
const k5mSocket = new WebSocket(k5mStreamURI);

// tradeSocket.onmessage = (message) => {
//     document.getElementById("trade").textContent = message.data;
// }

k5mSocket.onmessage = (message) => {
    let data = JSON.parse(message.data);
    let rtCandlestickPrice = {
        time: data.k.t / 1000,
        close: data.k.c,
        high: data.k.h,
        low: data.k.l,
        open: data.k.o
    }
    candleSeries.update(rtCandlestickPrice);
}

// real time update for sentry data
const sentryURI = 'wss://mooner.dace.dev/ws'
const sentrySocket = new WebSocket(sentryURI)

sentrySocket.onopen = (event) => {
    console.log("successfully connected to Sentry socket!")
}

sentrySocket.onclose = (event) => {
    console.log("the websocket connection closed for some reason.")
}
sentrySocket.onmessage = (message) => {
    let data = JSON.parse(message.data);
    console.log(data)
    if (data.m == "sentry") {
        let rtSentryData = {
            time: data.d.t,
            value: data.d.v
        }
        sentrySeries.update(rtSentryData)
    } else {
        console.log("PingPong message received!")
    }

}

sentrySocket.onerror = (event) => {
    console.log("There was some error with the websocket connection.")
}
// WORK IN PROGRESS
// var markers = [];
// markers.push({ time: data[88].time, position: 'aboveBar', color: '#e91e63', shape: 'arrowDown', text: 'Sell @ ' + data[88].high.toFixed(2) });
// markers.push({ time: data[13].time, position: 'belowBar', color: '#2196F3', shape: 'arrowUp', text: 'Buy @ ' + data[13].low.toFixed(2) });
// candleSeries.setMarkers(markers);