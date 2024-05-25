// /** @type {import('tailwindcss').Config} */
// const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
        content: ["../web/static/**/*.{html,js}"],
        theme: {
                extend: {
                        maxWidth: {
                                "1/4": "25%",
                        },
                        colors: {
                                transparent: "transparent",
                                current: "currentColor",
                                text: "#F0E9E9",
                                dark: "#0F0E0F",
                                "dark-99": "#111111",
                                "dark-98": "#141414",
                                "dark-97": "#161616",
                                "dark-96": "#191919",
                                "dark-95": "#1b1b1b",
                                "dark-94": "#1e1e1e",
                                "dark-93": "#202020",
                                "dark-92": "#232323",
                                "dark-91": "#252525",
                                "dark-90": "#282828",
                                "dark-80": "#414141",
                                "dark-70": "#5a5a5a",
                                "dark-60": "#737373",
                                "dark-50": "#8b8b8b",
                                primary: "#C19DD3",
                                secondary: "#52489C",
                                accent: "#4062BB",
                        },
                        fontFamily: {
                                times: ["Times"],
                                body: ['"Fira Code"'],
                        },
                },
        },
        plugins: [],
};
