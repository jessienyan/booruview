import type { Component } from "vue";
import { RouterLink } from "vue-router";
import { DISCORD_LINK, SURVEY_LINK } from "./config";

type Update = {
    date: Date;
    title: string;
    component: Component;
};

// NOTE: add new updates to the top of the list so it's sorted by most recent
const updates: Update[] = [
    {
        date: new Date("2026-04-25T19:04:00+00:00"),
        title: "tap to nav",
        component: (
            <>
                <p>Small update today, mostly for mobile users.</p>
                <ul>
                    <li>
                        You can change posts by clicking the left/right side of
                        the image instead of having to use the menu buttons.
                        This is on by default and can be disabled in the
                        settings
                    </li>
                    <li>
                        Panning and zooming can now be disabled in the settings
                    </li>
                    <li>Images now fill the entire screen</li>
                </ul>
            </>
        ),
    },
    {
        date: new Date("2026-04-21T19:08:00+00:00"),
        title: "discord server",
        component: (
            <>
                <p>
                    Booruview now has a discord server:{" "}
                    <a href={DISCORD_LINK} target="_blank">
                        {DISCORD_LINK}
                    </a>
                </p>
                <ul>
                    <li>Share and discuss feedback (public or private)</li>
                    <li>Share your favorite art/artists</li>
                    <li>Tell other people that they have bad taste</li>
                </ul>
                <p>
                    The{" "}
                    <a href={SURVEY_LINK} target="_blank">
                        anonymous survey
                    </a>{" "}
                    is still available if you prefer that.
                </p>
            </>
        ),
    },
    {
        date: new Date("2026-04-14T19:36:00+00:00"),
        title: "favorites update",
        component: (
            <>
                <p>Today's update has a few changes for favorites.</p>
                <p>
                    1. You can now{" "}
                    <strong>import favorites from Gelbooru</strong> (check the{" "}
                    <RouterLink to={{ name: "favorites" }}>
                        favs page
                    </RouterLink>
                    ).
                </p>
                <p>
                    2. Favorites now use pages the same as search results. This
                    should fix any lag issues if you have a lot of favs.
                </p>
                <p>
                    3. Favoriting a post while logged in now uses significantly
                    less bandwidth.
                </p>
            </>
        ),
    },
    {
        date: new Date("2026-03-28T20:30:16+00:00"),
        title: "improvements for slow connections",
        component: (
            <>
                <p>
                    Fixed an issue on slower connections where favorites were
                    sometimes not being saved. The site uses less bandwidth and
                    should be a bit more responsive, too.
                </p>
            </>
        ),
    },
    {
        date: new Date("2026-03-25T01:27:33+00:00"),
        title: "please login again",
        component: (
            <>
                <p>I've fixed a couple bugs affecting logged in users:</p>
                <ul>
                    <li>searches sometimes wouldn't use your blacklist</li>
                    <li>the page gets stuck or takes a while to load</li>
                </ul>
                <p>
                    Logging in again will fix the issue (
                    <a
                        href="https://codeberg.org/jessienyan/booruview/pulls/39"
                        target="_blank"
                        rel="noopener"
                    >
                        #39
                    </a>
                    )
                </p>
                <p>
                    Also, the site now uses a login cookie, so consider adding
                    an exception if you have cookies disabled. Thanks ඞ
                </p>
            </>
        ),
    },
    {
        date: new Date("2026-03-24T02:03:05+00:00"),
        title: "coming soon™",
        component: (
            <>
                <p>
                    I've been busy but still trying to work on the site when I
                    can. I'm hoping to have these out in the next couple weeks:
                </p>
                <ul>
                    <li>create multiple fav lists</li>
                    <li>import favs from gelbooru</li>
                    <li>save searches + search feeds (TBD)</li>
                    <li>
                        <a
                            href="https://codeberg.org/jessienyan/booruview/issues/36"
                            target="_blank"
                            rel="noopener"
                        >
                            variant set
                        </a>{" "}
                        grouping (combines visually similar images, opt-in
                        feature)
                    </li>
                </ul>
                <p>
                    I have some big plans for this year, more on that later :)
                </p>
                <p>
                    As always, you can post your feedback on the{" "}
                    <a href={SURVEY_LINK} target="_blank">
                        anonymous survey
                    </a>{" "}
                    or the{" "}
                    <a href="https://codeberg.org/jessienyan/booruview/issues">
                        codeberg repo
                    </a>
                </p>
            </>
        ),
    },
];

export default updates;
