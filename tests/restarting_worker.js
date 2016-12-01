var fivebeans = require('fivebeans');

setInterval(function () {
    var client = new fivebeans.client('localhost', 11300);
    client
        .on('connect', function()
        {
            console.log("connected");
            client.watch("test", function(err, numwatched) {
                        if (err == null) {
                            client.ignore("default", function(err, numwatched) {
                                if (err == null) {
                                    client.reserve(function(err, jobid, payload) {
                                        console.log("reserved job id", jobid, payload);
                                    });
                                }
                            });
                        }
                    });
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
}, 1000);