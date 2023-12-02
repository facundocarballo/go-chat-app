import React from "react";
import { HStack, VStack, Box, Text, Spacer, Divider } from "@chakra-ui/react";

interface IGroup {
  name: string;
  description: string;
}
export const Group = ({ name, description }: IGroup) => {
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
        <Text fontWeight="bold">{description}</Text>
        <Spacer />
      </HStack>
      <Box h="5px" />
    </VStack>
  );
};
