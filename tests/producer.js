var fivebeans = require('fivebeans');

var client = new fivebeans.client('localhost', 11300);
client
    .on('connect', function()
    {
        console.log("connected");
        client.put(1, 10, 5, JSON.stringify({"number": 1}), function(err, jobid) {
                                console.log("put job id", jobid);
                            });
//        client.use("test", function(err, tubename) {
//            if (err == null) {
//                //setInterval(function() {
//                // priority, delay, ttr, payload
//                    client.put(1, 20, 5, JSON.stringify({"number": 1}), function(err, jobid) {
//                        console.log("put job id", jobid);
//                    });
//                //}, 100);
//            }
//
//        });

    })
    .on('error', function(err)
    {
        // connection failure
    })
    .on('close', function()
    {
        // underlying connection has closed
    })
    .connect();