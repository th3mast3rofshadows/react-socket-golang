import React, {useState, useEffect} from 'react';
import io from 'socket.io-client';

const socket = io('http://localhost:3001/', {
    reconnectionDelayMax: 10000,
    "transports" : ["websocket"]
});
function App() {
    const [isConnected, setIsConnected] = useState(false);
    const [lastPong, setLastPong] = useState(null);
    const [togglePing, setTogglePing] = useState(false)

    useEffect(() => {
        socket.on("connect", () => {
            console.log('Connected!!!', socket.id)
            setIsConnected(socket.connected)
            socket.emit("echo", {text: "echoooo"}, (response) => {
                console.log(response)
            })
        })

        socket.on("list", (data) => {
            console.log("list")
            setLastPong(data.text + " at:" + new Date().toISOString())
        })

        socket.on("connect_error", (reason) => {
            console.log("Error connecting", reason)
        });

        socket.on("disconnect", (reason) => {
            console.log('Disconnection reason: ' + reason)
            setIsConnected(false)
        })

        return () => {
            console.log('Disconnecting socket...');
            if(socket) socket.disconnect();
        }
    }, []);

    function somePing() {
        setTogglePing(!togglePing)
        if(socket) {
            try {
                socket.emit("get_go_pods", {}, (response) => {
                    console.log("got smth back?")
                    console.log(response)

                    setLastPong(response.text + " at:" + new Date().toISOString())
                });
            } catch (e) {
                console.error(e)
            }
        } else {
            console.log('nu e socket')
        }
    }

    return (
        <div>
            <p>Connected: { '' + isConnected }</p>
            <p>Last pong: { lastPong || '-' }</p>
            <button onClick={ somePing }>Send ping</button>
        </div>
    );
}

export default App;