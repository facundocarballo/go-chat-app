import React from "react";
import { useProvider } from "../context";
import { Message } from "../models/message";
import { VStack } from "@chakra-ui/react";
import { ShowMessage } from "../components/ShowMessage";

interface IChat {
  socket: WebSocket | undefined;
  isGroup: boolean;
  toId: string;
  messages: Message[];
  setMessages: React.Dispatch<React.SetStateAction<Message[]>>;
}
export const Chat = ({
  socket,
  isGroup,
  toId,
  messages,
  setMessages,
}: IChat) => {
  // Attributes
  // Context
  const { user } = useProvider();
  // Methods
  const handleGetMessages = async () => {
    if (!user) return;
    const m = await user.GetMessages(toId, isGroup);
    if (!m) return;
    setMessages(m);
  };

  const GetNewMessageId = () => {
    if (!messages) return 1;
    if (messages.length === 0) return 1;
    return messages[messages.length - 1].id + 1;
  };

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

  React.useEffect(() => {
    handleGetMessages();
  }, []);
  // Component
  if (!user) return null;
  if (!messages) return null;
  return (
    <VStack w="full" h='690px' scrollBehavior='smooth' overflowY='scroll'>
      {messages?.map((m, idx) => (
        <ShowMessage key={idx} sender={m.sender_id} message={m.message} owner={user.id == m.sender_id} />
      ))}
    </VStack>
  );
};
