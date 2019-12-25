var nodes

//设置Modal中的数据
$('#detail').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget)
    var action = button.data('action')
    var modal = $(this)
    switch(action)
    {
    case 'add':
        modal.find('.modal-title').text('Add a node')
        setModalBody(modal, "extPath/address", "value", "report", "status")
        break
    case 'modify':
        modal.find('.modal-title').text('Modify')
        setModalBody(modal, "modify addr", "", "", "")
        break
    case 'info':
        modal.find('.modal-title').text('Info')
        var key = button.data('key')
        var src
        for (i in nodes)
        {
            if (nodes[i].Key == key)
            {
                src = nodes[i]
                break
            }
        }
        setModalBody(modal, src.Key, src.Value, src.Report, src.Status)
        modal.find('.modal-body #address').attr("disabled", true)
        modal.find('.modal-footer #save').click(function(){
            newData = {}
            if ($("#value").val() != src.Value)
            {
                newData["value"] = $("#value").val()
            }
            if ($("#report").val() != src.Report)
            {
                newData["report"] = $("#report").val()
            }
            if ($(".selectpicker").val() != src.Status)
            {
                newData["status"] = $(".selectpicker").val()
            }
            if (newData.length == 0)
            {
                return
            }
            newData["root"] = $("#root").val()
            newData["address"] = $("#address").val()
            $.ajax({
                url: "update",
                data: newData, 
                success: function (data) {
                    $("#detail").modal("hide");
                    $('body').removeClass('modal-open');
                    $('.modal-backdrop').remove();
                    $('#service-list').find('[switch-name="long"]').click();
                }
            });
        })
        break
    default:
    }
});

function setModalBody(modal, addr, value, report, status)
{
    modal.find('.modal-body #address').attr("disabled", false)
    modal.find('.modal-body #address').val(addr)
    modal.find('.modal-body #value').val(value)
    modal.find('.modal-body #report').val(report)
    modal.find('.modal-body .selectpicker').selectpicker('val', status)
}

function addnode()
{
     $.ajax({
        url: "addnode",
        data: {
            "root": $("#root").val(),
            "addr": $("#address").val(),
            "value": $("#value").val(),
            "report": $("#report").val(),
            "status": $("#status").val(),
        }, success: function (data) {
            $("#detail").modal("hide");
        }
    });
}

function Alert(errMsg) {
    bootbox.alert({
        buttons: {
            ok: {
                label: '知道了',
                className: 'btn-info'
            }
        },
        message: 'Error: '+errMsg,
        }
    );
}

$('.table-responsive .selectpicker').on('changed.bs.select', function(e) {
    var $selected = $(e.currentTarget).val();
    var $this = $(this)
    $.ajax({
        url: 'status',
        data: {
            'serviceName': $this.attr('data-service-name'),
            'addr': $this.attr('data-addr'),
            'status': $selected,
        }, success: function (data) {
            if (data.status == "OK") {
                $this.attr('data-old-val', $selected)
            } else {
                Alert(data.info)
                $(e.currentTarget).selectpicker('val', $this.attr('data-old-val'))
            }
        }, error: function (data) {
            Alert("set status error. "+data)
            $(e.currentTarget).selectpicker('val', $this.attr('data-old-val'))
        }
    });
})

$('#delete').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget)
    var addr = button.data('addr')
    $(this).find('.modal-body #del-address').val(addr)
})

function delnode()
{
     $.ajax({
        url: "delnode",
        data: {
            "service-name": $("#del-servicename").val(),
            "addr": $("#del-address").val(),
        }, success: function (data) {
            $("#delete").modal("hide");
        }
    });
}


