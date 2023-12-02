import { Grid, HStack, Heading, VStack } from "@chakra-ui/react";
import React from "react";
import { useProvider } from "../context";
import { Friend } from "../components/Friend";
import NextLink from 'next/link'

export const Friends = () => {
  // Attributes
  // Context
  const { user } = useProvider();
  // Methods
  // Component
  return (
    <VStack w="full">
      <Heading>Friends</Heading>
      <HStack w="90%">
        <Grid templateColumns="repeat(5, 1fr)" gap={6}>
          {user?.friends?.map((f, idx) => (
            <NextLink key={idx} href={`/user/${f.id}`}>
                <Friend name={f.name} email={f.email} />
            </NextLink>
          ))}
        </Grid>
      </HStack>
    </VStack>
  );
};
