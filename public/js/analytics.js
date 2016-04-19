var get_highcharts = function (tab_id) {
    if (tab_id == "bills-tab") {
        $.get("/Accounts/CsvData?t=expense", function (csv) {
            $('#container-bm').highcharts({
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
            $('#container-by').highcharts({
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
    } else if (tab_id == "account-tab") {
        alert(tab_id)
    } else if (tab_id == "catelog-tab") {
        alert(tab_id)
    }
}

