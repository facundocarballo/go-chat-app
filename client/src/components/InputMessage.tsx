import { ArrowForwardIcon } from "@chakra-ui/icons";
import { Button, HStack, Input, Spacer, Box, VStack } from "@chakra-ui/react";
import React from "react";

export const InputMessage = () => {
  // Attributes
  const [message, setMessage] = React.useState<string>("");
  // Context
  // Methods
  const handleSendMessage = async () => {};
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
