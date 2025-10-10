/**
 * ZeroTrace Design System Component Mapping
 * Maps design system tokens to React components
 */

export const designSystemMapping = {
  // Color Tokens from zerotrace-design-system.json
  colors: {
    light: {
      background: 'rgba(255, 255, 255, 1)',
      foreground: 'rgba(0, 0, 0, 1)',
      card: 'rgba(255, 255, 255, 1)',
      'card-foreground': 'rgba(0, 0, 0, 1)',
      primary: 'rgba(0, 0, 0, 1)',
      'primary-foreground': 'rgba(255, 255, 255, 1)',
      secondary: 'rgba(255, 107, 0, 1)',
      'secondary-foreground': 'rgba(255, 255, 255, 1)',
      muted: 'rgba(245, 245, 245, 1)',
      'muted-foreground': 'rgba(115, 115, 115, 1)',
      accent: 'rgba(255, 107, 0, 1)',
      'accent-foreground': 'rgba(255, 255, 255, 1)',
      destructive: 'rgba(239, 68, 68, 1)',
      'destructive-foreground': 'rgba(255, 255, 255, 1)',
      border: 'rgba(0, 0, 0, 1)',
      input: 'rgba(255, 255, 255, 1)',
      ring: 'rgba(0, 0, 0, 1)'
    },
    dark: {
      background: 'rgba(0, 0, 0, 1)',
      foreground: 'rgba(255, 255, 255, 1)',
      card: 'rgba(0, 0, 0, 1)',
      'card-foreground': 'rgba(255, 255, 255, 1)',
      primary: 'rgba(255, 255, 255, 1)',
      'primary-foreground': 'rgba(0, 0, 0, 1)',
      secondary: 'rgba(255, 107, 0, 1)',
      'secondary-foreground': 'rgba(0, 0, 0, 1)',
      muted: 'rgba(23, 23, 23, 1)',
      'muted-foreground': 'rgba(163, 163, 163, 1)',
      accent: 'rgba(255, 107, 0, 1)',
      'accent-foreground': 'rgba(0, 0, 0, 1)',
      destructive: 'rgba(239, 68, 68, 1)',
      'destructive-foreground': 'rgba(255, 255, 255, 1)',
      border: 'rgba(255, 255, 255, 1)',
      input: 'rgba(0, 0, 0, 1)',
      ring: 'rgba(255, 255, 255, 1)'
    }
  },

  // Typography from design system
  typography: {
    fontFamily: {
      sans: [
        'Space Grotesk',
        '-apple-system',
        'BlinkMacSystemFont',
        'Segoe UI',
        'Roboto',
        'Oxygen',
        'Ubuntu',
        'Cantarell',
        'Helvetica Neue',
        'sans-serif'
      ],
      mono: [
        'JetBrains Mono',
        'Fira Code',
        'Consolas',
        'Monaco',
        'Courier New',
        'monospace'
      ]
    },
    fontSize: {
      xs: ['0.75rem', { lineHeight: '1rem' }],
      sm: ['0.875rem', { lineHeight: '1.25rem' }],
      base: ['1rem', { lineHeight: '1.5rem' }],
      lg: ['1.125rem', { lineHeight: '1.75rem' }],
      xl: ['1.25rem', { lineHeight: '1.75rem' }],
      '2xl': ['1.5rem', { lineHeight: '2rem' }],
      '3xl': ['1.875rem', { lineHeight: '2.25rem' }],
      '4xl': ['2.25rem', { lineHeight: '2.5rem' }],
      '5xl': ['3rem', { lineHeight: '1' }],
      '6xl': ['3.75rem', { lineHeight: '1' }]
    },
    fontWeight: {
      light: '300',
      normal: '400',
      medium: '500',
      semibold: '600',
      bold: '700'
    }
  },

  // Spacing from design system
  spacing: {
    0: '0px',
    0.5: '0.125rem',
    1: '0.25rem',
    1.5: '0.375rem',
    2: '0.5rem',
    2.5: '0.625rem',
    3: '0.75rem',
    3.5: '0.875rem',
    4: '1rem',
    5: '1.25rem',
    6: '1.5rem',
    7: '1.75rem',
    8: '2rem',
    9: '2.25rem',
    10: '2.5rem',
    11: '2.75rem',
    12: '3rem',
    14: '3.5rem',
    16: '4rem',
    20: '5rem',
    24: '6rem',
    28: '7rem',
    32: '8rem',
    36: '9rem',
    40: '10rem',
    44: '11rem',
    48: '12rem',
    52: '13rem',
    56: '14rem',
    60: '15rem',
    64: '16rem',
    72: '18rem',
    80: '20rem',
    96: '24rem'
  },

  // Border radius from design system
  borderRadius: {
    none: '0px',
    sm: '0.125rem',
    DEFAULT: '0.25rem',
    md: '0.375rem',
    lg: '0.5rem',
    xl: '0.75rem',
    '2xl': '1rem',
    full: '9999px'
  },

  // Box shadows from design system
  boxShadow: {
    sm: '2px 2px 0px 0px rgba(0, 0, 0, 1)',
    DEFAULT: '4px 4px 0px 0px rgba(0, 0, 0, 1)',
    md: '6px 6px 0px 0px rgba(0, 0, 0, 1)',
    lg: '8px 8px 0px 0px rgba(0, 0, 0, 1)',
    xl: '12px 12px 0px 0px rgba(0, 0, 0, 1)',
    '2xl': '16px 16px 0px 0px rgba(0, 0, 0, 1)',
    none: 'none',
    orange: '4px 4px 0px 0px rgba(255, 107, 0, 1)',
    'orange-lg': '8px 8px 0px 0px rgba(255, 107, 0, 1)'
  },

  // Component styles from design system
  components: {
    button: {
      default: {
        backgroundColor: 'rgba(0, 0, 0, 1)',
        color: 'rgba(255, 255, 255, 1)',
        borderRadius: '0.25rem',
        padding: '0.75rem 1.5rem',
        fontSize: '0.875rem',
        fontWeight: '700',
        transition: 'all 0.15s ease',
        border: '3px solid rgba(0, 0, 0, 1)',
        boxShadow: '4px 4px 0px 0px rgba(0, 0, 0, 1)',
        textTransform: 'uppercase',
        letterSpacing: '0.025em',
        hover: {
          transform: 'translate(2px, 2px)',
          boxShadow: '2px 2px 0px 0px rgba(0, 0, 0, 1)'
        },
        active: {
          transform: 'translate(4px, 4px)',
          boxShadow: '0px 0px 0px 0px rgba(0, 0, 0, 1)'
        }
      },
      secondary: {
        backgroundColor: 'rgba(255, 107, 0, 1)',
        color: 'rgba(255, 255, 255, 1)',
        border: '3px solid rgba(0, 0, 0, 1)',
        boxShadow: '4px 4px 0px 0px rgba(0, 0, 0, 1)',
        hover: {
          transform: 'translate(2px, 2px)',
          boxShadow: '2px 2px 0px 0px rgba(0, 0, 0, 1)'
        }
      },
      destructive: {
        backgroundColor: 'rgba(239, 68, 68, 1)',
        color: 'rgba(255, 255, 255, 1)',
        border: '3px solid rgba(0, 0, 0, 1)',
        boxShadow: '4px 4px 0px 0px rgba(0, 0, 0, 1)',
        hover: {
          transform: 'translate(2px, 2px)',
          boxShadow: '2px 2px 0px 0px rgba(0, 0, 0, 1)'
        }
      },
      outline: {
        border: '3px solid rgba(0, 0, 0, 1)',
        backgroundColor: 'rgba(255, 255, 255, 1)',
        color: 'rgba(0, 0, 0, 1)',
        boxShadow: '4px 4px 0px 0px rgba(0, 0, 0, 1)',
        hover: {
          backgroundColor: 'rgba(255, 107, 0, 1)',
          color: 'rgba(255, 255, 255, 1)',
          transform: 'translate(2px, 2px)',
          boxShadow: '2px 2px 0px 0px rgba(0, 0, 0, 1)'
        }
      }
    },

    card: {
      base: {
        borderRadius: '0.25rem',
        border: '3px solid rgba(0, 0, 0, 1)',
        backgroundColor: 'rgba(255, 255, 255, 1)',
        color: 'rgba(0, 0, 0, 1)',
        boxShadow: '8px 8px 0px 0px rgba(0, 0, 0, 1)'
      },
      header: {
        padding: '1.5rem 1.5rem 0',
        display: 'flex',
        flexDirection: 'column',
        gap: '0.5rem'
      },
      title: {
        fontSize: '1.5rem',
        fontWeight: '700',
        lineHeight: '1.2',
        letterSpacing: '-0.025em',
        textTransform: 'uppercase'
      },
      content: {
        padding: '1.5rem'
      },
      footer: {
        padding: '0 1.5rem 1.5rem',
        display: 'flex',
        alignItems: 'center',
        borderTop: '3px solid rgba(0, 0, 0, 1)',
        marginTop: '1rem',
        paddingTop: '1rem'
      }
    },

    input: {
      base: {
        display: 'flex',
        height: '2.75rem',
        width: '100%',
        borderRadius: '0.25rem',
        border: '3px solid rgba(0, 0, 0, 1)',
        backgroundColor: 'rgba(255, 255, 255, 1)',
        padding: '0.75rem 1rem',
        fontSize: '0.875rem',
        fontWeight: '500',
        transition: 'all 0.15s',
        boxShadow: '4px 4px 0px 0px rgba(0, 0, 0, 1)',
        focus: {
          outline: 'none',
          boxShadow: '6px 6px 0px 0px rgba(255, 107, 0, 1)',
          borderColor: 'rgba(255, 107, 0, 1)'
        }
      }
    },

    badge: {
      default: {
        backgroundColor: 'rgba(0, 0, 0, 1)',
        color: 'rgba(255, 255, 255, 1)',
        border: '2px solid rgba(0, 0, 0, 1)',
        borderRadius: '0.125rem',
        padding: '0.25rem 0.75rem',
        fontSize: '0.75rem',
        fontWeight: '700',
        textTransform: 'uppercase',
        letterSpacing: '0.05em'
      },
      secondary: {
        backgroundColor: 'rgba(255, 107, 0, 1)',
        color: 'rgba(255, 255, 255, 1)',
        border: '2px solid rgba(0, 0, 0, 1)'
      },
      outline: {
        backgroundColor: 'transparent',
        color: 'rgba(0, 0, 0, 1)',
        border: '2px solid rgba(0, 0, 0, 1)'
      }
    }
  },

  // Animation from design system
  animation: {
    duration: {
      75: '75ms',
      100: '100ms',
      150: '150ms',
      200: '200ms',
      300: '300ms'
    },
    timingFunction: {
      linear: 'linear',
      in: 'cubic-bezier(0.4, 0, 1, 1)',
      out: 'cubic-bezier(0, 0, 0.2, 1)',
      'in-out': 'cubic-bezier(0.4, 0, 0.2, 1)'
    },
    keyframes: {
      shake: {
        '0%, 100%': { transform: 'translateX(0)' },
        '25%': { transform: 'translateX(-4px)' },
        '75%': { transform: 'translateX(4px)' }
      },
      'slide-in': {
        from: { transform: 'translateY(-100%)' },
        to: { transform: 'translateY(0)' }
      }
    }
  },

  // Breakpoints from design system
  breakpoints: {
    sm: '640px',
    md: '768px',
    lg: '1024px',
    xl: '1280px',
    '2xl': '1536px'
  },

  // Z-index from design system
  zIndex: {
    0: '0',
    10: '10',
    20: '20',
    30: '30',
    40: '40',
    50: '50',
    auto: 'auto'
  }
};

// Tailwind CSS class mappings for easy integration
export const tailwindMappings = {
  // Button classes
  button: {
    default: 'bg-black text-white border-3 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-x-1 active:translate-y-1 active:shadow-none transition-all duration-150 ease-in-out',
    secondary: 'bg-orange-500 text-white border-3 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] transition-all duration-150 ease-in-out',
    destructive: 'bg-red-500 text-white border-3 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] transition-all duration-150 ease-in-out',
    outline: 'bg-white text-black border-3 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:bg-orange-500 hover:text-white hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] transition-all duration-150 ease-in-out'
  },

  // Card classes
  card: 'bg-white border-3 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] rounded',

  // Input classes
  input: 'w-full h-11 border-3 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] focus:outline-none focus:border-orange-500 focus:shadow-[6px_6px_0px_0px_rgba(255,107,0,1)] transition-all duration-150',

  // Badge classes
  badge: {
    default: 'bg-black text-white border-2 border-black text-xs font-bold uppercase tracking-wider',
    secondary: 'bg-orange-500 text-white border-2 border-black text-xs font-bold uppercase tracking-wider',
    outline: 'bg-transparent text-black border-2 border-black text-xs font-bold uppercase tracking-wider'
  }
};

export default designSystemMapping;
