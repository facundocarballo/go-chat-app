import Head from "next/head";
import React from "react";
import { useProvider } from "@/src/context";
import { Box, Divider, Heading, VStack } from "@chakra-ui/react";
import { User } from "@/src/models/user";
import { useRouter } from "next/router";
import { InputMessage } from "@/src/components/InputMessage";

let readed = false;

export default function Friend() {
  const router = useRouter();
  const url = router.asPath.split("/");
  // Attributes
  const [friend, setFriend] = React.useState<User | undefined>(undefined);
  const friendId = url[url.length - 1];
  // Context
  const { user } = useProvider();
  // Methods
  const handleGetFriend = () => {
    if (!user || !user.friends) {
      return;
    }
    for (const f of user.friends) {
      if (f.id.toString() === friendId) {
        setFriend(f);
        return;
      }
    }
  };

  React.useEffect(() => {
    if (!readed) {
      handleGetFriend();
      readed = true;
    }
  });
  // Component
  return (
    <>
      <Head>
        <title>User</title>
        <meta name="description" content="Golang Chat app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <VStack w="full">
        \
        <Box h="10px" />
        <Heading>Chat with {friend?.name}</Heading>
        <Box h="10px" />
        <Divider />
        <Box h="100px" />
        <InputMessage />
      </VStack>
    </>
  );
}
