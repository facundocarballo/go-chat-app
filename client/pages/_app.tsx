import type { AppProps } from "next/app";
import { ChakraProvider } from "@chakra-ui/react";
import { ContextProvider } from "@/src/context";
import theme from "../styles/theme";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <ChakraProvider theme={theme}>
      <ContextProvider>
        <Component {...pageProps} />
      </ContextProvider>
    </ChakraProvider>
  );
}
