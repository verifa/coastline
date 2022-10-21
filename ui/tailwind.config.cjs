const config = {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		extend: {}
	},

	plugins: [
		require('@tailwindcss/forms'),
		require("@tailwindcss/typography"),
		require("daisyui")
	],
	daisyui: {
		themes: ["lofi", "dark"],
	},

};

module.exports = config;
