import { ArrowForwardIcon } from "@chakra-ui/icons";
import { Button, HStack, Input, Box, VStack } from "@chakra-ui/react";
import React from "react";
import { useProvider } from "../context";
import { Message } from "../models/message";

interface IInputMessage {
  socket: WebSocket | undefined;
  isGroup: boolean;
  toId: string;
  messages: Message[];
  setMessages: React.Dispatch<React.SetStateAction<Message[]>>;
}

export const InputMessage = ({
  socket,
  isGroup,
  toId,
  messages,
  setMessages,
}: IInputMessage) => {
  // Attributes
  const [message, setMessage] = React.useState<string>("");
  // Context
  const { user } = useProvider();
  // Methods
  const GetNewMessageId = () => {
    if (!messages) return 1;
    if (messages.length === 0) return 1;
    return messages[messages.length - 1].id + 1;
  };

  const handleSendMessage = async () => {
    if (!socket) return;
    if (!user) return;
    const obj = {
      user_id: user.id,
      is_group: isGroup,
      to_id: Number(toId),
      message: message,
    };
    const jsonString = JSON.stringify(obj);
    socket.send(jsonString);
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
  // Component
  return (
    <VStack w="full">
      <HStack w="full">
        <Box w="10px" />
        <HStack w="full" bg="black" borderRadius={10}>
          <Input
            placeholder="Write your message here..."
            w="100%"
            bg="black"
            value={message}
            onChange={(e) => setMessage(e.currentTarget.value)}
          />
          <Button variant="signup" w="50px" onClick={handleSendMessage}>
            <ArrowForwardIcon />
          </Button>
        </HStack>
        <Box w="10px" />
      </HStack>
    </VStack>
  );
};
