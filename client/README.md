## Client - NextJS App

To connect our React / NextJS app to our Golang Web Socket Server we just have to follow this guide.

### Connect with the Web Server Socket

The connection with the Web Socket Server is made every time that a user open a Chat. Can be a personal chat, or a group chat.

You will find the code of this implementation in this file 'pages/user/[id].tsx'

#### React Code

```tsx
const [socket, setSocket] = React.useState<WebSocket | undefined>(undefined);

const handleConnectSocket = async () => {
  if (!user) return;
  const s = new WebSocket("ws://localhost:3690/ws");

  s.onopen = function () {
    const obj = {
      user_id: user.id,
      is_jwt: true,
      is_group: false,
      to_id: -1,
      message: user.jwt,
    };
    const jsonString = JSON.stringify(obj);
    s.send(jsonString);
  };

  s.onclose = function () {
    alert("Connection has been closed.");
  };

  setSocket(s);
};

React.useEffect(() => {
  handleConnectSocket();
}, []);
```

#### Step by Step

1. Connect to the Web Socket.
2. Send the JWT to the Web Socket Server.
   > This second step is optional. For this Chat App, we need to know who is the client connected to the Web Server Socket.

### Send messages to the Web Server Socket

Everytime that a user sends a message, is sending a message to the Web Socket.

You will find the code of this implementation in this file 'src/components/InputMessage.tsx'

#### React Code

```tsx
const [message, setMessage] = React.useState<string>("");

const handleSendMessage = async () => {
  if (!socket) return;
  if (!user) return;
  const obj = {
    user_id: user.id,
    is_jwt: false,
    is_group: isGroup,
    to_id: Number(toId),
    message: message,
  };
  const jsonString = JSON.stringify(obj);

  socket.send(jsonString); // Sending the message to the Web Socket.

  setMessages([
    ...messages,
    new Message(
      GetNewMessageId(),
      user.id,
      Number(toId),
      message,
      new Date().toISOString()
    ),
  ]);
  setMessage("");
};
```

### Receive messages from the Web Server Socket

Everytime that a user sends a message, is sending a message to the Web Socket.

You will find the code of this implementation in this file 'src/subpages/Chat.tsx'

#### React Code

```tsx
const handleMessages = async () => {
  if (!socket) return;
  if (!user) return;
  if (!messages) return;
  socket.onmessage = function (e) {
    const json = JSON.parse(e.data);
    setMessages([
      ...messages,
      new Message(
        GetNewMessageId(),
        json.user_id,
        Number(toId),
        json.message,
        new Date().toISOString()
      ),
    ]);
  };
};

React.useEffect(() => {
  handleMessages();
});
```
