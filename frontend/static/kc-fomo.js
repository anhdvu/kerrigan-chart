const title = document.getElementById('title');
const legend = document.getElementById('legend');
function convertTime(t) {
    const now = new Date(t);
    return now.toUTCString().substr(17, 8);
}

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
        text: 'AESXII Chart TEST v0.7.0',
        fontSize: 24,
        horzAlign: 'left',
        vertAlign: 'bottom',
    },
});

// Price series
const candleSeriesConfig = {
    upColor: '#57D918',
    downColor: '#DB4300',
    borderVisible: true,
    wickVisible: true,
    borderColor: '#0F188C',
    wickColor: '#0F188C',
    borderUpColor: '#57D918',
    borderDownColor: '#DB4300',
    wickUpColor: '#57D918',
    wickDownColor: '#DB4300',
}
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

const lineSentryConfig = {
    priceLineVisible: false,
    priceLineColor: 'rgba(0, 150, 235, 0.3)',
    priceLineWidth: 1,
    priceLineStyle: LightweightCharts.LineStyle.SparseDotted,
    priceLineSource: LightweightCharts.PriceLineSource.LastBar,
    lastValueVisible: false,
    color: 'rgba(0, 150, 235, 0.3)',
    lineWidth: 1,
    lineStyle: LightweightCharts.LineStyle.Solid,
}

const zoneSentryConfig = {
    priceLineVisible: false,
    priceLineColor: 'rgba(99, 255, 0, 0.6)',
    priceLineWidth: 1,
    priceLineStyle: LightweightCharts.LineStyle.SparseDotted,
    priceLineSource: LightweightCharts.PriceLineSource.LastBar,
    lastValueVisible: true,
    color: 'rgba(99, 255, 0, 0.6)',
    lineWidth: 2,
    lineStyle: LightweightCharts.LineStyle.Solid,
}

const lineValue = 300;
const sentryLines = [-6, -5, -4, -3, -2, -1, 1, 2, 3, 4, 5, 6]
const sentryLineSeries = []
const sentryLineData = []
// Initialize data for sentryLineSeries and sentryLineData
const sentrySeries = chart.addLineSeries(mainSentrySeriesConfig);
let sentryData = [];
for (let i = 0; i < sentryLines.length; i++) {
    if (sentryLines[i] % 2 === 0) {
        const series = chart.addLineSeries(zoneSentryConfig);
        sentryLineSeries.push(series);
    } else {
        const series = chart.addLineSeries(lineSentryConfig);
        sentryLineSeries.push(series);
    }
    sentryLineData.push([]);
}

// ###########################
// Fetch static data
// ###########################
const fetchKline = async () => {
    const timeoffset = 60 * 60 * 24 * 30 * 1000;
    let startTime = Date.now() - timeoffset;
    const limit = 1000;
    const interval = '5m';
    const symbol = 'BTCUSDT';
    let volColor = '';

    while (true) {
        let response = await fetch(`https://api.binance.com/api/v3/klines?symbol=${symbol}&interval=${interval}&startTime=${startTime}&limit=${limit}`);
        if (!response.ok) {
            console.log("Error while fetching kline =( !!!")
            break;
        }
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

const fetchSentryHistory = async () => {
    let response = await fetch('/history');
    if (!response.ok) {
        console.log('Fetch failed!')
    }
    let data = await response.json();
    for (let e of data) {
        sentryData.push({
            time: e.time,
            value: e.value,
        });

        // for (let elem of sentryLines) {
        //     sentryLineData[sentryLines.indexOf(elem)].push({
        //         time: e.time,
        //         value: e.value + lineValue * elem
        //     });
        // };

        for (let i = 0; i < sentryLineData.length; i++) {
            sentryLineData[i].push({
                time: e.time,
                value: e.value + lineValue * sentryLines[i]
            })
        }
    };
    sentrySeries.setData(sentryData);
    for (let i = 0; i < sentryLineSeries.length; i++) {
        sentryLineSeries[i].setData(sentryLineData[i])
    }
}

fetchKline();
fetchSentryHistory();

// let markers = [];
// markers.push({ time: 1606897800, position: 'aboveBar', color: '#e91e63', shape: 'arrowDown', size: 2, text: 'S @ 19303.39' });
// markers.push({ time: 1606843200, position: 'belowBar', color: '#2196F3', shape: 'arrowUp', size: 2, text: 'B @ 18913.98' });
// candleSeries.setMarkers(markers);

// ###########################
// Websocket to get live data
// ###########################
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
    };
}

// real time data for sentry series
const sentryConnect = () => {
    const sentryURI = 'wss://mooner.dace.dev/ws';
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
            sentrySeries.update({
                time: data.d.t,
                value: data.d.v
            });
            for (let line of sentryLineSeries) {
                line.update({
                    time: data.d.t,
                    value: data.d.v + sentryLines[sentryLineSeries.indexOf(line)] * lineValue
                });

                if (sentryLines[sentryLineSeries.indexOf(line)] % 2 === 0 && sentryLines[sentryLineSeries.indexOf(line)] < 0) {
                    line.setMarkers([{ time: data.d.t, position: 'inBar', color: 'rgba(255, 255, 255, 0.8)', shape: 'arrowUp', size: 0, text: sentryLines[sentryLineSeries.indexOf(line)] / 2 }])
                } else if (sentryLines[sentryLineSeries.indexOf(line)] % 2 === 0 && sentryLines[sentryLineSeries.indexOf(line)] > 0) {
                    line.setMarkers([{ time: data.d.t, position: 'inBar', color: 'rgba(255, 255, 255, 0.8)', shape: 'arrowDown', size: 0, text: sentryLines[sentryLineSeries.indexOf(line)] / 2 }])
                } else {
                    continue
                }
            }
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
