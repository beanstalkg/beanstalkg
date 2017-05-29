var fivebeans = require('fivebeans');

var args = process.argv.slice(2);
var client = new fivebeans.client('localhost', 11300);
client
    .on('connect', function()
    {
        console.log("connected");
        client.use("test", function(err, tubename) {
            if (err == null) {
              client.kick_job(args[0], function(err, jobid) {
                  console.log("kicking job ", args[0]);
                  process.exit();
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
