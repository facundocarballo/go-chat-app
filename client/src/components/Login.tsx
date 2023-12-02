import React from "react";
import { useProvider } from "../context";
import {
  Box,
  Button,
  Divider,
  HStack,
  Spacer,
  Spinner,
  VStack,
} from "@chakra-ui/react";
import { CustomInput } from "./CustomInput";
import { User } from "../models/user";
import { GetHash } from "../handlers/crypto";

export const Login = () => {
  // Attributes
  const [loading, setLoading] = React.useState<boolean>(false);
  const [email, setEmail] = React.useState<string>("");
  const [password, setPassword] = React.useState<string>("");
  // Context
  const { setUser } = useProvider();
  // Methods
  const handleEmail = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.currentTarget.value);
  };
  const handlePassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.currentTarget.value);
  };
  const handleLogin = async () => {
    setLoading(true);
    const user = await User.Login(email, GetHash(password));
    setUser(user);
    setLoading(false);
  };
  const handleSignUp = async () => {};
  // Component
  return (
    <VStack w="full">
      <CustomInput
        title="Email"
        type="email"
        placeholder="facu@gmail.com"
        value={email}
        handler={handleEmail}
      />
      <Divider />
      <CustomInput
        title="Password"
        type="password"
        placeholder="Write your password"
        value={password}
        handler={handlePassword}
      />
      <Divider />
      <Box h="20px" />
      {loading ? (
        <Spinner />
      ) : (
        <HStack w="full">
          <Spacer />
          <Button variant="login" onClick={handleLogin}>
            Login
          </Button>
          <Button variant="signup" onClick={handleSignUp}>
            Sign Up
          </Button>
          <Box w="20px" />
        </HStack>
      )}
    </VStack>
  );
};
