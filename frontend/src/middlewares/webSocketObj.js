// const heartbeatInterval = 30 * 1000

class WebSocketObj {
    constructor() {
        this.socket = null
        // this.heartbeatTimer = null
    }

    connect(url) {
        if (!this.socket) {
            this.socket = new WebSocket(url)
            // this.heartbeatTimer = setInterval(() => {
            //     this.socket.send("ping")
            // }, heartbeatInterval)
        }
    }

    disconnect() {
        if (this.socket) {
            console.log("disconnect...")
            this.socket.close()
            this.socket = null
            // clearInterval(this.heartbeatTimer)
            // this.heartbeatTimer = null
        }
    }

    send(message) {
        if (this.socket) {
            this.socket.send(JSON.stringify(message))
        }
    }

    on(eventName, callback) {
        if (this.socket) {
            this.socket.addEventListener(eventName, callback)
        }
    }
}

export { WebSocketObj }
