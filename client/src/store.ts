import { reactive } from "vue";
import { SearchQuery, type SerializedSearchQuery } from "./search";

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
        autoplayVideo: boolean;
        blacklist: Tag[];
        closeSidebarOnSearch: boolean;
        columnCount: number;
        columnSizing: ColumnSizing;
        columnWidth: number;
        consented: boolean;
        favorites: Post[];
        fullscreenViewMenuAnchor: FullscreenViewMenuAnchorPoint;
        fullscreenViewMenuRotate: boolean;
        highResImages: boolean;
        muteVideo: boolean;
        queryHistory: SearchHistory[];
        searchOnPageLoad: PageLoadAutoSearch;
        sidebarTabsHidden: boolean;
    };

    loadSettings(): void;
    saveSettings(): void;

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
        autoplayVideo: true,
        blacklist: [],
        closeSidebarOnSearch: true,
        columnCount: 3,
        columnSizing: "dynamic",
        columnWidth: 400,
        consented: false,
        favorites: [],
        fullscreenViewMenuAnchor: "bottomcenter",
        fullscreenViewMenuRotate: false,
        highResImages: true,
        muteVideo: true,
        queryHistory: [],
        searchOnPageLoad: onPageLoadDefaultValue(),
        sidebarTabsHidden: false,
    },

    loadSettings() {
        for (const _k in this.settings) {
            const k = _k as keyof typeof store.settings;

            let val = JSON.parse(localStorage.getItem(k) ?? "null");
            if (val == null) {
                continue;
            }

            if (k === "queryHistory") {
                type serializedHistory = {
                    date: string;
                    query: SerializedSearchQuery;
                };

                // Transform query history JSON
                val = (val as serializedHistory[]).map(({ date, query }) => {
                    const entry = {
                        date: new Date(date),
                        query: new SearchQuery(),
                    };

                    query.include.forEach((tag) => entry.query.includeTag(tag));
                    query.exclude.forEach((tag) => entry.query.excludeTag(tag));

                    return entry;
                });

                val = val.slice(0, QUERY_HISTORY_KEEP_RECENT_LIMIT);
            }

            (this.settings as any)[k] = val;
        }
    },

    saveSettings() {
        Object.entries(this.settings).forEach(([k, v]) => {
            localStorage.setItem(k, JSON.stringify(v));
        });
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
        this.saveSettings();
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
