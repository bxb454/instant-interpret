const socket: WebSocket = new WebSocket("ws://localhost:8080/ws");

interface IMessageEvent extends MessageEvent {
  data: string;
}

let connect = (cb: (msg: IMessageEvent) => void) => {
  console.log("connecting");

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = (msg: IMessageEvent) => {
    console.log(msg);
    cb(msg);
  };

  socket.onclose = (event: CloseEvent) => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = (error: Event) => {
    console.log("Socket Error: ", error);
  };
};

let sendMsg = (msg: string) => {
  console.log("sending msg: ", msg);
  socket.send(msg);
};

export { connect, sendMsg };