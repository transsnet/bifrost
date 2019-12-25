
//Ajax加载
$('[data-toggle="switch"]').click(function(e) {
    var $this = $(this),
        loadurl = "/nodes?name=" + $this.attr('switch-name'),
        targ = $this.attr('data-target');

    $.get(loadurl, function(data) {
        $(targ).html(data);
    });

    $this.tab('show');
    return false;
});

//第一次加载列表
$(document).ready(function(e) {
    $('#service-list').find("li a").first().click();
});
