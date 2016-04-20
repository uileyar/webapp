var get_account_hc = function () {
    $.get("/Accounts/CsvData?t=expense", function (csv) {
        $('#container-ae').highcharts({
            chart: {
                plotBackgroundColor: null,
                plotBorderWidth: null,
                plotShadow: false,
                type: 'pie'
            },
            title: {
                text: '账户支出汇总'
            },
            credits: {
                enabled: false
            },
            tooltip: {
                pointFormat: '{series.name}: <b>￥{point.y};{point.percentage:.1f}%</b>'
            },
            plotOptions: {
                pie: {
                    allowPointSelect: true,
                    cursor: 'pointer',
                    dataLabels: {
                        enabled: true,
                        format: '<b>{point.name}</b>:{point.percentage:.1f}%<br/>￥{point.y}',
                        style: {
                            color: (Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black'
                        }
                    }
                }
            },
            data: {
                csv: csv
            }
        });
    });

    $.get("/Accounts/CsvData?t=income", function (csv) {
        $('#container-ai').highcharts({
            chart: {
                plotBackgroundColor: null,
                plotBorderWidth: null,
                plotShadow: false,
                type: 'pie'
            },
            title: {
                text: '账户收入汇总'
            },
            credits: {
                enabled: false
            },
            tooltip: {
                pointFormat: '{series.name}: <b>￥{point.y};{point.percentage:.1f}%</b>'
            },
            plotOptions: {
                pie: {
                    allowPointSelect: true,
                    cursor: 'pointer',
                    dataLabels: {
                        enabled: true,
                        format: '<b>{point.name}</b>:{point.percentage:.1f}%<br/>￥{point.y}',
                        style: {
                            color: (Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black'
                        }
                    }
                }
            },
            data: {
                csv: csv
            }
        });
    });
};

var get_bills_hc = function () {
    $.get("/Bills/JsonData?t=bm", function (data) {
        $('#container-bm').highcharts({
            chart: {
                zoomType: 'xy',
                selectionMarkerFill: 'rgba(0,0,0, 0.2)',
                type: 'line'
            },
            credits: {
                enabled: false
            },
            title: {
                text: '月汇总'
            },
            xAxis: {
                //categories: data.categories,
                type: 'datetime',
                dateTimeLabelFormats: { // don't display the dummy year
                    month: '%e. %b',
                    year: '%b'
                }
            },

            yAxis: {
                title: {
                    text: '金额 (￥)'
                }
            },
            plotOptions: {
                line: {
                    cursor: 'pointer',
                    dataLabels: {
                        enabled: true,
                        formatter: function () {
                            return '￥' + this.y;
                        }
                    },
                    enableMouseTracking: true
                }
            },
            tooltip: {
                shared: true,
                useHTML: true,
                headerFormat: '<small>{point.key}</small><table>',
                pointFormat: '<tr><td style="color: {series.color}">{series.name}: </td>' +
                '<td style="text-align: right">{point.y}</td></tr>',
                footerFormat: '</table>',
                valueDecimals: 2,
                valuePrefix: '￥'
            },
            series: data.series
        });
    });
};

var set_hc_color = function () {
    Highcharts.getOptions().colors = Highcharts.map(Highcharts.getOptions().colors, function (color) {
        return {
            radialGradient: {cx: 0.5, cy: 0.3, r: 0.7},
            stops: [
                [0, color],
                [1, Highcharts.Color(color).brighten(-0.3).get('rgb')] // darken
            ]
        };
    });
};

var get_highcharts = function (tab_id) {

    if (tab_id == "bills-tab") {
        get_bills_hc();
    } else if (tab_id == "account-tab") {
        get_account_hc();
    } else if (tab_id == "catelog-tab") {
        get_catelog_hc();
    }
};

