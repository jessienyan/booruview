import type { Component} from "vue";

type Update = {
	date: Date;
	title: string;
	component: Component;
}

// NOTE: add new updates to the top of the list so it's sorted by most recent
const updates: Update[] = [
	{
		date: new Date("2026-03-24T02:03:05+00:00"),
		title: "example title",
		component: <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed sed feugiat eros. Vivamus mattis dolor at nisi finibus condimentum. Maecenas sodales non dui id interdum. Morbi dignissim justo vitae fermentum accumsan. Aliquam faucibus ex ex, a lobortis lectus dignissim quis. Aliquam in viverra mauris.</p>
	}
];

export default updates;
