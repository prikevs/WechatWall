$(function() {
    $('body').height($(window).height());

    // 判断是否已经开始
    var isBegin = false;
    var TIMEOUT = 500;
    var chosed = [];
    //获取随机数，并且此数不在之前中奖名单里
    function randomNum() {
        var hasChosed = true;
        var nu;
        while(hasChosed){
            nu = Math.floor(Math.random() * avatar.length);
            if(chosed.length == 0)  hasChosed = false;
            for(var a in chosed){
                if (nu == chosed[a]) {
                    break;
                }
                if (a == chosed.length-1) {
                    hasChosed = false;
                }
            }
        }
        return nu;
    }
    // 速度变化的随机滚动抽奖
    function timeChange() {
        num = randomNum();
        $('.avatar').attr('src',avatar[num].src);
        $(".prizer").text(people[num].username);
        immer = setTimeout(function() {
            TIMEOUT -= 50;
            if(TIMEOUT < 25) TIMEOUT = 25;
            timeChange(TIMEOUT);
        }, TIMEOUT);
    }

    // 图片预加载
    var people = [];
    var avatar = [];
    $.ajax({
        type: "GET",
        async: true,
        url: "https://hk2.prikevs.com/lottery/list",
        dataType:'jsonp',
        success:function(result) {
            console.log(result);
            if(result.ret_code == 200) {
                people = result.user_list;
                loadImgs();
            }
        },
        // timeout:3000, //请求超时时间
        error: function(jqXHR){
            console.log("发生错误：" + jqXHR.status);
        }
    });
    var has_bd = false;
    var bd_num;
    function loadImgs(){

        var loadedImages = 0;
        var numImages = people.length;
        for (var i=0;i<numImages;i++) {
            avatar[i] = new Image();
            avatar[i].onload = function() {
                console.log(loadedImages + '/' + numImages);
                if (++loadedImages >= numImages) {
                    // callback
                    alert("图片加载成功");
                    $('.start').click(function(){
                        if(isBegin) return false;
                        isBegin = true;
                        TIMEOUT = 500;
                        timeChange();
                        $.ajax({
                            type: "GET",
                            async: true,
                            url: "https://hk2.prikevs.com/lottery/bd",
                            dataType:'jsonp',
                            success:function(result) {
                                if(result.has_bd) {
                                    has_bd = true;
                                    var open_id = result.openid;
                                    for(var i = 0; i < people.length; i++) {
                                        if(people[i].user_openid == open_id) {
                                            bd_num = i;
                                        }
                                    }
                                }else {
                                    has_bd = false;
                                }
                            },
                            // timeout:3000, //请求超时时间
                            error: function(jqXHR){
                                console.log("发生错误：" + jqXHR.status);
                            }
                        });
                    });

                    $('.stop').click(function(){
                        isBegin = false;
                        clearTimeout(immer);
                        if(has_bd) {
                            $('.avatar').attr('src',avatar[bd_num].src);
                            $(".prizer").text(people[bd_num].username);
                            chosed.push(bd_num);
                        }else {
                            chosed.push(num);
                        }
                    });
                }
            };
            avatar[i].src = 'https://hk2.prikevs.com/' + people[i].img_url;
        }
    }
})