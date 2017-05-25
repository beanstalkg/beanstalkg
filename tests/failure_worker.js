'use strict';

var fivebeans = require('fivebeans');
var co = require('co');
var bb = require('bluebird');

var client = new fivebeans.client('localhost', 11300);
bb.promisifyAll(client, {multiArgs: true});

client.on('connect', function () {
  co(function* () {
    console.log("connected");
    yield client.watchAsync("test");
    yield client.ignoreAsync("default");
    yield client.useAsync("test");
    yield doit(client);
  }).catch(function(e) {
    console.log(e);
  });
});

client.connect();

function doit(client) {
  return co(function* () {
    var res = yield client.reserveAsync();
    console.log(res[1].toString());
    var ob = JSON.parse(res[1].toString());
    ob.number = ob.number + 1;
    if (ob.number < 4) {
      yield client.destroyAsync(res[0]);
      yield client.putAsync(1, ob.number, 5, JSON.stringify(ob));
    } else {
      yield client.buryAsync(res[0], 1);
    }
    setTimeout(function() {doit(client)}, 1);
  });
}
