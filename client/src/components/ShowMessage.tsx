import { HStack, Spacer, Text, Box, VStack, Circle } from "@chakra-ui/react";
import React from "react";

interface IShowMessage {
  owner: boolean;
  message: string;
  sender: number;
}
export const ShowMessage = ({ owner, message, sender }: IShowMessage) => {
  // Attributes
  // Context
  // Methods
  // Component
  if (owner)
    return (
      <HStack w="full">
        <Spacer />
        <HStack bg="green.800" borderRadius={10}>
          <Box w="1px" />
          <Text p="5px">{message}</Text>
          <Box w="1px" />
        </HStack>
        <Box w="10px" />
      </HStack>
    );
  return (
    <HStack w="full">
      <Box w='10px' />
      <Circle size='40px' bg='green'>
        <Text>{sender}</Text>
      </Circle>
      <HStack bg="gray.800" borderRadius={10}>
        <Box w="1px" />
        <Text p="5px">{message}</Text>
        <Box w="1px" />
      </HStack>
      <Spacer />
    </HStack>
  );
};
