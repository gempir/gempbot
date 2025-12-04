import { createTheme, MantineColorsTuple } from '@mantine/core';

// Custom purple color palette matching the existing branding
const purple: MantineColorsTuple = [
  '#f3e5f5',
  '#e1bee7',
  '#ce93d8',
  '#ba68c8',
  '#ab47bc',
  '#9c27b0', // Base purple
  '#8e24aa',
  '#7b1fa2',
  '#6a1b9a',
  '#4a148c',
];

// Custom indigo accent
const indigo: MantineColorsTuple = [
  '#e8eaf6',
  '#c5cae9',
  '#9fa8da',
  '#7986cb',
  '#5c6bc0',
  '#3f51b5',
  '#3949ab',
  '#303f9f',
  '#283593',
  '#1a237e',
];

export const theme = createTheme({
  primaryColor: 'purple',
  primaryShade: 7,

  colors: {
    purple,
    indigo,
  },

  fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',

  headings: {
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
    fontWeight: '700',
    sizes: {
      h1: { fontSize: '2.25rem', lineHeight: '2.5rem' },
      h2: { fontSize: '1.875rem', lineHeight: '2.25rem' },
      h3: { fontSize: '1.5rem', lineHeight: '2rem' },
      h4: { fontSize: '1.25rem', lineHeight: '1.75rem' },
    },
  },

  defaultRadius: 'md',

  components: {
    Button: {
      defaultProps: {
        radius: 'md',
      },
    },
    Card: {
      defaultProps: {
        radius: 'md',
        padding: 'lg',
        shadow: 'sm',
      },
    },
    TextInput: {
      defaultProps: {
        radius: 'md',
      },
    },
    NumberInput: {
      defaultProps: {
        radius: 'md',
      },
    },
    Select: {
      defaultProps: {
        radius: 'md',
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
          borderRadius: '8px',
        },
      },
    },
  },
});
