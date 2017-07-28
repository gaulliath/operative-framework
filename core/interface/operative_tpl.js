var api_url = 'http://127.0.0.1:{{PORT_API}}';
function call_api(command, arguments, fn) {
    if(arguments == undefined)
        arguments = [];
    console.log(arguments)
    var result = false;
    if (arguments.length < 1) {
        console.log('command')
        $.post(api_url, {'exec': command}, function (data) {
            console.log('ok')
            if(command == 'list_module') {
                if ($('#front') != undefined) {
                    if (data.data.module_list.length > 0) {
                        $('#front').html('');
                        $.each(data.data.module_list, function (index, value) {
                            var tpl = '<div class="module">\
                                <p class="module_title">' + value['name'] + '</p>\
                            <p class="module_description">' + value['description'] + '</p>\
                            <div class="right_module">\
                                <a href="module.php?use='+value['name']+'"><button class="module_use">use this module</button></a>\
                            </div>\
                            </div>';
                            $('#front').append(tpl)
                        })
                    }
                }
            }
        })
    }
    else{
        arguments['exec'] = command;
        $.post(api_url, arguments, function(data){
            if(data.status == 'ERROR')
            {
                fn(false)
            }
            else if(data.status == "OK") {
                fn(data)
            }
        });
    }
}