const title = document.getElementById('title');
const legend = document.getElementById('legend');
const dyde = document.getElementById('dyde');
function convertTime(t) {
    const now = new Date(t);
    return now.toUTCString().substr(17, 8);
}

const tableRecords = [];
const table5m = document.getElementById('timeframe-5m');

// Create base chart
const chart = LightweightCharts.createChart(document.getElementById('kchart'), {
    width: 1600,
    height: 720,
    layout: {
        backgroundColor: 'rgba(19, 23, 34, 1)',
        textColor: 'rgba(255, 255, 255, 0.8)',
        fontSize: 16,
    },
    priceScale: {
        scaleMargins: {
            top: 0.1,
            bottom: 0.1
        },
        entireTextOnly: true
    },
    timeScale: {
        timeVisible: true,
        rightOffset: 16,
        fixLeftEdge: true,
        rightBarStaysOnScroll: true
    },
    grid: {
        horzLines: {
            color: '#F0F3FA',
            style: LightweightCharts.LineStyle.Dotted,
            visible: false
        },
        vertLines: {
            color: '#F0F3FA',
            style: LightweightCharts.LineStyle.Dotted,
            visible: false
        },
    },
    crosshair: {
        mode: LightweightCharts.CrosshairMode.Normal,
    },
    watermark: {
        color: 'rgba(0, 150, 235, 1)',
        visible: true,
        text: 'AESXII Chart v0.5.5',
        fontSize: 24,
        horzAlign: 'left',
        vertAlign: 'bottom',
    },
});

// Price series
const candleSeries = chart.addCandlestickSeries();
let chartData = [];
const volumeSeries = chart.addHistogramSeries({
    priceScaleId: '',
    priceFormat: {
        type: 'volume',
        precision: 3,
    },
    scaleMargins: {
        top: 0.75,
        bottom: 0,
    },
    lastValueVisible: false,
    priceLineVisible: false
});
let volumeData = [];

// Sentry related series
const sentrySeriesConfig = {
    priceLineVisible: false,
    priceLineColor: 'rgba(0, 150, 235, 1)',
    priceLineWidth: 1,
    priceLineStyle: LightweightCharts.LineStyle.SparseDotted,
    priceLineSource: LightweightCharts.PriceLineSource.LastBar,
    lastValueVisible: true,
    color: 'rgba(0, 150, 235, 1)',
    lineWidth: 2,
    lineStyle: LightweightCharts.LineStyle.Solid,
}

const mainSentrySeriesConfig = {
    priceLineVisible: true,
    priceLineColor: 'rgba(125, 23, 166, 1)',
    priceLineWidth: 1,
    priceLineStyle: LightweightCharts.LineStyle.SparseDotted,
    priceLineSource: LightweightCharts.PriceLineSource.LastBar,
    lastValueVisible: true,
    color: 'rgba(125, 23, 166, 1)',
    lineWidth: 2,
    lineStyle: LightweightCharts.LineStyle.Solid,
}

const subSentrySeriesConfig = {
    priceLineVisible: false,
    priceLineColor: 'rgba(0, 150, 235, 0.2)',
    priceLineWidth: 1,
    priceLineStyle: LightweightCharts.LineStyle.SparseDotted,
    priceLineSource: LightweightCharts.PriceLineSource.LastBar,
    lastValueVisible: false,
    color: 'rgba(0, 150, 235, 0.2)',
    lineWidth: 1,
    lineStyle: LightweightCharts.LineStyle.Solid,
}
const sentrySeries = chart.addLineSeries(mainSentrySeriesConfig);
let historySentryData = [];
const sN300Series = chart.addLineSeries(sentrySeriesConfig);
let sN300SentryData = [];
const sN600Series = chart.addLineSeries(sentrySeriesConfig);
let sN600SentryData = [];
const sP300Series = chart.addLineSeries(sentrySeriesConfig);
let sP300SentryData = [];
const sP600Series = chart.addLineSeries(sentrySeriesConfig);
let sP600SentryData = [];
const sN150Series = chart.addLineSeries(subSentrySeriesConfig);
let sN150SentryData = [];
const sP150Series = chart.addLineSeries(subSentrySeriesConfig);
let sP150SentryData = [];
const sN900Series = chart.addLineSeries(sentrySeriesConfig);
let sN900SentryData = [];
const sP900Series = chart.addLineSeries(sentrySeriesConfig);
let sP900SentryData = [];
const sN1200Series = chart.addLineSeries(sentrySeriesConfig);
let sN1200SentryData = [];
const sN1500Series = chart.addLineSeries(sentrySeriesConfig);
let sN1500SentryData = [];

