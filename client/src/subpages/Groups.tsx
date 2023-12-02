import { Grid, HStack, Heading, VStack } from "@chakra-ui/react";
import React from "react";
import { useProvider } from "../context";
import { Group } from "../components/Group";
import NextLink from "next/link";

export const Groups = () => {
  // Attributes
  // Context
  const { user } = useProvider();
  // Methods
  // Component
  return (
    <VStack w="full">
      <Heading>Groups</Heading>
      <HStack w="90%">
        <Grid templateColumns="repeat(5, 1fr)" gap={6}>
          {user?.groups?.map((g, idx) => (
            <NextLink key={idx} href={`/group/${g.id}`}>
              <Group name={g.name} description={g.description} />
            </NextLink>
          ))}
        </Grid>
      </HStack>
    </VStack>
  );
};
