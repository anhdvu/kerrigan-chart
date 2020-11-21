// Live price on title
const title = document.getElementById("title");
const legend = document.getElementById("legend");
let currentOHLC = {
    open: 0,
    high: 0,
    low: 0,
    close: 0
}

// Create base chart
const chart = LightweightCharts.createChart(document.getElementById('chart-v1'), {
    width: 1280,
    height: 720,
    layout: {
        backgroundColor: '#FFF',
        textColor: '#000',
        fontSize: 16,
    },
    timeScale: {
        timeVisible: true,
        rightOffset: 2
    },
    crosshair: {
        mode: LightweightCharts.CrosshairMode.Normal,
    },
    watermark: {
        color: 'rgba(139, 65, 236, 0.6)',
        visible: true,
        text: 'Kerrigan Chart v0.4',
        fontSize: 32,
        horzAlign: 'left',
        vertAlign: 'bottom',
    },
});

// Candlestick series
const candleSeries = chart.addCandlestickSeries();
let chartData = [];
fetch('https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=5m&limit=1000')
    .then(function (response) {
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
            });
        };
        candleSeries.setData(chartData);
    })
    .catch(function (error) {
        console.log('Request failed', error)
    });

// All sentry-related series
const sentrySeries = chart.addLineSeries();
let historySentryData = [];
const outstandingSeries = chart.addLineSeries({
    color: '#f48fb1',
    lineWidth: 1,
    lineStyle: LightweightCharts.LineStyle.Dashed
});
let outstandingSentryData = [];
const lowballSeries = chart.addLineSeries({
    color: '#f48fb1',
    lineWidth: 1,
    lineStyle: LightweightCharts.LineStyle.Dashed
});
let lowballSentryData = [];
fetch('/history')
    .then(function (response) {
        return response.json();
    })
    .then(function (data) {
        for (let e of data) {
            historySentryData.push({
                time: e.time,
                value: e.value,
            });
            outstandingSentryData.push({
                time: e.time,
                value: e.value + 300
            });
            lowballSentryData.push({
                time: e.time,
                value: e.value - 300
            });
        };
        sentrySeries.setData(historySentryData);
        outstandingSeries.setData(outstandingSentryData);
        lowballSeries.setData(lowballSentryData);
    })
    .catch(function (error) {
        console.log('Request failed', error)
    });

// ################
// Websocket
// ################
const baseURL = 'wss://stream.binance.com:9443';
// const tradeStream = baseURL + '/ws/btcusdt@trade';
// const tradeSocket = new WebSocket(tradeStream);
const k5mStreamURI = baseURL + '/ws/btcusdt@kline_5m';
const k5mSocket = new WebSocket(k5mStreamURI);

k5mSocket.onmessage = (message) => {
    let data = JSON.parse(message.data);
    let rtCandlestickPrice = {
        time: data.k.t / 1000,
        close: data.k.c,
        high: data.k.h,
        low: data.k.l,
        open: data.k.o
    };
    candleSeries.update(rtCandlestickPrice);
    title.textContent = "Kerrigan - " + parseFloat(data.k.c).toFixed(2);
    currentOHLC = {
        open: parseFloat(rtCandlestickPrice.open).toFixed(2),
        high: parseFloat(rtCandlestickPrice.high).toFixed(2),
        low: parseFloat(rtCandlestickPrice.low).toFixed(2),
        close: parseFloat(rtCandlestickPrice.close).toFixed(2)
    }
    legend.textContent = `O ${currentOHLC.open} - H ${currentOHLC.high} - L ${currentOHLC.low} - C ${currentOHLC.close}`;
};


// real time data for sentry series
const sentryURI = 'wss://mooner.dace.dev/ws';
const sentrySocket = new WebSocket(sentryURI);

sentrySocket.onopen = (event) => {
    console.log("successfully connected to Sentry socket!")
}

sentrySocket.onclose = (event) => {
    console.log("the websocket connection closed for some reason.")
}
sentrySocket.onmessage = (message) => {
    let data = JSON.parse(message.data);
    if (data.m == "sentry") {
        let rtSentryData = {
            time: data.d.t,
            value: data.d.v
        };
        let rtOutstandingSentryData = {
            time: data.d.t,
            value: data.d.v + 300
        };
        let rtLowballSentryData = {
            time: data.d.t,
            value: data.d.v - 300
        };
        sentrySeries.update(rtSentryData);
        outstandingSeries.update(rtOutstandingSentryData);
        lowballSeries.update(rtLowballSentryData);
    } else {
        console.log("PingPong message received!")
    }

}

sentrySocket.onerror = (event) => {
    console.log("There was some error with the websocket connection.")
}

// OHLC with crosshair move
chart.subscribeCrosshairMove(function (param) {
    if (param === undefined || param.time === undefined || param.point.x < 0 || param.point.y < 0) {
        legend.textContent = `O ${currentOHLC.open} - H ${currentOHLC.high} - L ${currentOHLC.low} - C ${currentOHLC.close}`
    } else {
        let price = param.seriesPrices.get(candleSeries);
        try {
            let OHLC = {
                open: parseFloat(price.open).toFixed(2),
                high: parseFloat(price.high).toFixed(2),
                low: parseFloat(price.low).toFixed(2),
                close: parseFloat(price.close).toFixed(2)
            }
            legend.textContent = `O ${OHLC.open} - H ${OHLC.high} - L ${OHLC.low} - C ${OHLC.close}`;
        }
        catch (err) {
            console.log("No data at given time yet.")
        }
    }
});

// WORK IN PROGRESS
// let markers = [];
// markers.push({ time: data[88].time, position: 'aboveBar', color: '#e91e63', shape: 'arrowDown', text: 'Sell @ ' + data[88].high.toFixed(2) });
// markers.push({ time: data[13].time, position: 'belowBar', color: '#2196F3', shape: 'arrowUp', text: 'Buy @ ' + data[13].low.toFixed(2) });
// candleSeries.setMarkers(markers);