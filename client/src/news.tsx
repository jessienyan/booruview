import type { Component} from "vue";
import { SURVEY_LINK } from "./config";

type Update = {
	date: Date;
	title: string;
	component: Component;
}

// NOTE: add new updates to the top of the list so it's sorted by most recent
const updates: Update[] = [
	{
		date: new Date("2026-03-25T01:27:33+00:00"),
		title: "please login again",
		component: (
			<>
			<p>There were some issues with searching while logged in that should be resolved now. These happened when opening a new search page or refreshing:</p>
				<ul>
			<li>searches sometimes wouldn't use your blacklist</li>
			<li>the page gets stuck or takes a while to load</li>
			</ul>
		<p>Please login again to fix these issues. The site now uses a login cookie, so consider adding an exception if you have cookies disabled. Thanks ඞ</p>
			</>
		)
	},
	{
		date: new Date("2026-03-24T02:03:05+00:00"),
		title: "coming soon™",
		component: (
			<>
				<p>I've been busy but still trying to work on the site when I can. I'm hoping to have these out in the next couple weeks:</p>
				<ul>
					<li>create multiple fav lists</li>
					<li>import favs from gelbooru</li>
					<li>save searches + search feeds (TBD)</li>
					<li><a href="https://codeberg.org/jessienyan/booruview/issues/36" target="_blank" rel="noopener">variant set</a> grouping (combines visually similar images, opt-in feature)</li>
				</ul>
				<p>I have some big plans for this year, more on that later :)</p>
				<p>As always, you can post your feedback on the <a href={SURVEY_LINK} target="_blank">anonymous survey</a> or the <a href="https://codeberg.org/jessienyan/booruview/issues">codeberg repo</a></p>
			</>
		)
	}
];

export default updates;
