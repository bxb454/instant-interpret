const createWebSocket = () => new WebSocket("ws://localhost:8080/ws");

let socket: WebSocket | null = createWebSocket();

interface IMessageEvent extends MessageEvent {
  data: string;
}

let connect = (cb: (msg: IMessageEvent) => void) => {
    console.log("connecting");
  
    if (!socket) {
      console.log("Socket is null, reconnecting...");
      socket = createWebSocket();
    }
  
    socket.onopen = () => {
      console.log("Successfully Connected");
    };
  
    socket.onmessage = (msg: IMessageEvent) => {
      console.log(msg);
      cb(msg);
    };
  
    socket.onclose = (event: CloseEvent) => {
      console.log("Socket Closed Connection: ", event);
      console.log("Was the connection closed cleanly? ", event.wasClean);
      console.log("Code: ", event.code);
  
      console.log("Reconnecting...");
      socket = createWebSocket();
      connect(cb);
    };
  
    socket.onerror = (error: Event) => {
      console.log("Socket Error: ", error);
    };
  };
  
  let sendMsg = (msg: string) => {
    console.log("sending msg: ", msg);
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(msg);
    } else {
      console.log("Can't send message, the WebSocket is not open. Current state is: " + (socket ? socket.readyState : "Socket is null"));
    }
  };
  
  export { connect, sendMsg };