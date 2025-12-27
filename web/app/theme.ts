import { createTheme, type MantineColorsTuple } from "@mantine/core";

// Muted green - darker, more readable
const terminal: MantineColorsTuple = [
  "#d4f4e8",
  "#a9e9d1",
  "#7edeba",
  "#53d3a3",
  "#28c88c", // Primary muted green
  "#20a070",
  "#187854",
  "#105038",
  "#08281c",
  "#000000",
];

// Steel gray - mechanical, industrial
const steel: MantineColorsTuple = [
  "#f8f9fa",
  "#e9ecef",
  "#ced4da",
  "#adb5bd",
  "#868e96",
  "#6c757d",
  "#495057",
  "#343a40",
  "#212529",
  "#16191d",
];

// Accent blue - system highlight
const accent: MantineColorsTuple = [
  "#e7f5ff",
  "#d0ebff",
  "#a5d8ff",
  "#74c0fc",
  "#4dabf7",
  "#339af0",
  "#228be6",
  "#1c7ed6",
  "#1971c2",
  "#1864ab",
];

export const theme = createTheme({
  primaryColor: "terminal",
  primaryShade: 5,

  colors: {
    terminal,
    steel,
    accent,
  },

  fontFamily: "'JetBrains Mono', 'Fira Code', 'SF Mono', Consolas, monospace",

  headings: {
    fontFamily: "'JetBrains Mono', 'Fira Code', 'SF Mono', Consolas, monospace",
    fontWeight: "600",
    sizes: {
      h1: { fontSize: "1.5rem", lineHeight: "2rem" },
      h2: { fontSize: "1.25rem", lineHeight: "1.75rem" },
      h3: { fontSize: "1.125rem", lineHeight: "1.5rem" },
      h4: { fontSize: "1rem", lineHeight: "1.375rem" },
    },
  },

  defaultRadius: 0,

  spacing: {
    xs: "0.25rem",
    sm: "0.5rem",
    md: "0.75rem",
    lg: "1rem",
    xl: "1.5rem",
  },

  components: {
    Button: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        root: {
          fontWeight: 500,
          textTransform: "uppercase" as const,
          letterSpacing: "0.05em",
          fontSize: "0.75rem",
          border: "1px solid",
        },
      },
    },
    Card: {
      defaultProps: {
        radius: 0,
        padding: "md",
      },
      styles: {
        root: {
          border: "1px solid var(--mantine-color-steel-7)",
          backgroundColor: "var(--mantine-color-steel-9)",
        },
      },
    },
    TextInput: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        input: {
          backgroundColor: "var(--mantine-color-steel-9)",
          border: "1px solid var(--mantine-color-steel-6)",
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.875rem",
        },
        label: {
          fontWeight: 500,
          fontSize: "0.75rem",
          textTransform: "uppercase" as const,
          letterSpacing: "0.05em",
          marginBottom: "0.25rem",
        },
      },
    },
    Textarea: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        input: {
          backgroundColor: "var(--mantine-color-steel-9)",
          border: "1px solid var(--mantine-color-steel-6)",
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.875rem",
        },
        label: {
          fontWeight: 500,
          fontSize: "0.75rem",
          textTransform: "uppercase" as const,
          letterSpacing: "0.05em",
          marginBottom: "0.25rem",
        },
      },
    },
    NumberInput: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        input: {
          backgroundColor: "var(--mantine-color-steel-9)",
          border: "1px solid var(--mantine-color-steel-6)",
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.875rem",
        },
        label: {
          fontWeight: 500,
          fontSize: "0.75rem",
          textTransform: "uppercase" as const,
          letterSpacing: "0.05em",
          marginBottom: "0.25rem",
        },
      },
    },
    Select: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        input: {
          backgroundColor: "var(--mantine-color-steel-9)",
          border: "1px solid var(--mantine-color-steel-6)",
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.875rem",
        },
        label: {
          fontWeight: 500,
          fontSize: "0.75rem",
          textTransform: "uppercase" as const,
          letterSpacing: "0.05em",
          marginBottom: "0.25rem",
        },
        dropdown: {
          backgroundColor: "var(--mantine-color-steel-9)",
          border: "1px solid var(--mantine-color-steel-6)",
          borderRadius: 0,
        },
      },
    },
    Checkbox: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        input: {
          borderRadius: 0,
        },
        label: {
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.875rem",
        },
      },
    },
    Switch: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        track: {
          borderRadius: 0,
        },
        thumb: {
          borderRadius: 0,
        },
      },
    },
    Table: {
      styles: {
        table: {
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.8125rem",
        },
        th: {
          fontWeight: 600,
          textTransform: "uppercase" as const,
          letterSpacing: "0.05em",
          fontSize: "0.6875rem",
          borderBottom: "1px solid var(--mantine-color-steel-6)",
          padding: "0.5rem 0.75rem",
        },
        td: {
          borderBottom: "1px solid var(--mantine-color-steel-7)",
          padding: "0.5rem 0.75rem",
        },
      },
    },
    NavLink: {
      styles: {
        root: {
          borderRadius: 0,
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.8125rem",
          padding: "0.5rem 0.75rem",
        },
      },
    },
    Modal: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        content: {
          backgroundColor: "var(--mantine-color-steel-9)",
          border: "1px solid var(--mantine-color-steel-6)",
        },
        header: {
          backgroundColor: "var(--mantine-color-steel-9)",
          borderBottom: "1px solid var(--mantine-color-steel-7)",
        },
        title: {
          fontFamily: "'JetBrains Mono', monospace",
          fontWeight: 600,
          textTransform: "uppercase" as const,
          letterSpacing: "0.05em",
          fontSize: "0.875rem",
        },
      },
    },
    Badge: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        root: {
          fontFamily: "'JetBrains Mono', monospace",
          fontWeight: 500,
          textTransform: "uppercase" as const,
          letterSpacing: "0.05em",
          fontSize: "0.625rem",
        },
      },
    },
    Pagination: {
      styles: {
        control: {
          borderRadius: 0,
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.75rem",
        },
      },
    },
    Tooltip: {
      defaultProps: {
        radius: 0,
      },
      styles: {
        tooltip: {
          fontFamily: "'JetBrains Mono', monospace",
          fontSize: "0.75rem",
          backgroundColor: "var(--mantine-color-steel-8)",
          border: "1px solid var(--mantine-color-steel-6)",
        },
      },
    },
    Loader: {
      defaultProps: {
        type: "bars",
        color: "terminal",
      },
    },
    Divider: {
      styles: {
        root: {
          borderColor: "var(--mantine-color-steel-7)",
        },
      },
    },
    Anchor: {
      styles: {
        root: {
          fontFamily: "'JetBrains Mono', monospace",
        },
      },
    },
    ActionIcon: {
      defaultProps: {
        radius: 0,
      },
    },
  },
});