const fetchKline = async () => {
    const timeoffset = 60 * 60 * 24 * 60 * 1000;
    let startTime = Date.now() - timeoffset;
    const limit = 1000;
    const interval = '5m';
    const symbol = 'BTCUSDT';
    let volColor = '';

    while (true) {
        let response = await fetch(`https://api.binance.com/api/v3/klines?symbol=${symbol}&interval=${interval}&startTime=${startTime}&limit=${limit}`);
        if (!response.ok) break;
        let data = await response.json();
        for (let e of data) {
            chartData.push({
                time: e[0] / 1000,
                open: parseFloat(e[1]),
                high: parseFloat(e[2]),
                low: parseFloat(e[3]),
                close: parseFloat(e[4])
            });
            if (parseFloat(e[4]) < parseFloat(e[1])) {
                volColor = 'rgba(255,82,82, 0.2)';
            } else {
                volColor = 'rgba(0, 150, 136, 0.2)';
            }
            volumeData.push({
                time: e[0] / 1000,
                value: parseFloat(e[5]),
                color: volColor
            })
        };
        if (data.length < limit) {
            candleSeries.setData(chartData);
            volumeSeries.setData(volumeData);
            break;
        }
        startTime = data[data.length - 1][6];
    }
}

const fill5mRecordTable = async () => {
    let response = await fetch(`https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=5m&limit=5`);
    if (!response.ok) {
        console.log('Fetch failed!')
    }
    let data = await response.json();
    tempRecordData = data.slice(0, 4).reverse();
    let table5mRecords = []
    for (let e of tempRecordData) {
        table5mRecords.push({
            startTime: convertTime(e[0]),
            open: parseFloat(e[1]).toFixed(2),
            price: parseFloat(e[4]).toFixed(2),
            vol: parseFloat(e[5]).toFixed(3),
            volRatio: parseFloat(e[9] / e[5]).toFixed(4),
            tradeNumber: e[8]
        })
    }
    for (let i = 2; i < table5m.rows.length; i++) {
        table5m.rows[i].cells[0].textContent = table5mRecords[i - 2].startTime;
        table5m.rows[i].cells[1].textContent = table5mRecords[i - 2].price;
        if (table5mRecords[i - 2].price < table5mRecords[i - 2].open) {
            table5m.rows[i].cells[1].style.backgroundColor = 'rgba(255,82,82, 0.8)';
        } else {
            table5m.rows[i].cells[1].style.backgroundColor = 'rgba(0, 150, 136, 0.8)';
        }
        table5m.rows[i].cells[2].textContent = table5mRecords[i - 2].vol;
        if (table5mRecords[i - 2].vol > 399.0) {
            table5m.rows[i].cells[2].style.backgroundColor = 'rgba(255, 217, 0, 0.9)';
        } else if (table5mRecords[i - 2].vol > 199.0) {
            table5m.rows[i].cells[2].style.backgroundColor = 'rgba(255, 217, 0, 0.6)';
        } else if (table5mRecords[i - 2].vol > 99.0) {
            table5m.rows[i].cells[2].style.backgroundColor = 'rgba(255, 217, 0, 0.3)';
        } else {
            table5m.rows[i].cells[2].style.backgroundColor = 'rgba(255, 217, 0, 0.1)';
        }
        table5m.rows[i].cells[3].textContent = table5mRecords[i - 2].volRatio;
        table5m.rows[i].cells[4].textContent = table5mRecords[i - 2].tradeNumber;
    }
}

