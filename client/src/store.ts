import { reactive } from "vue";
import { SearchQuery, type SerializedSearchQuery } from "./search";

type SettingsKey = keyof Omit<Store["settings"], "save" | "write">;
type SearchHistory = {
    date: Date;
    query: SearchQuery;
};

export type FullscreenViewMenuAnchorPoint =
    | "topleft"
    | "topcenter"
    | "topright"
    | "right"
    | "bottomright"
    | "bottomcenter"
    | "bottomleft"
    | "left";

const QUERY_HISTORY_KEEP_RECENT_LIMIT = 100;

function loadValue<K extends SettingsKey>(
    key: K,
    defaultValue: Store["settings"][K],
    transform: (val: string) => Store["settings"][K],
): Store["settings"][K] {
    const val = localStorage?.getItem(key);
    if (val === null) {
        return defaultValue;
    }
    return transform(val);
}

export type ColumnSizing = "fixed" | "dynamic";
export type PageLoadAutoSearch = "always" | "if-query" | "never";

// Page load setting was changed from a checkbox to a dropdown. Convert the
// existing bool to the new dropdown selection
function onPageLoadDefaultValue(): PageLoadAutoSearch {
    const oldKey = "searchOnLoad";
    const val = localStorage.getItem(oldKey);

    switch (val) {
        case "true":
            return "if-query";
        case "false":
            return "never";
        default:
            return "always";
    }
}

type Store = {
    currentPage: number;
    totalPostCount: number;
    resultsPerPage: number;
    hasSearched: boolean;
    fetchingPosts: boolean;
    postsBeingViewed: "search-results" | "favorites";

    toast: {
        msg: string;
        type: "error" | "info";
    };

    tagMenu: {
        tag: Tag;
        ref: HTMLElement | null;
    } | null;
    fullscreenPost: Post | null;
    sidebarClosed: boolean;

    userIsSwipingToChangePage: boolean;

    settings: {
        consented: boolean;

        columnSizing: ColumnSizing;
        columnCount: number;
        columnWidth: number;

        sidebarTabsHidden: boolean;
        closeSidebarOnSearch: boolean;
        searchOnPageLoad: PageLoadAutoSearch;
        highResImages: boolean;

        autoplayVideo: boolean;
        muteVideo: boolean;

        fullscreenViewMenuAnchor: FullscreenViewMenuAnchorPoint;
        fullscreenViewMenuRotate: boolean;

        blacklist: Tag[];
        favorites: Post[];

        queryHistory: SearchHistory[];

        save(): void;
        write<K extends SettingsKey>(
            key: K,
            transform: (val: Store["settings"][K]) => string,
        ): void;
    };

    query: SearchQuery;
    lastQuery: SearchQuery;

    /** mapping of page number to posts */
    posts: Map<number, Post[]>;
    cachedTags: Map<string, Tag>;

    onPostsCleared: EventTarget;

    tagsForPost(post: Post): Promise<Tag[]>;
    loadTags(tags: string[]): Promise<void>;
    maxPage(): number;
    nextPage(): Promise<void> | null;
    postsForCurrentPage(): Post[] | undefined;
    prevPage(): Promise<void> | null;
    searchPosts(page?: number): Promise<void>;
    setQueryParams(): void;
    addQueryToHistory(): void;
    shouldSearchOnPageLoad(): boolean;
    clearPosts(): void;
};

