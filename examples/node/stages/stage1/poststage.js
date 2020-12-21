var net = require('net');

const postStage = (arg) => arg

const client = net.connect("\\\\.\\pipe\\wt_stream", function() {
    console.log('connected');
})
 
client.on('data', (data) => {
    const stageResponse = postStage(data)
    
    client.end(stageResponse);
});

client.on('end', () => {
    console.log('hangup');
})