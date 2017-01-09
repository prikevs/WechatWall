(function(global) {
    var easing, animate, animation, css, pool, run, queue, dequeue, cache, guid, expando;

    guid = 1;

    expando = 'lxjwlt' + (Math.random() + '').replace(/\D/g, '');

    pool = [];

    cache = {};

    // t: current time, b: begInnIng value, c: change In value, d: duration
    easing = {
        def: 'easeOutQuad',
        swing: function (t, b, c, d) {
            //alert(jQuery.easing.default);
            return easing[easing.def](t, b, c, d);
        },
        easeInQuad: function (t, b, c, d) {
            return c*(t/=d)*t + b;
        },
        easeOutQuad: function (t, b, c, d) {
            return -c *(t/=d)*(t-2) + b;
        },
        easeInCubic: function (t, b, c, d) {
            return c*(t/=d)*t*t + b;
        },
        easeOutBack: function (t, b, c, d, s) {
            // if (s == undefined) s = 1.70158;
            if (s == undefined) s = 6;
            return c*((t=t/d-1)*t*((s+1)*t + s) + 1) + b;
        },
        easeOutBounce: function (t, b, c, d) {
            if ((t/=d) < (1/2.75)) {
                return c*(7.5625*t*t) + b;
            } else if (t < (2/2.75)) {
                return c*(7.5625*(t-=(1.5/2.75))*t + .75) + b;
            } else if (t < (2.5/2.75)) {
                return c*(7.5625*(t-=(2.25/2.75))*t + .9375) + b;
            } else {
                return c*(7.5625*(t-=(2.625/2.75))*t + .984375) + b;
            }
        }
    };

    css = function(elem, obj) {
        var prop;

        if (typeof obj === 'object') {
            for (prop in obj) {
                elem.style[prop] = obj[prop];
            }
        } else if (arguments.length === 3) {
            elem.style[arguments[1]] = arguments[2];
        } else {
            return elem.currentStyle ?
                elem.currentStyle[obj] : getComputedStyle(elem, null)[obj];
        }
    };

    run = function(pool, easing) {
        if (pool[0] === 'run' || !pool.length) return;
        pool.unshift('run');

        console.log('run');

        var timeId = setInterval(function() {
            var obj, val, i, t, b, c, d;

            for(i = pool.length - 1; i > 0; i--) {
                obj = pool[i];
                obj['bTime'] = obj['bTime'] || new Date().getTime();

                t = new Date().getTime() - obj['bTime'];
                b = obj['beginVal'];
                c = obj['changeVal'];
                d = obj['duration'];
                type = obj['type'];

                if (t >= d) {
                    val = easing[type](d, b, c, d);
                    obj.over();
                    pool.splice(i, 1);
                } else {
                    val = easing[type](t, b, c, d);
                }

                css(obj['elem'], obj['propName'], val + obj['unit']);
            }

            if (pool.length === 1) {
                clearInterval(timeId);
                pool.pop();
            }

        }, 16);
    };

    animation = function(elem, attr, duration, type, callback) {
        var n = 0,
            beginVal, targetVal, prop, n, unit, cssVal;

        for (prop in attr) {
            cssVal = css(elem, prop);
            beginVal = +parseFloat(cssVal).toFixed(1);
            targetVal = +parseFloat(attr[prop]).toFixed(1);
            unit = +cssVal === beginVal ? '' : cssVal.match(/[a-z]+/);

            if (targetVal !== beginVal) {
                n += 1;
                pool.push({
                    elem: elem,
                    propName: prop,
                    beginVal: beginVal,
                    changeVal: targetVal - beginVal,
                    duration: duration || 400,
                    type: type || 'swing',
                    unit: unit,
                    over: function() {
                        n -= 1;
                        if (n === 0){
                            if(typeof callback == 'function'){
                                callback();
                            }
                            dequeue(elem);
                        }
                    }
                });
            }
        }

        run(pool, easing);
    };

    animate = function(elem, attr, duration, type, callback) {
        var fnc;
        fnc = animation.bind(window, elem, attr, duration, type, callback);
        queue(elem, fnc);
    };

    queue = function(elem, fnc) {
        var theQueue, internalKey;
        if (!elem[expando]) {
            elem[expando] = guid++;
            cache[elem[expando]] = [];
        }
        internalKey = elem[expando];
        theQueue = cache[internalKey];

        theQueue.push(fnc);

        if (theQueue[0] !== 'run') dequeue(elem);
    };

    dequeue = function(elem) {
        var theQueue = cache[elem[expando]],
            state = 'run',
            fnc;

        if (theQueue[0] === state && theQueue.length === 1) {
            theQueue.shift();
            return;
        }

        while (theQueue.length) {
            fnc = theQueue.shift();
            if (typeof fnc === 'function') {
                fnc();
                theQueue.unshift(state);
                break;
            }
        }
    };

    global.animate = animate;
    global.animation = animation;
    global.css = css;

})(window);

// var block = document.getElementById('block'),
//     button = document.getElementById('button'),
//     reset = document.getElementById('reset');
//
// button.addEventListener('click', function(e) {
//     animate(block, {
//         'left': '200px',
//         'opacity': 0.5
//     }, 800, 'easeOutBounce');
//     animate(block, {
//         'top': '100px',
//         'opacity': 1
//     }, 800, 'easeOutBounce');
// });
//
// reset.addEventListener('click', function(e) {
//     css(block, {
//         'top': '0px',
//         'left': '100px'
//     });
// });
