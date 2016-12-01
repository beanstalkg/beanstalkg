var fivebeans = require('fivebeans');

setInterval(function () {
    var client = new fivebeans.client('localhost', 11300);
    client
        .on('connect', function()
        {
            console.log("connected");
            client.use("test", function(err, tubename) {
                if (err == null) {
                    client.put(1, 0, 1, "payload", function(err, jobid) {
                        console.log("put job id", jobid);
                    });
                }
            });
            client.close();
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
}, 100);