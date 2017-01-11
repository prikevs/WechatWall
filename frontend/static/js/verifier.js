/**
 * Created by xuekaima on 2017/1/8.
 */
$(function(){
    var MARGIN = 24.6;

    var tips_time_keeper;
    var tips = $('.tips-bar'),
        tips_success = $('.tips-bar .tips-success'),
        tips_fail = $('.tips-bar .tips-fail'),
        tips_loading = $('.tips-bar .tips-loading');
    var socket = new WebSocket("wss://hk2.prikevs.com/ws/verifier");
    socket.onopen = function(event) {
        // alert('success');
        tips.css('transform', 'translateY(-100%)');
        tips_loading.css('transform', 'translateY(-100%)');
    }

    var i = 1;
    socket.onmessage = function(event) {
        // console.log(event.data);
        // var data = JSON.parse(event.data);
        // console.log(data);

        var tdata = event.data;
        var data_arr = tdata.split('\n');
        // console.log(data_arr);
        var msg = JSON.parse(data_arr[0]);
        // console.log(msg);
        if(msg['ret_code']) {
            switch(msg['ret_code']){
                case 200:
                    // alert('success');
                    // if(tips_time_keeper) {
                    //     clearTimeout(tips_time_keeper);
                    // }
                    tips_success.css('transform', 'translateY(0)');
                    tips.css('transform', 'translateY(0)');
                    tips_time_keeper = setTimeout(function(){
                        tips.css('transform', 'translateY(-100%)');
                        tips_success.css('transform', 'translateY(-100%)');
                    }, 1000);
                    break;
                case 500:
                    // alert('fail');
                    console.log(msg['err_msg']);
                    // if(tips_time_keeper) {
                    //     clearTimeout(tips_time_keeper);
                    // }
                    tips_fail.css('transform', 'translateY(0)');
                    tips.css('transform', 'translateY(0)');
                    tips_time_keeper = setTimeout(function(){
                        tips.css('transform', 'translateY(-100%)');
                        tips_fail.css('transform', 'translateY(-100%)');
                    }, 1000);
                    break;
                default:
                    alert('other');
                    break;
            }
            $('#' + msg['msg_id']).remove();
        }else {
            var len = data_arr.length;
            for(var i = 0; i < len; i++) {
                var data = JSON.parse(data_arr[i]);
                console.log(data);
                addmessage(data);
            }
        }
    }

    // var dtemp = {img_url: 'img/1.png', username: 'ma', content: 'hello'};
    function addmessage(data) {
        var t = new Date(data.create_time*1000);
        var ct = t.toLocaleDateString() + ' ' + t.toLocaleTimeString();
        $('.mess-list').append(
            '<li class="mess-item" id="' + data.msg_id + '">' +
            '<i class="icon close" title="删除此信息"></i>' +
            '<div class="item-avatar">' +
            '<img src="https://hk2.prikevs.com' + data.img_url + '" alt="头像" class="avatar">' +
            '</div>' +
            '<div class="item-mess">' +
            '<h2 class="user-name">' + data.username + '</h2>' +
            '<p class="mess-content">' + data.content + '</p>' +
            '<p class="ct">' + ct + '</p>' +
            '</div>' +
            '<div class="item-btn">' +
            '<a href="javascript:;" class="btn show-now">立即上墙</a>' +
            '<a href="javascript:;" class="btn">通过审核</a>' +
            '</div>' +
            '</li>'
        );
        setTimer(data.msg_id, data.ttl);
    }

    function setTimer(id, TTL) {
        var timer = setTimeout(function(){
            $('#' + id).remove();
        }, TTL*0.9);
    }

    $('.mess-list').on('click', '.btn', function(){
        var show_now = false;
        if($(this).hasClass('show-now')) show_now = true;
        var data = {
            "msg_id": $(this).parents('.mess-item').attr('id'),
            "verified_time": new Date().getTime(),
            "show_now": show_now
        }
        console.log(data);
        socket.send(JSON.stringify(data));
    })

    $('.mess-list').on('click', '.close', function(){
        $(this).parents('.mess-item').remove();
        var msg = {
            'msg_id': $(this).parents('.mess-item').attr('id'),
            'del': true
        }
        socket.send(JSON.stringify(msg));
    });

    // setInterval(function(){
    // 	addmessage(dtemp);
    // 	i++;
    // }, 2000);

})