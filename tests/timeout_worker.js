'use strict';

var fivebeans = require('fivebeans');
var co = require('co');
var bb = require('bluebird');

var client = new fivebeans.client('localhost', 11300);
bb.promisifyAll(client, {multiArgs: true});

client.on('connect', function () {
     co(function* () {
     console.log("connected");
//                yield client.watchAsync("test");
//                yield client.ignoreAsync("default");
//                yield client.useAsync("test");
               yield doit(client);

            }).catch(function(e) {
                console.log(e);
            });
});

client.connect();





    function doit(client) {
        return co(function* () {
             var res = yield client.reserve_with_timeoutAsync(10);
             console.log(res[1].toString());
             var ob = JSON.parse(res[1].toString());
             ob.number = ob.number + 1;
             yield client.destroyAsync(res[0]);
             client.putAsync(1, ob.number, 5, JSON.stringify(ob));
             setTimeout(function() {doit(client)}, 1);
        });
    }
