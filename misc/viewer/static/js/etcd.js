var warn = false

//设置Modal中的数据
$('#etcd').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget)
    var action = button.data('action')
    var modal = $(this)
    switch(action)
    {
    case 'add':
        modal.find('.modal-title').text('Add an etcd')
        break
    default:
    }
});

function addetcd()
{
    if (warn) {
        $("#warn").addClass("hidden")
    }
     $.ajax({
        url: "addetcd",
        data: {
            "name": $("#name").val(),
            "service-prefix": $("#service-prefix").val(),
            "addresses": $("#addresses").val(),
            "username": $("#username").val(),
            "password": $("#password").val(),
        }, success: function (data) {
            if (data.status == "OK") {
                $("#etcd").modal("hide");
                location.reload()
            } else {
                $("#warn").text(data.error)
                warn = true
                $('#warn').removeClass('hidden')
            }
        }, error: function(data) {
            $("#warn").text(JSON.stringify(data))
            warn = true
            $('#warn').removeClass('hidden')
        }
    });
}

