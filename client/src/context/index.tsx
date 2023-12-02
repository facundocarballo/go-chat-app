import React from "react";
import { IContext } from "./interface";
import { User } from "../models/user";

const Context = React.createContext<IContext>({
  user: undefined,

  setUser: () => {},
});

export const ContextProvider: React.FC<any> = (props: any) => {
  // Attributes
  const [user, setUser] = React.useState<User | undefined>(undefined);

  // Methods
  const values = { user, setUser };

  const memo = React.useMemo(() => values, [user]);

  return <Context.Provider value={memo}>{props.children}</Context.Provider>;
};

export function useProvider(): IContext {
  const context = React.useContext(Context);
  if (!context)
    throw new Error(
      "useProvider have to be called from a component rendered inside of the Context."
    );
  return context;
}
