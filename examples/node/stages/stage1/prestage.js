var net = require('net');

const preStage = (arg) => {
    return {
        name: "Test Name"
    }
}

const client = net.connect("\\\\.\\pipe\\wt_stream", function() {
    console.log('connected');
})
 
client.on('data', (data) => {
    const stageResponse = preStage(data)
    
    client.end(stageResponse);
});

client.on('end', () => {
    console.log('hangup');
})