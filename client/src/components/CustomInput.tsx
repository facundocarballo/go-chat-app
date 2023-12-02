import React from "react";
import { Box, HStack, Text, VStack, Spacer, Input } from "@chakra-ui/react";

interface ICustomInput {
  title: string;
  placeholder: string;
  type: "text" | "email" | "number" | "password";
  value: string;
  handler: (e: React.ChangeEvent<HTMLInputElement>) => void
}
export const CustomInput = ({ title, placeholder, type, value, handler }: ICustomInput) => {
  // Attributes
  // Context
  // Methods
  // Component
  return (
    <VStack w='full'>
      <HStack w="full">
        <Box w="10px" />
        <Text fontWeight="bold">{title}</Text>
        <Spacer />
      </HStack>
      <Input 
        w='90%'
        placeholder={placeholder}
        type={type}
        value={value}
        onChange={handler}
      />
    </VStack>
  );
};
