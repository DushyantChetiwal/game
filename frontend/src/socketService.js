class SocketService {
  constructor() {
    this.socket = null;
  }

  static getSocket() {
    if (!this.socket) {
      this.socket = new WebSocket("ws://localhost:8080/v1/play");
    }
    return this.socket;
  }

  static closeSocket() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }
}

export default SocketService;