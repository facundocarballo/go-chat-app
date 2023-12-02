import { HStack, VStack, Box, Text, Spacer, Divider } from "@chakra-ui/react";
import React from "react";

interface IFriend {
  name: string;
  email: string;
}
export const Friend = ({ name, email }: IFriend) => {
  // Attributes
  // Context
  // Methods
  // Component
  return (
    <VStack w="300px" bg="primary" borderRadius={10}>
      <Box h="5px" />
      <HStack w="full">
        <Box w="10px" />
        <Text fontWeight="bold">{name}</Text>
        <Spacer />
      </HStack>
      <Divider />
      <HStack w="full">
        <Box w="10px" />
        <Text fontWeight="bold">{email}</Text>
        <Spacer />
      </HStack>
      <Box h="5px" />
    </VStack>
  );
};
