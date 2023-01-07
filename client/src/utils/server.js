import io from 'socket.io-client';

const setupClient = (namespace = "") => {
    let socket = io('http://localhost:3001/' + namespace, {
        "transports" : ["websocket"]
    });

    socket.on("connect", (cb) => {
        console.log('Connected to namespace: ' + namespace , socket.id)
        cb()
        setIsConnected(socket.connected)
        console.log('Getting pods...')
        socket.emit("go_get_pods", "", (data) => {
            podNo.current = data.total
            generateMoles(data.podList)
        })
    })
}

export default setupClient
