import { createTheme, MantineColorsTuple } from "@mantine/core";

// Custom cyan/green color palette - #00fa91
const cyan: MantineColorsTuple = [
  "#ccfff0",
  "#99ffe1",
  "#66ffd2",
  "#33ffc3",
  "#00ffb4",
  "#00fa91", // Base color #00fa91
  "#00e082",
  "#00c773",
  "#00ad64",
  "#009455",
];

// Custom blue accent
const blue: MantineColorsTuple = [
  "#dbeafe",
  "#bfdbfe",
  "#93c5fd",
  "#60a5fa",
  "#3b82f6", // Base blue
  "#2563eb",
  "#1d4ed8",
  "#1e40af",
  "#1e3a8a",
  "#1e3a8a",
];

export const theme = createTheme({
  primaryColor: "cyan",
  primaryShade: { light: 7, dark: 7 }, // Use darker shade (7) for better contrast

  colors: {
    cyan,
    blue,
  },

  fontFamily:
    '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',

  headings: {
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
    fontWeight: "700",
    sizes: {
      h1: { fontSize: "2.25rem", lineHeight: "2.5rem" },
      h2: { fontSize: "1.875rem", lineHeight: "2.25rem" },
      h3: { fontSize: "1.5rem", lineHeight: "2rem" },
      h4: { fontSize: "1.25rem", lineHeight: "1.75rem" },
    },
  },

  defaultRadius: "md",

  components: {
    Button: {
      defaultProps: {
        radius: "md",
      },
    },
    Card: {
      defaultProps: {
        radius: "md",
        padding: "lg",
        shadow: "sm",
      },
    },
    TextInput: {
      defaultProps: {
        radius: "md",
      },
    },
    NumberInput: {
      defaultProps: {
        radius: "md",
      },
    },
    Select: {
      defaultProps: {
        radius: "md",
      },
    },
    Table: {
      styles: {
        th: {
          fontWeight: 600,
        },
      },
    },
    NavLink: {
      styles: {
        root: {
          borderRadius: "8px",
        },
      },
    },
  },
});
