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
import { Chat } from "@/src/subpages/Chat";
import { Message } from "@/src/models/message";
import NextLink from "next/link";

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
    }
    s.onclose = function () {
      alert("Connection has been closed.");
    };
    setSocket(s);
  };

  React.useEffect(() => {
      handleGetFriend();
      handleConnectSocket();
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
          <Heading>Chat with {friend?.name}</Heading>
          <Spacer />
        </HStack>
        <Divider />
        <Chat
          toId={friendId}
          isGroup={false}
          socket={socket}
          messages={messages}
          setMessages={setMessages}
        />
        <Box h={{lg: "100px", base: "10px"}} />
        <InputMessage
          socket={socket}
          isGroup={false}
          toId={friendId}
          messages={messages}
          setMessages={setMessages}
        />
      </VStack>
    </>
  );
}