const store = reactive<Store>({
    currentPage: 1,
    totalPostCount: 0,
    resultsPerPage: 0,
    hasSearched: false,
    fetchingPosts: false,
    postsBeingViewed: "search-results",

    toast: {
        msg: "",
        type: "info",
    },

    tagMenu: null,
    fullscreenPost: null,
    sidebarClosed: false,

    userIsSwipingToChangePage: false,

    settings: {
        consented: loadValue("consented", false, JSON.parse),

        columnSizing: loadValue("columnSizing", "dynamic", (v) => v as any),
        columnCount: loadValue("columnCount", 3, parseInt),
        columnWidth: loadValue("columnWidth", 400, parseInt),

        sidebarTabsHidden: loadValue("sidebarTabsHidden", false, JSON.parse),
        searchOnPageLoad: loadValue(
            "searchOnPageLoad",
            onPageLoadDefaultValue(),
            JSON.parse,
        ),
        closeSidebarOnSearch: loadValue(
            "closeSidebarOnSearch",
            true,
            JSON.parse,
        ),
        highResImages: loadValue("highResImages", true, JSON.parse),

        autoplayVideo: loadValue("autoplayVideo", true, JSON.parse),
        muteVideo: loadValue("muteVideo", true, JSON.parse),

        fullscreenViewMenuAnchor: loadValue(
            "fullscreenViewMenuAnchor",
            "bottomcenter",
            JSON.parse,
        ),
        fullscreenViewMenuRotate: loadValue(
            "fullscreenViewMenuRotate",
            false,
            JSON.parse,
        ),

        blacklist: loadValue("blacklist", [], JSON.parse),
        favorites: loadValue("favorites", [], JSON.parse),

        queryHistory: loadValue("queryHistory", [], (val) => {
            let ret: SearchHistory[] = [];
            const json = JSON.parse(val) as {
                date: string;
                query: SerializedSearchQuery;
            }[];

            for (const data of json) {
                const entry = {
                    date: new Date(data.date),
                    query: new SearchQuery(),
                };

                data.query.include.forEach((tag) =>
                    entry.query.includeTag(tag),
                );
                data.query.exclude.forEach((tag) =>
                    entry.query.excludeTag(tag),
                );

                ret = ret.concat(entry);
            }

            return ret;
        }),

        save() {
            this.write("consented", JSON.stringify);
            this.write("columnSizing", (v) => v);
            this.write("columnCount", (v) => v.toString());
            this.write("columnWidth", (v) => v.toString());
            this.write("sidebarTabsHidden", JSON.stringify);
            this.write("searchOnPageLoad", JSON.stringify);
            this.write("closeSidebarOnSearch", JSON.stringify);
            this.write("highResImages", JSON.stringify);
            this.write("autoplayVideo", JSON.stringify);
            this.write("muteVideo", JSON.stringify);
            this.write("fullscreenViewMenuAnchor", JSON.stringify);
            this.write("fullscreenViewMenuRotate", JSON.stringify);
            this.write("blacklist", JSON.stringify);
            this.write("favorites", JSON.stringify);
            this.write("queryHistory", (val) =>
                JSON.stringify(val.slice(0, QUERY_HISTORY_KEEP_RECENT_LIMIT)),
            );
        },

        write<K extends keyof Store["settings"]>(
            key: K,
            transform: (val: Store["settings"][K]) => string,
        ) {
            localStorage?.setItem(key, transform(this[key]));
        },
    },

    query: new SearchQuery(),
    lastQuery: new SearchQuery(),
    posts: new Map<number, Post[]>(),
    cachedTags: new Map<string, Tag>([
        // Fake rating:* as metadata tags
        [
            "rating:general",
            { name: "rating:general", count: 0, type: "metadata" },
        ],
        [
            "rating:sensitive",
            { name: "rating:sensitive", count: 0, type: "metadata" },
        ],
        [
            "rating:questionable",
            { name: "rating:questionable", count: 0, type: "metadata" },
        ],
        [
            "rating:explicit",
            { name: "rating:explicit", count: 0, type: "metadata" },
        ],
    ]),

    onPostsCleared: new EventTarget(),

    tagsForPost(post: Post): Promise<Tag[]> {
        return new Promise<Tag[]>((resolve, reject) => {
            store
                .loadTags(post.tags)
                .then(() =>
                    resolve(
                        post.tags
                            .map((t) => store.cachedTags.get(t))
                            .filter((t) => t != null),
                    ),
                )
                .catch(reject);
        });
    },

    setQueryParams() {
        const params = new URLSearchParams();
        params.set("page", this.currentPage.toString());
        params.set("q", this.query.asList().join(","));

        const newUrl = new URL(window.location.href);
        newUrl.hash = params.toString();

        if (window.location.href !== newUrl.toString()) {
            // Slight hack: pushState() does not trigger the window hashchange event.
            // This made it easier to add routing without having to redo a lot of logic.
            // This means searchPosts() is called both by the UI (i.e. clicking search)
            // and by the router (page load/forward/back)
            window.history.pushState(null, "", newUrl);
        }
    },

    searchPosts(page?: number): Promise<void> {
        type PostListResponse = {
            count_per_page: number;
            total_count: number;
            results: Post[];
        };

        return new Promise((resolve, reject) => {
            page = page ?? this.currentPage;
            const sameQuery = this.query.equals(this.lastQuery);

            // Don't refetch posts we already have
            if (sameQuery && this.posts.has(page)) {
                this.fetchingPosts = false;
                this.currentPage = page;
                this.setQueryParams();
                resolve();
                return;
            }

            const query = this.query
                .asList()
                .concat(this.settings.blacklist.map((t) => `-${t.name}`));
            const queryParams =
                `q=${encodeURIComponent(query.join(","))}` + `&page=${page}`;

            this.fetchingPosts = true;

            fetch("/api/posts?" + queryParams)
                .then((resp) => {
                    if (resp.status >= 400) {
                        resp.json()
                            .then((val) => {
                                let msg = "Something went wrong";

                                if ("error" in val) {
                                    msg = val["error"];
                                }

                                store.toast = {
                                    msg,
                                    type: "error",
                                };

                                reject();
                            })
                            .catch(() => {
                                store.toast = {
                                    msg: "Something went wrong",
                                    type: "error",
                                };
                                reject();
                            })
                            .finally(() => (store.hasSearched = true));
                        return;
                    }

                    resp.json().then((json: PostListResponse) => {
                        this.posts.set(page!, json.results);
                        this.resultsPerPage = json.count_per_page;
                        this.totalPostCount = json.total_count;
                        this.currentPage = page!;
                        this.postsBeingViewed = "search-results";

                        if (this.settings.closeSidebarOnSearch) {
                            this.sidebarClosed = true;
                        }

                        this.addQueryToHistory();
                        this.lastQuery = this.query.copy();
                        this.setQueryParams();

                        resolve();
                    });
                })
                .catch((err) => {
                    console.error(err);
                    store.toast = {
                        msg: "Something went wrong",
                        type: "error",
                    };
                    reject(err);
                })
                .finally(() => {
                    this.fetchingPosts = false;
                    store.hasSearched = true;
                });
        });
    },

    maxPage(): number {
        return Math.ceil(this.totalPostCount / this.resultsPerPage);
    },

    loadTags(tags: string[]): Promise<void> {
        type TagResponse = {
            results: Tag[];
        };

        const maxTagsPerRequest = 100;

        // Fetch tags in parallel if there are too many for one request
        if (tags.length > maxTagsPerRequest) {
            let requests: Promise<void>[] = [];

            for (let i = 0; i < tags.length; i += maxTagsPerRequest) {
                const start = i;
                const end = i + maxTagsPerRequest;
                requests = requests.concat(
                    this.loadTags(tags.slice(start, end)),
                );
            }

            return new Promise((resolve, reject) =>
                Promise.all(requests)
                    .then(() => resolve())
                    .catch(() => reject()),
            );
        }

        return new Promise((resolve, reject) => {
            const missing = tags.filter((t) => !this.cachedTags.has(t));

            if (missing.length === 0) {
                resolve();
                return;
            }

            fetch("/api/tags?q=" + encodeURIComponent(missing.join(" ")))
                .then((resp) => {
                    resp.json().then((json: TagResponse) => {
                        json.results.forEach((t) =>
                            this.cachedTags.set(t.name, t),
                        );
                        resolve();
                    });
                })
                .catch((err) => {
                    console.error(err);
                    reject();
                });
        });
    },

    postsForCurrentPage(): Post[] | undefined {
        return this.posts.get(this.currentPage);
    },

    nextPage(): Promise<void> | null {
        if (this.currentPage >= this.maxPage()) {
            return null;
        }

        return this.searchPosts(this.currentPage + 1).catch(() => {});
    },

    prevPage(): Promise<void> | null {
        if (this.currentPage <= 1) {
            return null;
        }

        return this.searchPosts(this.currentPage - 1).catch(() => {});
    },

    addQueryToHistory() {
        if (this.query.isEmpty() || this.query.equals(this.lastQuery)) {
            return;
        }

        let entry: SearchHistory | null = null;

        // To avoid cluttering the search history, reuse an existing entry
        // if it exists with the same query
        for (let i = 0; i < this.settings.queryHistory.length; i++) {
            const ithEntry = this.settings.queryHistory[i];

            if (!this.query.equals(ithEntry.query)) {
                continue;
            }

            entry = ithEntry;
            entry.date = new Date();

            // Remove the entry from the list, it will be re-added to the front
            this.settings.queryHistory.splice(i, 1);
            break;
        }

        if (entry === null) {
            entry = reactive({
                date: new Date(),
                query: this.query.copy(),
            });
        }

        // Newest entries are added to the front of the list
        this.settings.queryHistory = [entry].concat(this.settings.queryHistory);
        this.settings.save();
    },

    shouldSearchOnPageLoad(): boolean {
        switch (store.settings.searchOnPageLoad) {
            case "always":
                return true;
            case "if-query":
                return !store.query.isEmpty();
            case "never":
                return false;
            default:
                return false;
        }
    },

    clearPosts() {
        this.posts.clear();
        this.onPostsCleared.dispatchEvent(new CustomEvent("postsCleared"));
    },
});

export default store;