const fetchSentryHistory = async () => {
    let response = await fetch('/history');
    if (!response.ok) {
        console.log('Fetch failed!')
    }
    let data = await response.json();
    for (let e of data) {
        historySentryData.push({
            time: e.time,
            value: e.value,
        });
        sN300SentryData.push({
            time: e.time,
            value: e.value + 300
        });
        sN600SentryData.push({
            time: e.time,
            value: e.value + 600
        })
        sP300SentryData.push({
            time: e.time,
            value: e.value - 300
        });
        sP600SentryData.push({
            time: e.time,
            value: e.value - 600
        });
        sP150SentryData.push({
            time: e.time,
            value: e.value - 150
        });
        sN150SentryData.push({
            time: e.time,
            value: e.value + 150
        });
        sN900SentryData.push({
            time: e.time,
            value: e.value + 900
        });
        sP900SentryData.push({
            time: e.time,
            value: e.value - 900
        });
        sN1200SentryData.push({
            time: e.time,
            value: e.value + 1200
        });
        sN1500SentryData.push({
            time: e.time,
            value: e.value + 1500
        });
    };
    sentrySeries.setData(historySentryData);
    sN300Series.setData(sN300SentryData);
    sN600Series.setData(sN600SentryData);
    sP300Series.setData(sP300SentryData);
    sP600Series.setData(sP600SentryData);
    sN150Series.setData(sN150SentryData);
    sP150Series.setData(sP150SentryData);
    sN900Series.setData(sN900SentryData);
    sP900Series.setData(sP900SentryData);
    sN1200Series.setData(sN1200SentryData);
    sN1500Series.setData(sN1500SentryData);
}

fetchKline();
fetchSentryHistory();
fill5mRecordTable();
setInterval(function () { fill5mRecordTable() }, 30 * 1000);

// let markers = [];
// markers.push({ time: 1606897800, position: 'aboveBar', color: '#e91e63', shape: 'arrowDown', size: 2, text: 'S @ 19303.39' });
// markers.push({ time: 1606843200, position: 'belowBar', color: '#2196F3', shape: 'arrowUp', size: 2, text: 'B @ 18913.98' });
// candleSeries.setMarkers(markers);

// ################
// Websocket
// ################
const tf1mtds = document.getElementsByClassName('stats 1m');
const tf5mtds = document.getElementsByClassName('stats 5m');
const klineConnect = () => {
    const baseURL = 'wss://stream.binance.com:9443';
    const kline5mStreamName = 'btcusdt@kline_5m';
    const kline1mStreamName = 'btcusdt@kline_1m';
    const multiKlineStreamURI = `${baseURL}/stream?streams=${kline5mStreamName}/${kline1mStreamName}`;
    const klineStreamURI = baseURL + '/ws/' + kline5mStreamName;
    const klineSocket = new WebSocket(klineStreamURI);
    let volColor = '';

    klineSocket.onopen = (event) => {
        console.log('successfully connected to Kline socket!');
    }

    klineSocket.onclose = (event) => {
        console.log('The Kline websocket connection closed for some reason.\nReconnecting in 5 seconds...');
        setTimeout(function () {
            fetchKline();
            klineConnect();
        }, 5000);
    }

    klineSocket.onmessage = (message) => {
        let data = JSON.parse(message.data);
        let rtCandlestickPrice = {
            time: data.k.t / 1000,
            close: data.k.c,
            high: data.k.h,
            low: data.k.l,
            open: data.k.o
        };
        if (data.k.c < data.k.o) {
            volColor = 'rgba(255,82,82, 0.2)';
        } else {
            volColor = 'rgba(0, 150, 136, 0.2)';
        }
        let rtVolumeData = {
            time: data.k.t / 1000,
            value: data.k.v,
            color: volColor
        }
        candleSeries.update(rtCandlestickPrice);
        volumeSeries.update(rtVolumeData);
        title.textContent = 'AESXII Chart - ' + parseFloat(data.k.c).toFixed(2);
        tf5mtds[1].innerHTML = '<p>' + convertTime(data.k.t) + '</p>';
        tf5mtds[2].innerHTML = '<p>' + parseFloat(data.k.c).toFixed(2) + '</p>';
        if (data.k.c < data.k.o) {
            tf5mtds[2].style.backgroundColor = 'rgba(255,82,82, 0.8)';
        } else {
            tf5mtds[2].style.backgroundColor = 'rgba(0, 150, 136, 0.8)';
        }

        tf5mtds[3].innerHTML = '<p>' + parseFloat(data.k.v).toFixed(3) + " BTC" + '</p>';
        if (data.k.v > 399.0) {
            tf5mtds[3].style.backgroundColor = 'rgba(255, 217, 0, 0.9)';
        } else if (data.k.v > 199.0) {
            tf5mtds[3].style.backgroundColor = 'rgba(255, 217, 0, 0.6)';
        } else if (data.k.v > 99.0) {
            tf5mtds[3].style.backgroundColor = 'rgba(255, 217, 0, 0.3)';
        } else {
            tf5mtds[3].style.backgroundColor = 'rgba(255, 217, 0, 0.1)';
        }
        tf5mtds[4].innerHTML = '<p>' + parseFloat(data.k.V / data.k.v).toFixed(4) + '</p>';
        tf5mtds[5].innerHTML = '<p>' + data.k.n + '</p>';
    };
}

