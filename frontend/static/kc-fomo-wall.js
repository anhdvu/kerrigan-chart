const generateBarColors = (current) => {
    result = [];
    for (let i = 0; i < 50; i++) {
        if (i < 10) {
            result.push('rgb(204, 0, 0)')
        } else if (9 < i && i < 20) {
            result.push('rgb(230, 92, 0)')
        } else if (19 < i && i < 30) {
            result.push('rgb(102, 0, 204)')
        } else if (29 < i && i < 40) {
            result.push('rgb(0, 128, 43)')
        } else {
            result.push('rgb(0, 255, 0)')
        }
    };
    result[current] = 'rgb(77, 148, 255)';
    return result
}
const generateTickText = (values) => {
    const result = []
    for (let value of values) {
        if (value % 600 === 0) {
            result.push('' + (value / 600 * -1))
        } else {
            result.push('')
        }
    }
    return result
}

const fetchBook = async () => {
    let response = await fetch('/book');
    if (!response.ok) {
        console.log('Fetch failed!')
    }
    let data = await response.json();
    const chartdata = JSON.parse(data)["chart_data"];
    let trace = [{
        x: chartdata.bar,
        y: chartdata.value,
        type: 'bar',
        marker: {
            color: generateBarColors(chartdata.current_bar)
        }
    }];
    let layout = {
        width: 1600,
        paper_bgcolor: 'rgba(19, 23, 34, 1)',
        plot_bgcolor: 'rgba(19, 23, 34, 1)',
        font: {
            size: 16,
            color: 'rgba(255, 255, 255, 0.8)'
        },
        title: '',
        xaxis: {
            title: 'Lines',
            tickmode: 'array',
            tickvals: [0, 10, 20, 30, 40],
            ticktext: generateTickText(chartdata.xticks),
            ticklabelposition: "outside"
        },
        yaxis: {
            title: 'BTC'
        },
        annotations: [
            {
                x: chartdata.current_bar,
                y: 40,
                xref: 'x',
                yref: 'y',
                text: 'Current Price',
                showarrow: true,
                arrowhead: 16,
                ax: 0,
                ay: -120
            }
        ]
    };
    Plotly.react(kcwall, trace, layout);
};

const kcwall = document.getElementById("kchart-wall");
Plotly.newPlot(kcwall, [{
    type: 'bar',
    x: [],
    y: []
}], {
    width: 1600,
    paper_bgcolor: 'rgba(19, 23, 34, 1)',
    plot_bgcolor: 'rgba(19, 23, 34, 1)',
    font: {
        size: 16,
        color: 'rgba(255, 255, 255, 0.8)'
    },
    title: '',
    xaxis: {
        title: 'Lines',
        tickmode: 'array',
        tickvals: [0, 10, 20, 30, 40],
        ticktext: [600, 300, 0, -300, -600],
        ticklabelposition: "outside"
    },
    yaxis: {
        title: 'BTC'
    },
});

fetchBook();
setInterval(() => {
    fetchBook()
}, 30 * 1000);