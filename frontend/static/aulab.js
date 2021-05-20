// Symbols
const btcusdt = 'BTCUSDT';
const ethusdt = 'ETHUSDT';
const bnbusdt = 'BNBUSDT';

// Binance public REST API settings
const timeoffset = 60 * 60 * 24 * 30 * 1000;
let startTime = Date.now() - timeoffset;
const limit = 1000;
const interval = '5m';

// URIs
const btcusdtURI = `https://api.binance.com/api/v3/klines?symbol=${btcusdt}&interval=${interval}&startTime=${startTime}&limit=${limit}`
const ethusdtURI = `https://api.binance.com/api/v3/klines?symbol=${ethusdt}&interval=${interval}&startTime=${startTime}&limit=${limit}`
const bnbusdtURI = `https://api.binance.com/api/v3/klines?symbol=${bnbusdt}&interval=${interval}&startTime=${startTime}&limit=${limit}`

// Chart HTML Elements
const title = document.getElementById('title');
const legend = document.getElementById('legend');

// Utility functions
function convertTime(t) {
    const now = new Date(t);
    return now.toUTCString().substr(17, 8);
}

// Charts
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

const fetchKline = async (candleSeries, volumeSeries, uri) => {
    let volColor = '';
    while (true) {
        let response = await fetch(uri);
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

const btcusdtSentryURI = '/history'
const fetchSentryHistory = async (uri) => {
    let response = await fetch(uri);
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

fetchKline(candleSeries, volumeSeries, btcusdtURI);
fetchSentryHistory(btcusdtSentryURI);




















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