// real time data for sentry series
const sentryConnect = () => {
    const sentryURI = 'wss://mooner.dace.dev/ws';
    // const sentryURI = 'ws://localhost:8080/ws';
    const sentrySocket = new WebSocket(sentryURI);

    sentrySocket.onopen = (event) => {
        console.log('successfully connected to Sentry socket!');
    }

    sentrySocket.onclose = (event) => {
        console.log('The Sentry websocket connection closed for some reason.\nReconnecting in 5 seconds...');
        setTimeout(function () {
            fetchSentryHistory();
            sentryConnect();
        }, 5000);
    }
    sentrySocket.onmessage = (message) => {
        let data = JSON.parse(message.data);
        if (data.m == 'sentry') {
            let rtSentryData = {
                time: data.d.t,
                value: data.d.v
            };
            let rtN300SentryData = {
                time: data.d.t,
                value: data.d.v + 300
            };
            let rtN600SentryData = {
                time: data.d.t,
                value: data.d.v + 600
            };
            let rtP300SentryData = {
                time: data.d.t,
                value: data.d.v - 300
            };
            let rtP600SentryData = {
                time: data.d.t,
                value: data.d.v - 600
            };
            let rtP150SentryData = {
                time: data.d.t,
                value: data.d.v - 150
            };
            let rtN150SentryData = {
                time: data.d.t,
                value: data.d.v + 150
            };
            let rtN900SentryData = {
                time: data.d.t,
                value: data.d.v + 900
            };
            let rtP900SentryData = {
                time: data.d.t,
                value: data.d.v - 900
            };
            let rtN1200SentryData = {
                time: data.d.t,
                value: data.d.v + 1200
            };
            let rtN1500SentryData = {
                time: data.d.t,
                value: data.d.v + 1500
            };
            sentrySeries.update(rtSentryData);
            sN300Series.update(rtN300SentryData);
            sN600Series.update(rtN600SentryData);
            sP300Series.update(rtP300SentryData);
            sP600Series.update(rtP600SentryData);
            sN150Series.update(rtN150SentryData);
            sP150Series.update(rtP150SentryData);
            sN900Series.update(rtN900SentryData);
            sP900Series.update(rtP900SentryData);
            sN1200Series.update(rtN1200SentryData);
            sN1500Series.update(rtN1500SentryData);
        } else if (data.m == 'dyde') {
            dyde.textContent = data.d.v.toFixed(2);
        } else {
            console.log('ping');
        }
    }
    sentrySocket.onerror = (event) => {
        console.log('There was some error with the websocket connection.');
        console.log(event.data);
    }
}

klineConnect();
sentryConnect();

// OHLC with crosshair move
chart.subscribeCrosshairMove(function (param) {
    if (param === undefined || param.time === undefined || param.point.x < 0 || param.point.y < 0) {
        // legend.textContent = `O ${currentOHLC.open} - H ${currentOHLC.high} - L ${currentOHLC.low} - C ${currentOHLC.close}`;
        legend.textContent = ``;
    } else {
        let price = param.seriesPrices.get(candleSeries);
        let volume = param.seriesPrices.get(volumeSeries);
        try {
            let OHLC = {
                open: parseFloat(price.open).toFixed(2),
                high: parseFloat(price.high).toFixed(2),
                low: parseFloat(price.low).toFixed(2),
                close: parseFloat(price.close).toFixed(2)
            };
            legend.innerHTML = `<p>O ${OHLC.open} - H ${OHLC.high} - L ${OHLC.low} - C ${OHLC.close} - Vol ${volume}</p>`;
        }
        catch (err) {
            console.log('No data at given time yet.')
        }
    }
});