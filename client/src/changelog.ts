export type ChangeLogEntry = {
	date: string;
	changes: string[];
};

export const changeLog: ChangeLogEntry[] = [
	{
		date: "2026-02-26",
		changes: ["added user accounts"]
	},
	{
		date: "2026-02-16",
		changes: ["enabled experimental media proxy to workaround gelbooru requiring referrer"]
	},
	{
		date: "2026-01-20",
		changes: ["fixed another case of API returning 500 instead of 503"],
	},
	{
		date: "2026-01-01",
		changes: ["improved stack traces for backend errors", "happy new year"],
	},
	{
		date: "2025-12-13",
		changes: ["fixed favorites not loading"],
	},
	{
		date: "2025-12-06",
		changes: ["added 'open in new tab' button when clicking on tags"],
	},
	{
		date: "2025-12-02",
		changes: ["changed hosting providers (linode -> hetzner)"],
	},
	{
		date: "2025-11-30",
		changes: ["fixed API returning 500 instead of 503 when gelbooru request times out"],
	},
	{
		date: "2025-11-28",
		changes: [
			"fixed rating tags having the wrong color when they are excluded",
			"fixed header icon position on mobile",
			"removed coal generator, you were nice this year 🎅",
		],
	},
	{
		date: "2025-11-18",
		changes: [
			"removed search input getting autofocused when opening the sidebar",
			"changed formatter/linter (prettier -> biomejs)",
		],
	},
	{
		date: "2025-11-14",
		changes: ["fixed videos not loading"],
	},
	{
		date: "2025-11-06",
		changes: ["added desktop support for reordering favorites (experimental)"],
	},
	{
		date: "2025-11-05",
		changes: ["fixed favorited tags being duplicated when importing settings"],
	},
	{
		date: "2025-11-02",
		changes: [
			"moved blacklist into settings tab",
			"added some example filters",
			"changed recent 'search' button to be a hyperlink",
		],
	},
	{
		date: "2025-10-30",
		changes: [
			"added a link to view post on gelbooru",
			"fixed swipe arrow being half off-screen on some mobile browsers",
		],
	},
	{
		date: "2025-10-20",
		changes: ["fixed extra hyphens being added to raw excluded tags", "added a default opt-in nsfw blacklist"],
	},
	{
		date: "2025-10-14",
		changes: ["added a max post height setting"],
	},
	{
		date: "2025-10-07",
		changes: [
			"fixed images not loading on some mobile browsers",
			"fixed videos not changing when viewing in fullscreen and changing to prev/next post",
		],
	},
	{
		date: "2025-10-03",
		changes: ["added country code to caddy logs to help determine where to host server"],
	},
	{
		date: "2025-10-01",
		changes: ["fixed OR search not working"],
	},
	{
		date: "2025-09-29",
		changes: [
			"(EXPERIMENTAL) eu.booruview.com 🇩🇪 new server in Germany. you can share your settings between booruview.com and eu.booruview.com",
			'fixed posts getting "stuck" and not resetting when doing a new search',
			"added keyboard shortcut for changing pages (left arrow / right arrow)",
			"fixed the arrow graphic not appearing when swiping left/right",
		],
	},
	{
		date: "2025-09-19",
		changes: ["increase settings export size limit to 5MB", "revert cloudflare migration, return to monke"],
	},
	{
		date: "2025-09-14",
		changes: ["tweak post API limit"],
	},
	{
		date: "2025-09-13",
		changes: [
			"added a lot of documentation to the help tab",
			"add link to tag chips that opens a new tab with the tag as the search query. only triggers if " +
				"control clicking, middle clicking, or right click > open in new tab",
		],
	},
	{
		date: "2025-09-12",
		changes: [
			"added favorite tags",
			"improve color contrast of tag chips",
			"fix tag names not being cleaned up properly",
			"fix dropdown menu sometimes going off screen",
			"refactor dropdown menus",
		],
	},
	{
		date: "2025-09-11",
		changes: ["added a clear tags button to quickly clear parts of the search query"],
	},
	{
		date: "2025-09-10",
		changes: [
			"scroll position should now be remembered between pages",
			"added a fullscreen loading spinner",
			'removed the "auto search on page load" setting',
			"added vue-router and changed the URLs (any old search URLs should automatically redirect)",
			"fixed some issues with the page not grabbing focus which prevented scrolling via keyboard",
			"thank you to everyone who has submitted feedback!",
		],
	},
	{
		date: "2025-09-09",
		changes: ["fixed thumbnails for video content not loading"],
	},
	{
		date: "2025-09-07",
		changes: [
			"added some subtle animations",
			"adjusted tag chip text to be easier to read",
			"added error message when trying to view posts past page 200 (gelbooru blocks the request)",
		],
	},
	{
		date: "2025-09-05",
		changes: ["upgraded the coal machine"],
	},
	{
		date: "2025-09-03",
		changes: [
			"fixed searches/tag lookups breaking if a tag has a quote in it because apparently gelbooru html escapes them? very cool",
			"gelbooru should now be reported as down if cloudflare returns a 521 error",
			"fixed some cases where the API would return an empty response",
			"added logging to track down any other empty responses",
		],
	},
	{
		date: "2025-08-31",
		changes: [
			"tweaked rate limit",
			"fixed some api endpoints returning 404 instead of 400",
			"fixed api not setting Content-Type header",
			"increased the size limit for settings export",
		],
	},
	{
		date: "2025-08-29",
		changes: [
			"fixed duplicate history entry when searching for the first time",
			"added reminder to disable hd images if loading is slow",
		],
	},
	{
		date: "2025-08-23",
		changes: ["added icon to show which posts are favorited", "cache search suggestions on frontend"],
	},
	{
		date: "2025-08-22",
		changes: ["added survey CTA to sidebar", "fixed alignment of sidebar buttons on mobile"],
	},
	{
		date: "2025-08-21",
		changes: [
			"tweaked rate limit to be more forgiving initially",
			"changed rate limit error message to include time until unbanned",
			"added setting to disable checking for updates",
		],
	},
	{
		date: "2025-08-20",
		changes: [
			"added support for multiple gelbooru API keys to avoid rate limiting",
			"fixed search suggestions breaking if a tag contains a hyphen",
		],
	},
	{
		date: "2025-08-18",
		changes: ["improved gelbooru error handling to avoid vague 'something went wrong' error"],
	},
	{
		date: "2025-08-15",
		changes: ["moved source code to codeberg", "made a list and checked it twice 🎅"],
	},
	{
		date: "2025-08-10",
		changes: ["added a loading spinner to the prev/next footer buttons"],
	},
	{
		date: "2025-08-06",
		changes: ["added search hints for commonly misused filters, e.g. order:favcount instead of sort:score"],
	},
	{
		date: "2025-08-05",
		changes: [
			"fixed rate limit sometimes being triggered when it shouldn't",
			"tweaked rate limit parameters to prevent scraping",
		],
	},
	{
		date: "2025-08-02",
		changes: [
			"added import/export to settings",
			"added tag editing",
			"fixed search tag being flagged as 'raw' if it contained a space",
			"fixed search input losing focus when using arrow keys + enter to choose a suggestion",
			"tweaked help tab and search button to make entering one tag at a time more clear",
		],
	},
	{
		date: "2025-08-01",
		changes: [
			"added favorites",
			"fixed tags being obscured when using a vertical menu in fullscreen view",
			"tweaked tag colors to be a bit easier on the eyes",
			"fixed some lag when changing pages",
			"tweaked min/max zoom level",
		],
	},
	{
		date: "2025-07-29",
		changes: ["added swipe gesture to change the page", "reduced flashing when changing pages"],
	},
	{
		date: "2025-07-28",
		changes: [
			"fixed scroll position not being remembered when using back/forward browser navigation",
			"fixed keyboard nav (pageup/pagedown) not working",
		],
	},
	{
		date: "2025-07-24",
		changes: ["added settings to customize the menu position when viewing a post in fullscreen"],
	},
	{
		date: "2025-07-21",
		changes: [
			'"rating:" filter now appears as a metadata tag',
			"fixed uppercase tags not being registered correctly. the search input now only accepts lowercase letters",
			"changed images to start loading before they are visible (should improve how the site feels on slower connections)",
		],
	},
	{
		date: "2025-07-20",
		changes: [
			"added consent option to only view SFW content (blacklists sensitive, questionable, and explicit ratings)",
			'added an "always" option to auto search on page load (default)',
		],
	},
	{
		date: "2025-07-11",
		changes: ["added a notification when making too many requests and triggering the rate limit"],
	},
	{
		date: "2025-07-09",
		changes: ['added links to the "about" tab', "lots of style tweaks"],
	},
	{
		date: "2025-07-08",
		changes: ["added a notification when a new version is released", "increased results per page from 50 to 100"],
	},
	{
		date: "2025-07-06",
		changes: [
			"autocomplete no longer commits the search input if you're doing an OR search (you have to press 'search' or hit enter now to explicitly add the OR to your search)",
			"added setting for video autoplay and mute",
		],
	},
	{
		date: "2025-07-04",
		changes: ["added this changelog", "fixed autocomplete not working when doing an OR search"],
	},
	{
		date: "2025-07-01",
		changes: [
			'added the "recent" search tab',
			"fixed fullscreen view staying open when using back/forward browser nav",
		],
	},
	{
		date: "2025-06-24",
		changes: ["added blacklist"],
	},
	{
		date: "2025-??-??",
		changes: ["too lazy to go back and write them all out", "check the commit history lol"],
	},
	{
		date: "2025-05-11",
		changes: ["initial release"],
	},
];
