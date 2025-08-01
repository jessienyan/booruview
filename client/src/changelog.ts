export type ChangeLogEntry = {
    date: string;
    changes: string[];
};

export const changeLog: ChangeLogEntry[] = [
    {
        date: "2025-08-01",
        changes: [
            "added favorites",
            "fixed tags being obscured when using a vertical menu in fullscreen view",
            "tweaked tag colors to be a bit easier on the eyes",
        ],
    },
    {
        date: "2025-07-29",
        changes: [
            "added swipe gesture to change the page",
            "reduced flashing when changing pages",
        ],
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
        changes: [
            "added settings to customize the menu position when viewing a post in fullscreen",
        ],
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
        changes: [
            "added a notification when making too many requests and triggering the rate limit",
        ],
    },
    {
        date: "2025-07-09",
        changes: ['added links to the "about" tab', "lots of style tweaks"],
    },
    {
        date: "2025-07-08",
        changes: [
            "added a notification when a new version is released",
            "increased results per page from 50 to 100",
        ],
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
        changes: [
            "added this changelog",
            "fixed autocomplete not working when doing an OR search",
        ],
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
        changes: [
            "too lazy to go back and write them all out",
            "check the commit history lol",
        ],
    },
    {
        date: "2025-05-11",
        changes: ["initial release"],
    },
];
