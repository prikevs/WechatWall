/**
 * Created by xuekaima on 2017/1/8.
 */
$(function(){
    var MARGIN = 24.6;

    var socket = new WebSocket("ws://wechat.hk2.prikevs.com/ws/verifier");
    socket.onopen = function(event) {
        alert('success');
    }

    var i = 1;
    socket.onmessage = function(event) {
        var data = JSON.parse(event.data);
        console.log(data);
        if(data['ret_code']) {
            switch(data['ret_code']){
                case 200:
                    alert('success');
                    break;
                case 500:
                    alert('fail');
                    console.log(data['err_msg']);
                    break;
                default:
                    alert('other');
                    break;
            }
            $('#' + data['msg_id']).remove();
        }else {
            addmessage(data);
        }
    }

    // var dtemp = {img_url: 'img/1.png', username: 'ma', content: 'hello'};
    function addmessage(data) {
        $('.mess-list').append(
            '<li class="mess-item" id="' + data.msg_id + '">' +
            '<i class="icon close"></i>' +
            '<div class="item-avatar">' +
            '<img src="http://wechat.hk2.prikevs.com' + data.img_url + '" alt="头像" class="avatar">' +
            '</div>' +
            '<div class="item-mess">' +
            '<h2 class="user-name">' + data.username + '</h2>' +
            '<p class="mess-content">' + data.content + '</p>' +
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
    })

    // setInterval(function(){
    // 	addmessage(dtemp);
    // 	i++;
    // }, 2000);

})