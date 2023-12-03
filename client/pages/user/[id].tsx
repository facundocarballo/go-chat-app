import Head from "next/head";
import React from "react";
import { useProvider } from "@/src/context";
import {
  Box,
  Button,
  Divider,
  HStack,
  Heading,
  Spacer,
  VStack,
} from "@chakra-ui/react";
import { User } from "@/src/models/user";
import { useRouter } from "next/router";
import { InputMessage } from "@/src/components/InputMessage";
import { ChevronLeftIcon } from "@chakra-ui/icons";
import { FriendChat } from "@/src/subpages/FriendChat";
import { Message } from "@/src/models/message";
import NextLink from "next/link";

let readed = false;

export default function Friend() {
  const router = useRouter();
  const url = router.asPath.split("/");
  // Attributes
  const [friend, setFriend] = React.useState<User | undefined>(undefined);
  const [socket, setSocket] = React.useState<WebSocket | undefined>(undefined);
  const [messages, setMessages] = React.useState<Message[]>([]);
  const friendId = url[url.length - 1];
  // Context
  const { user } = useProvider();
  // Methods
  const handleGetFriend = () => {
    if (!user || !user.friends) {
      return;
    }
    for (const f of user.friends) {
      if (f.id.toString() === friendId) {
        setFriend(f);
        return;
      }
    }
  };

  const handleConnectSocket = async () => {
    console.log("ws://localhost:3690/ws?jwt=" + user?.jwt);
    const s = new WebSocket("ws://localhost:3690/ws?jwt=" + user?.jwt);
    s.onclose = function () {
      alert("Connection has been closed.");
    };
    setSocket(s);
  };

  React.useEffect(() => {
    if (!readed) {
      handleGetFriend();
      handleConnectSocket();
      readed = true;
    }
  });

  React.useEffect(() => {
    // handleConnectSocket();
    // console.log("hola...");
    // const socket = io('http://localhost:3690')
    // console.log(socket)
    // // Manejar eventos del WebSocket
    // socket.on("connect", () => {
    //   console.log("Conectado al servidor WebSocket");
    // });
    // socket.on("message", (data) => {
    //   console.log("Mensaje del servidor:", data);
    // });
    // socket.on("disconnect", () => {
    //   console.log("Desconectado del servidor WebSocket");
    // });
    // // Limpiar la conexión al desmontar el componente
    // return () => {
    //   socket.disconnect();
    // };
  }, []);
  // Component
  return (
    <>
      <Head>
        <title>User</title>
        <meta name="description" content="Golang Chat app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <VStack w="full">
        <Box h="10px" />
        <HStack w="full">
          <Box w="10px" />
          <NextLink href={'/'}>
            <Button bg="bg">
              <ChevronLeftIcon />
            </Button>
          </NextLink>
          <Spacer />
          <Heading>Chat with {friend?.name}</Heading>
          <Spacer />
        </HStack>
        <Box h="10px" />
        <Divider />
        <FriendChat
          friendId={friendId}
          socket={socket}
          messages={messages}
          setMessages={setMessages}
        />
        <Box h="100px" />
        <InputMessage
          socket={socket}
          friendId={friendId}
          messages={messages}
          setMessages={setMessages}
        />
      </VStack>
    </>
  );
}
