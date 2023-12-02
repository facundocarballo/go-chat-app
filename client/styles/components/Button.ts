export const Button = {
  variants: {
    login: () => ({
      bg: "primary",
      color: "white",
      margin: "2px",
      _hover: {
        boxShadow: "md",
        transform: "scale(1.12)",
        bg: "primary",
      },
    }),
    signup: () => ({
        bg: "secundary",
        color: "white",
        margin: "2px",
        _hover: {
          boxShadow: "md",
          transform: "scale(1.12)",
          bg: "secundary",
        },
      }),
  },
};
