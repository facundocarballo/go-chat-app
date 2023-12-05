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
import { useRouter } from "next/router";
import { InputMessage } from "@/src/components/InputMessage";
import { ChevronLeftIcon } from "@chakra-ui/icons";
import { Chat } from "@/src/subpages/Chat";
import { Message } from "@/src/models/message";
import NextLink from "next/link";
import { Group } from "@/src/models/group";

export default function Friend() {
  const router = useRouter();
  const url = router.asPath.split("/");
  // Attributes
  const [group, setGroup] = React.useState<Group | undefined>(undefined);
  const [socket, setSocket] = React.useState<WebSocket | undefined>(undefined);
  const [messages, setMessages] = React.useState<Message[]>([]);
  const groupId = url[url.length - 1];
  // Context
  const { user } = useProvider();
  // Methods
  const handleGetFriend = () => {
    if (!user || !user.groups) {
      return;
    }
    for (const group of user.groups) {
      if (group.id.toString() === groupId) {
        setGroup(group);
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
    };
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
        <title>Group</title>
        <meta name="description" content="Golang Chat app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <VStack w="full">
        <Box h="10px" />
        <HStack w="full">
          <Box w="10px" />
          <NextLink href={"/"}>
            <Button bg="bg">
              <ChevronLeftIcon />
            </Button>
          </NextLink>
          <Spacer />
          <Heading>Chat with {group?.name}</Heading>
          <Spacer />
        </HStack>
        <Box h="10px" />
        <Divider />
        <Chat
          toId={groupId}
          isGroup={true}
          socket={socket}
          messages={messages}
          setMessages={setMessages}
        />
        <Box h="100px" />
        <InputMessage
          socket={socket}
          isGroup={true}
          toId={groupId}
          messages={messages}
          setMessages={setMessages}
        />
      </VStack>
    </>
  );
}
