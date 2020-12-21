var net = require('net');

const stage = (arg) => {
    const { name } = arg

    const [firstName, lastName] = name.split(' ')
    return { firstName, lastName }
}

const client = net.connect("\\\\.\\pipe\\wt_stream", function() {
    console.log('connected');
})
 
client.on('data', (data) => {
    const stageResponse = stage(data)
    
    client.end(stageResponse);
});

client.on('end', () => {
    console.log('hangup');
})