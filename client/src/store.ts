import { type ComputedRef, computed, reactive, watch } from "vue";
import type { RouteLocation } from "vue-router";
import { router } from "./router";
import { SearchQuery, type SerializedSearchQuery } from "./search";

export type SearchHistory = {
    date: Date;
    query: SearchQuery;
};

export type AccountData = {
    favorite_posts: Post[];
    favorite_tags: Tag[];
    blacklist: Tag[];
    search_history: SearchHistory[];
};

export type AddAccountDataPayload = {
    favorite_posts: Post[];
    favorite_tags: Tag[];
    blacklist: Tag[];
    search_history: SearchHistory[];
};

export type RemoveAccountDataPayload = {
    favorite_post_ids: number[];
    favorite_tag_names: string[];
    blacklist_names: string[];
    search_history: SerializedSearchQuery[];
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

export const APP_VERSION_TTL_MS = 60 * 1000;

export type ColumnSizing = "fixed" | "dynamic";

type Store = {
    account: {
        authToken: string;
        username: string;
        data: AccountData;
    } | null;
    fetchingAccountData: boolean;

    login(username: string, password: string): Promise<void>;
    saveAccountCredentials(): void;
    addToAccountData(data: Partial<AddAccountDataPayload>): Promise<void>;
    removeFromAccountData(data: Partial<RemoveAccountDataPayload>): Promise<void>;
    fetchAccountData(): Promise<void>;

    currentPage: number;
    totalPostCount: number;
    resultsPerPage: number;
    hasSearched: boolean;
    fetchingPosts: boolean;
    justClickedSearchButton: boolean;

    cdnHosts: {
        image: string;
        video: string;
        media_proxy: boolean;
    } | null;

    updateCDNHosts(): void;

    toast: {
        msg: string;
        type: "error" | "info";
    };

    fullscreenPost: Post | null;
    sidebarClosed: boolean;

    userIsSwipingToChangePage: boolean;

    settings: {
        autoplayVideo: boolean;
        blacklist: Tag[];
        checkForUpdates: boolean;
        closeSidebarOnSearch: boolean;
        columnCount: number;
        columnSizing: ColumnSizing;
        columnWidth: number;
        consented: boolean;
        favorites: Post[];
        favoriteTags: Tag[];
        fullscreenViewMenuAnchor: FullscreenViewMenuAnchorPoint;
        fullscreenViewMenuRotate: boolean;
        highResImages: boolean;
        muteVideo: boolean;
        queryHistory: SearchHistory[];
        sidebarTabsHidden: boolean;
        maxPostHeight: number | null;
        newsLastViewedAt: Date;
    };

    loadSettings(): void;
    saveSettings(): void;
    appVersion(): Promise<string>;

    query: SearchQuery;
    lastQuery: SearchQuery;

    lastSearchRoute: RouteLocation | null;

    /** mapping of page number to posts */
    posts: Map<number, Post[]>;
    cachedTags: Map<string, Tag>;
    cachedTagSearch: Map<string, Tag[]>;

    onEditTag: EventTarget;
    onPostsCleared: EventTarget;

    editTag(tag: Tag): void;
    tagsForPost(post: Post): Promise<Tag[]>;
    loadTags(tags: string[]): Promise<void>;
    maxPage(): number;
    nextPage(): Promise<void>;
    postsForCurrentPage(): Post[] | undefined;
    prevPage(): Promise<void>;
    searchPosts(opts: { query: SearchQuery, page?: number; force?: boolean }): Promise<void>;
    addQueryToHistory(): void;
    clearPosts(): void;
    getTag(name: string): Tag | undefined;

    favoritePosts(): ComputedRef<Post[]>;
    addFavoritePosts(posts: Post[]): Promise<void>;
    removeFavoritePosts(postIds: number[]): Promise<void>;

    favoriteTags(): ComputedRef<Tag[]>;
    addFavoriteTags(tags: Tag[]): Promise<void>;
    removeFavoriteTags(tagNames: string[]): Promise<void>;

    blacklist(): ComputedRef<Tag[]>;
    addToBlacklist(tags: Tag[]): Promise<void>;
    removeFromBlacklist(tagNames: string[]): Promise<void>;

    searchHistory(): ComputedRef<SearchHistory[]>;
    addToSearchHistory(history: SearchHistory[]): Promise<void>;
    removeFromSearchHistory(queries: SerializedSearchQuery[]): Promise<void>;
};

const store = reactive<Store>({
    account: JSON.parse(localStorage.getItem("account") || "null"),
    fetchingAccountData: false,

    async login(username: string, password: string): Promise<void> {
        try {
            const resp = await fetch("/api/login", {
                method: "POST",
                body: JSON.stringify({ username, password }),
                headers: { "Content-Type": "application/json" },
            });
            const data = await resp.json();
            if (data.error) {
                this.toast = {
                    msg: data.error,
                    type: "error",
                };
                return;
            }


            this.account = {
                authToken: data.auth_token,
                username: data.username,
                data: {
                    blacklist: [],
                    favorite_posts: [],
                    search_history: [],
                    favorite_tags: [],
                }
            };
            this.saveAccountCredentials();
            await this.fetchAccountData();
        } catch(err) {
            console.error(err);
            this.toast = {
                msg: "Failed to login, is there a connection problem?",
                type: "error",
            };
            throw err;
        }
    },

    saveAccountCredentials() {
        let payload = null;
        if(this.account) {
            const { authToken, username } = this.account;
            payload = { authToken, username }
        }
        localStorage.setItem("account", JSON.stringify(payload));
    },

    async addToAccountData(data: Partial<AddAccountDataPayload>): Promise<void> {
        if(this.account === null) {
            return;
        }

        const { authToken } = this.account;
        const payload = { add: data };

        const resp = await fetch("/api/account/data", {
            body: JSON.stringify(payload),
            method: "PATCH",
            headers: {
                Authorization: `Bearer ${authToken}`,
                "Content-Type": "application/json"
            }
        });
        if(!resp.ok) {
            console.error("error adding data", resp);
            throw new Error("error adding data");
        }
    },

    async removeFromAccountData(data: Partial<RemoveAccountDataPayload>): Promise<void> {
        if(this.account === null) {
            return;
        }

        const { authToken } = this.account;
        const payload = { remove: data };

        const resp = await fetch("/api/account/data", {
            body: JSON.stringify(payload),
            method: "PATCH",
            headers: {
                Authorization: `Bearer ${authToken}`,
                "Content-Type": "application/json"
            }
        });
        if(!resp.ok) {
            console.error("error removing data", resp);
            throw new Error("error removing data");
        }
    },

    async fetchAccountData()  {
        if(!this.account) {
            return;
        }

        // If the data is already available in the HTML, use it directly
        const preloadedData = JSON.parse(document.getElementById("account-data")!.innerText || "null");

        if(preloadedData) {
            this.account.data = preloadedData;
            return;
        }

        const { authToken } = this.account;

        try {
            this.fetchingAccountData = true;

            const resp = await fetch("/api/account/data", {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${authToken}`
                }
            });

            // Token may have expired or account was deleted.
            if(resp.status === 401) {
                this.account = null;
                this.saveAccountCredentials();
                this.toast = {
                    msg: "Please login again",
                    type: "error"
                }
                return;
            }

            if(!resp.ok) {
                console.error("error fetching account data", resp);
                return;
            }

            const data = await resp.json() as AccountData;
            this.account.data = data;
        } catch(e) {
            console.error(e);
            this.toast = {
                msg: "Failed to fetch account data",
                type: "error"
            }
        } finally {
            this.fetchingAccountData = false;
        }
    },

    currentPage: 1,
    totalPostCount: 0,
    resultsPerPage: 0,
    hasSearched: false,
    fetchingPosts: false,
    justClickedSearchButton: false,

    cdnHosts: null,

    updateCDNHosts() {
        this.cdnHosts = JSON.parse(document.getElementById("cdn-hosts")!.innerText);
    },

    toast: {
        msg: "",
        type: "info",
    },

    fullscreenPost: null,
    sidebarClosed: false,

    userIsSwipingToChangePage: false,

    settings: {
        autoplayVideo: true,
        blacklist: [],
        checkForUpdates: true,
        closeSidebarOnSearch: true,
        columnCount: 3,
        columnSizing: "dynamic",
        columnWidth: 400,
        consented: false,
        favorites: [],
        favoriteTags: [],
        fullscreenViewMenuAnchor: "bottomcenter",
        fullscreenViewMenuRotate: false,
        highResImages: true,
        muteVideo: true,
        queryHistory: [],
        sidebarTabsHidden: false,
        maxPostHeight: 600,
        newsLastViewedAt: new Date(0),
    },

    loadSettings() {
        for (const _k in this.settings) {
            const k = _k as keyof typeof this.settings;

            let raw = localStorage.getItem(k);
            if (raw == null) {
                continue;
            }

            // The old settings code didn't stringify columnSizing, fix it
            if (k === "columnSizing" && raw.length > 0 && raw[0] !== '"') {
                raw = JSON.stringify(raw);
                localStorage.setItem(k, raw);
            }

            let val: any;

            // Ignore bad values and just use the default
            try {
                val = JSON.parse(raw);
            } catch (e) {
                console.warn("bad setting value", { k, val, e });
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

                    query.include.forEach(tag => {
                        entry.query.includeTag(tag);
                    });
                    query.exclude.forEach(tag => {
                        entry.query.excludeTag(tag);
                    });

                    return entry;
                });

                val = val.slice(0, QUERY_HISTORY_KEEP_RECENT_LIMIT);
            } else if (k === "newsLastViewedAt") {
                val = new Date(val);
            }

            (this.settings as any)[k] = val;
        }
    },

    saveSettings() {
        Object.entries(this.settings).forEach(([k, v]) => {
            localStorage.setItem(k, JSON.stringify(v));
        });
    },

    async appVersion(): Promise<string> {
        type versionStorage = {
            version: string;
            checkedAt: number;
        }

        const CACHE_KEY = "appversion"
        const cached = localStorage.getItem(CACHE_KEY);

        if(cached) {
            const storedVal = JSON.parse(cached) as versionStorage;

            // Cached version is still fresh
            if(storedVal.checkedAt + APP_VERSION_TTL_MS >= Date.now()) {
                return storedVal.version;
            }
        }

        // Refetch if cache is stale or empty
        const resp = await fetch("/api/version");
        const { version }: { version: string } = await resp.json();
        const val: versionStorage = {
            version,
            checkedAt: Date.now(),
        };
        localStorage.setItem(CACHE_KEY, JSON.stringify(val));

        return version;
    },

    query: new SearchQuery(),
    lastQuery: new SearchQuery(),
    lastSearchRoute: null,
    posts: new Map<number, Post[]>(),
    cachedTags: new Map<string, Tag>(),
    cachedTagSearch: new Map<string, Tag[]>(),

    onEditTag: new EventTarget(),
    onPostsCleared: new EventTarget(),

    editTag(tag: Tag) {
        this.onEditTag.dispatchEvent(new CustomEvent("edit_tag", { detail: tag }));
    },

    async tagsForPost(post: Post): Promise<Tag[]> {
        await store.loadTags(post.tags);
        return post.tags.map(t => this.cachedTags.get(t)).filter(t => t != null);
    },

    async searchPosts({ query, page, force }: { query: SearchQuery, page?: number; force?: boolean; }): Promise<void> {
        type PostListResponse = {
            count_per_page: number;
            total_count: number;
            results: Post[];
        };

        this.fetchingPosts = true;
            page = page ?? this.currentPage;
            const sameQuery = query.equals(this.lastQuery);

            // Don't refetch posts we already have
            if (!force && sameQuery && this.posts.has(page)) {
                this.fetchingPosts = false;
                this.currentPage = page;
                return;
            }

            const searchTags = query.asList().concat(this.blacklist().value.map(t => `-${t.name}`));

            const queryParams = new URLSearchParams();
            queryParams.append("page", page.toString());

            for (const t of searchTags) {
                queryParams.append("q", t);
            }

            const doSearch = async () => {
                try {
                    const resp = await fetch(`/api/posts?${queryParams.toString()}`);
                    if (resp.status >= 400) {
                        try {
                            const val = await resp.json();
                            let msg = "Something went wrong";

                            if ("error" in val) {
                                msg = val.error;
                            }

                            this.toast = {
                                msg,
                                type: "error",
                            };
                        } catch {
                            this.toast = {
                                msg: "Something went wrong",
                                type: "error",
                            };
                        }
                        this.hasSearched = true;
                        throw new Error("search failed");
                    }

                    const json = await resp.json() as PostListResponse;
                    if (!sameQuery) {
                        this.posts.clear();
                    }

                    this.posts.set(page!, json.results);
                    this.resultsPerPage = json.count_per_page;
                    this.totalPostCount = json.total_count;
                    this.currentPage = page!;

                    if (this.settings.closeSidebarOnSearch) {
                        this.sidebarClosed = true;
                    }

                    this.addQueryToHistory();
                    this.lastQuery = this.query.copy();
                } catch(err) {
                    console.error(err);
                    this.toast = {
                        msg: "Something went wrong",
                        type: "error",
                    };
                    throw err;
                } finally {
                    this.fetchingPosts = false;
                    this.hasSearched = true;
                }
            }

            if(!this.fetchingAccountData) {
                await doSearch();
            } else {
                // Wait to search until we're done fetching account data since we need the blacklist
                watch(() => this.fetchingAccountData, doSearch, { once: true });
            }
    },

    maxPage(): number {
        return Math.ceil(this.totalPostCount / this.resultsPerPage);
    },

    async loadTags(tags: string[]): Promise<void> {
        if(tags.length === 0) {
            return;
        }

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
                requests = requests.concat(this.loadTags(tags.slice(start, end)));
            }

            await Promise.all(requests);
            return;
        }

        const missing = tags.filter(t => !this.cachedTags.has(t));

        if (missing.length === 0) {
            return;
        }

        const queryParams = new URLSearchParams();

        for (const t of missing) {
            queryParams.append("t", t);
        }

        try {
            const resp = await fetch(`/api/tags?${queryParams.toString()}`);
            const json = await resp.json() as TagResponse;
            json.results.forEach(t => {
                this.cachedTags.set(t.name, t);
            });
        } catch(err) {
            console.error(err);
            throw err;
        }
    },

    postsForCurrentPage(): Post[] | undefined {
        return this.posts.get(this.currentPage);
    },

    async nextPage(): Promise<void> {
        if (this.currentPage >= this.maxPage()) {
            throw new Error("already at max page");
        }

        await router.push({
            name: "search",
            params: {
                page: this.currentPage + 1,
                query: router.currentRoute.value.params.query,
            },
        });
    },

    async prevPage(): Promise<void> {
        if (this.currentPage <= 1) {
            throw new Error("already at first page");
        }

        await router.push({
            name: "search",
            params: {
                page: this.currentPage - 1,
                query: router.currentRoute.value.params.query,
            },
        });
    },

    addQueryToHistory() {
        if (this.query.isEmpty() || this.query.equals(this.lastQuery)) {
            return;
        }

        this.addToSearchHistory([{
            date: new Date(),
            query: this.query.copy(),
        }]);
    },

    clearPosts() {
        this.posts.clear();
        this.onPostsCleared.dispatchEvent(new CustomEvent("postsCleared"));
    },

    getTag(name: string): Tag | undefined {
        if (name.startsWith("-")) {
            name = name.slice(1);
        }
        return this.cachedTags.get(name);
    },

    favoritePosts(): ComputedRef<Post[]> {
        return computed(() => {
            if (this.account !== null) {
                return this.account.data.favorite_posts;
            }
            return this.settings.favorites;
        })
    },

    async addFavoritePosts(posts: Post[]) {
        if(this.account !== null) {
            return this.addToAccountData({favorite_posts: posts});
        }
        this.settings.favorites = posts.concat(this.settings.favorites);
        this.saveSettings();
    },

    async removeFavoritePosts(postIds: number[]) {
        if(this.account !== null) {
            return this.removeFromAccountData({favorite_post_ids: postIds});
        }
        this.settings.favorites = this.settings.favorites.filter(p => !postIds.includes(p.id));
        this.saveSettings();
    },

    favoriteTags(): ComputedRef<Tag[]> {
        return computed(() => {
            if (this.account !== null) {
                return this.account.data.favorite_tags;
            }
            return this.settings.favoriteTags;
        });
    },

    async addFavoriteTags(tags: Tag[]) {
        if(this.account !== null) {
            return this.addToAccountData({favorite_tags: tags});
        }
        this.settings.favoriteTags = tags.concat(this.settings.favoriteTags);
        this.saveSettings();
    },

    async removeFavoriteTags(tagNames: string[]) {
        if(this.account !== null) {
            return this.removeFromAccountData({favorite_tag_names: tagNames});
        }
        this.settings.favoriteTags = this.settings.favoriteTags.filter(t => !tagNames.includes(t.name));
        this.saveSettings();
    },

    blacklist(): ComputedRef<Tag[]> {
        return computed(() => {
            if (this.account !== null) {
                return this.account.data.blacklist;
            }
            return this.settings.blacklist;
        });
    },

    async addToBlacklist(tags: Tag[]) {
        if(this.account !== null) {
            return this.addToAccountData({blacklist: tags});
        }
        this.settings.blacklist = tags.concat(this.settings.blacklist);
        this.saveSettings();
    },

    async removeFromBlacklist(tagNames: string[]) {
        if(this.account !== null) {
            return this.removeFromAccountData({blacklist_names: tagNames});
        }
        this.settings.blacklist = this.settings.blacklist.filter(t => !tagNames.includes(t.name));
        this.saveSettings();
    },

    searchHistory(): ComputedRef<SearchHistory[]> {
        return computed(() => {
            if (this.account !== null) {
                return this.account.data.search_history;
            }
            return this.settings.queryHistory;
        });
    },

    async addToSearchHistory(history: SearchHistory[]) {
        if(this.account !== null) {
            return this.addToAccountData({search_history: history});
        }
        this.settings.queryHistory = history.concat(this.settings.queryHistory);
        this.saveSettings();
    },

    async removeFromSearchHistory(queries: SerializedSearchQuery[]) {
        if(this.account !== null) {
            return this.removeFromAccountData({search_history: queries});
        }
        this.settings.queryHistory = this.settings.queryHistory.filter(h =>
            !queries.some(q => h.query.equals(new SearchQuery(q)))
        );
        this.saveSettings();
    },
});

export default store;
