import Head from "next/head";
import React from "react";
import { useProvider } from "@/src/context";
import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  useDisclosure,
  Divider,
  Box,
} from "@chakra-ui/react";
import { Login } from "@/src/components/Login";
import { Friends } from "@/src/subpages/Friends";
import { Groups } from "@/src/subpages/Groups";

export default function Home() {
  // Attributes
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = React.useRef(null);
  // Context
  const { user } = useProvider();
  // Methods
  // Component
  return (
    <>
      <Head>
        <title>Go Chat</title>
        <meta name="description" content="Golang Chat app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      {/* Alert Dialog - Login / SignUp */}
      <AlertDialog
        isOpen={!user}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent bg="black">
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Login - SignUp
            </AlertDialogHeader>

            <AlertDialogBody>
              <Login />
            </AlertDialogBody>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>

      <Box h="10px" />
      <Friends />
      <Box h="20px" />
      <Divider />
      <Box h="20px" />
      <Groups />
    </>
  );
}
