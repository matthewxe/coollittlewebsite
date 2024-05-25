/** @type {import('tailwindcss').Config} */
module.exports = {
        content: ["../web/static/**/*.{html,js}"],
        theme: {
                extend: {
                        // animation: {
                        //         flow: "flow 1s ease infinite",
                        // },
                        // keyframes: {
                        //         flow: {
                        //                 "0%, 100%": {
                        //                         background: "linear-gradient(90deg, rgba(2,0,36,1) 0%, rgba(9,9,121,1) 35%, rgba(0,212,255,1) 100%)",
                        //                 },
                        //                 "50%": {
                        //                         background: "linear-gradient(90deg, rgba(0,212,255,1) 0%, rgba(9,9,121,1) 35%, rgba(2,0,36,1) 100%)",
                        //                 },
                        //         },
                        // },
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
                                display: ["Tex Gyre Adventor Regular"],
                                body: ["Ubuntu"],
                        },
                },
        },
        plugins: [],
};
