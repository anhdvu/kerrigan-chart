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

const fetchBook = async () => {
    let response = await fetch('/book');
    if (!response.ok) {
        console.log('Fetch failed!')
    }
    let data = await response.json();
    const chartdata = JSON.parse(data)["chart_data"];
    console.log(chartdata);
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
        title: 'Walls',
        xaxis: {
            title: 'Delta',
            tickmode: 'array',
            tickvals: [0, 10, 20, 30, 40],
            ticktext: chartdata.xticks,
            ticklabelposition: "outside left"
        },
        yaxis: {
            title: 'BTC'
        },
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
    title: 'Walls',
});

fetchBook();
setInterval(() => {
    fetchBook()
}, 30 * 1000);