$(function(){
	var MARGIN = 24.6;

	var socket = new WebSocket("ws://wechat.hk2.prikevs.com/ws/wall");
	socket.onopen = function(event) {
		alert('success');
	}

    var i = 1;
	socket.onmessage = function(event) {
		var data = JSON.parse(event.data);
		console.log(data);
		addmessage(data);
	}

	var dtemp = {img_url: 'img/1.png', Username: 'ma', Content: 'hello'};
    // http://wechat.hk2.prikevs.com
	function addmessage(data) {
        $('.mess-list').append(
            '<li class="mess-item">' +
            '<div class="item-avatar">' +
            '<img src="http://wechat.hk2.prikevs.com' + data.img_url + '" alt="头像" class="avatar">' +
            '</div>' +
            '<div class="item-mess">' +
            '<h2 class="user-name">' + data.username + '</h2>' +
            '<p class="mess-content">' + data.content + '</p>' +
            '</div>' +
            '</li>'
        );
        if($('.mess-list').height() > $(window).height()) {
            // $('.mess-list').css('transform',  'translateY(-' + (MARGIN*i) + 'vh)');
            $('.mess-list').animate({'top': '-' + (MARGIN) + 'vh'}, function () {
                $('.mess-list .mess-item').eq(0).remove();
                $('.mess-list').css('top', 0);
                // $('.mess-list').css('transform',  'translateY(0)');
                // $('.mess-list').css('transform',  'translateY(' + (MARGIN) + 'vh)');
            })
            // animate($('.mess-list')[0], {'transform': '-' + (MARGIN) + 'vh'}, 600, 'swing', function(){
            //     $('.mess-list').css('transform',  'translateY(0)');
            // })
            // i++;
        }
	}

    // setInterval(function(){
    //     addmessage(dtemp);
    //     i++;
    // }, 2000);

})