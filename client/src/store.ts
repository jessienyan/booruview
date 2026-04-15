import { type ComputedRef, computed, reactive, watch } from "vue";
import type { RouteLocation } from "vue-router";
import { router } from "./router";
import { SearchQuery, type SerializedSearchQuery, type SimpleSerializedSearchQuery } from "./search";

export type SearchHistory = {
    date: Date;
    query: SearchQuery;
};

export type SavedSearch = SearchHistory;

export type AccountData = {
    favorite_posts: Post[];
    favorite_tags: Tag[];
    blacklist: Tag[];
    search_history: SearchHistory[];
    saved_searches: SavedSearch[];
};

export type AddAccountDataPayload = {
    favorite_posts: Post[];
    favorite_tags: Tag[];
    blacklist: Tag[];
    search_history: {
        date: string;
        query: SerializedSearchQuery;
    }[];
    saved_searches: {
        date: string;
        query: SerializedSearchQuery;
    }[];
};

export type RemoveAccountDataPayload = {
    favorite_post_ids: number[];
    favorite_tag_names: string[];
    blacklist_names: string[];
    search_queries: SimpleSerializedSearchQuery[];
    saved_queries: SimpleSerializedSearchQuery[];
};

export type PostListResponse = {
    count_per_page: number;
    total_count: number;
    results: Post[];
};

type AccountDataResponse = {
    favorite_posts: Post[];
    favorite_tags: Tag[];
    blacklist: Tag[];
    search_history: {
        date: string;
        query: {
            include: Tag[];
            exclude: Tag[];
        };
    }[];
    saved_searches: {
        date: string;
        query: {
            include: Tag[];
            exclude: Tag[];
        };
    }[];
};

function parseAccountDataFromAPI(resp: Partial<AccountDataResponse>): Partial<AccountData> {
    const parsed: Partial<AccountData> = {};

    if(resp.favorite_posts != null)
        parsed.favorite_posts = resp.favorite_posts;

    if(resp.favorite_tags != null)
        parsed.favorite_tags = resp.favorite_tags;

    if(resp.blacklist != null)
        parsed.blacklist = resp.blacklist;

    if(resp.search_history != null)
        parsed.search_history = resp.search_history.map(h => ({
            date: new Date(h.date),
            query: new SearchQuery(h.query),
        }));

    if(resp.saved_searches != null)
        parsed.saved_searches = resp.saved_searches.map(h => ({
            date: new Date(h.date),
            query: new SearchQuery(h.query),
        }));

    return parsed;
}

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
        data?: AccountData;
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
        favHeaderOpen: boolean;
        favoriteTags: Tag[];
        favorites: Post[];
        fullscreenViewMenuAnchor: FullscreenViewMenuAnchorPoint;
        fullscreenViewMenuRotate: boolean;
        highResImages: boolean;
        maxPostHeight: number | null;
        muteVideo: boolean;
        newsLastViewedAt: Date;
        queryHistory: SearchHistory[];
        savedSearches: SearchHistory[];
    };

    loadSettings(): void;
    saveSettings(): void;
    appVersion(): Promise<string>;

    query: SearchQuery;
    lastQuery: SearchQuery;
    lastSearchRoute: RouteLocation | null;
    lastFavPage: number;

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
    searchPosts(tags: SimpleSerializedSearchQuery, page: number): Promise<PostListResponse>;
    searchAndUpdateResults(opts: { query: SearchQuery, page?: number; force?: boolean }): Promise<void>;
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
    removeFromSearchHistory(queries: SimpleSerializedSearchQuery[]): Promise<void>;

    savedSearches(): ComputedRef<SavedSearch[]>;
    addToSavedSearches(searches: SavedSearch[]): Promise<void>;
    removeFromSavedSearches(queries: SimpleSerializedSearchQuery[]): Promise<void>;
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
                    saved_searches: [],
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
        if(!this.account?.data) {
            return;
        }

        // PATCH does not return favorited posts to reduce bandwidth
        if(data.favorite_posts) {
            this.account.data.favorite_posts = data.favorite_posts.concat(this.account.data.favorite_posts);
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

        this.account.data = {
            ...this.account.data,
            ...parseAccountDataFromAPI(await resp.json() as Partial<AccountDataResponse>)
        }
    },

    async removeFromAccountData(data: Partial<RemoveAccountDataPayload>): Promise<void> {
        if(!this.account?.data) {
            return;
        }

        // PATCH does not return favorited posts to reduce bandwidth
        if(data.favorite_post_ids) {
            const ids = new Set(data.favorite_post_ids);
            this.account.data.favorite_posts = this.account.data.favorite_posts.filter(p => !ids.has(p.id));
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

        this.account.data = {
            ...this.account.data,
            ...parseAccountDataFromAPI(await resp.json() as Partial<AccountDataResponse>)
        }
    },

    async fetchAccountData()  {
        if(!this.account) {
            return;
        }

        // If the data is already available in the HTML, use it directly
        const preloadedData = JSON.parse(document.getElementById("account-data")!.innerText || "null") as AccountDataResponse;

        if(preloadedData) {
            this.account.data = {
                favorite_posts: [],
                favorite_tags: [],
                blacklist: [],
                search_history: [],
                saved_searches: [],
                ...parseAccountDataFromAPI(preloadedData),
            }
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

            this.account.data = {
                favorite_posts: [],
                favorite_tags: [],
                blacklist: [],
                search_history: [],
                saved_searches: [],
                ...parseAccountDataFromAPI(await resp.json() as AccountDataResponse),
            }
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
        favHeaderOpen: true,
        favoriteTags: [],
        favorites: [],
        fullscreenViewMenuAnchor: "bottomcenter",
        fullscreenViewMenuRotate: false,
        highResImages: true,
        maxPostHeight: 600,
        muteVideo: true,
        newsLastViewedAt: new Date(0),
        queryHistory: [],
        savedSearches: [],
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

            if (k === "queryHistory" || k === "savedSearches") {
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

                if(k === "queryHistory") {
                    val = val.slice(0, QUERY_HISTORY_KEEP_RECENT_LIMIT);
                }
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
    lastFavPage: 1,
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

    async searchPosts(tags: SimpleSerializedSearchQuery, page: number): Promise<PostListResponse> {
        const queryParams = new URLSearchParams();

        queryParams.append("page", page.toString());
        tags.include.forEach(t => {queryParams.append("q", t)});
        tags.exclude.forEach(t => {queryParams.append("q", `-${t}`)});

        const resp = await fetch(`/api/posts?${queryParams.toString()}`);
        let val: any;

        try {
            val = await resp.json();
        } catch(e) {
            console.error("Failed to parse API response", e);
            throw e;
        }

        if (resp.status >= 400) {
            if ("error" in val) {
                throw val.error;
            }
            throw "Something went wrong";
        }

        return val;
    },

    async searchAndUpdateResults({ query, page, force }: { query: SearchQuery, page?: number; force?: boolean; }): Promise<void> {
        this.fetchingPosts = true;
        page = page ?? this.currentPage;
        const sameQuery = query.equals(this.lastQuery);

        // Don't refetch posts we already have
        if (!force && sameQuery && this.posts.has(page)) {
            this.fetchingPosts = false;
            this.currentPage = page;
            return;
        }

        const queryTags = query.toJSONSimple();
        queryTags.exclude = queryTags.exclude.concat(this.blacklist().value.map(t => t.name));

        const doSearch = async () => {
            let posts: PostListResponse;

            try {
                posts = await this.searchPosts(queryTags, page);
            } catch(e) {
                this.toast = {
                    msg: typeof e === "string" ? e : "Something went wrong",
                    type: "error",
                };
                this.fetchingPosts = false;
                return;
            }

            this.hasSearched = true;
            this.fetchingPosts = false;

            if (!sameQuery) {
                this.posts.clear();
            }

            this.posts.set(page, posts.results);
            this.resultsPerPage = posts.count_per_page;
            this.totalPostCount = posts.total_count;
            this.currentPage = page;

            if (this.settings.closeSidebarOnSearch) {
                this.sidebarClosed = true;
            }

            this.addQueryToHistory();
            this.lastQuery = this.query.copy();
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

        // Filter out meta queries that don't have a tag
        const searchFilters = /^(fav|height|id|pool|score|sort|source|updated|user|width):/;
        tags = tags.filter(t => !searchFilters.test(t))
        tags = tags.filter(t => !this.cachedTags.has(t));

        if (tags.length === 0) {
            return;
        }

        const queryParams = new URLSearchParams();

        for (const t of tags) {
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
            if (this.account?.data) {
                return this.account.data.favorite_posts;
            }
            return this.settings.favorites;
        })
    },

    async addFavoritePosts(posts: Post[]) {
        if(this.account?.data) {
            return this.addToAccountData({favorite_posts: posts});
        }
        this.settings.favorites = posts.concat(this.settings.favorites);
        this.saveSettings();
    },

    async removeFavoritePosts(postIds: number[]) {
        if(this.account?.data) {
            return this.removeFromAccountData({favorite_post_ids: postIds});
        }
        this.settings.favorites = this.settings.favorites.filter(p => !postIds.includes(p.id));
        this.saveSettings();
    },

    favoriteTags(): ComputedRef<Tag[]> {
        return computed(() => {
            if (this.account?.data) {
                return this.account.data.favorite_tags;
            }
            return this.settings.favoriteTags;
        });
    },

    async addFavoriteTags(tags: Tag[]) {
        if(this.account?.data) {
            return this.addToAccountData({favorite_tags: tags});
        }
        this.settings.favoriteTags = tags.concat(this.settings.favoriteTags);
        this.saveSettings();
    },

    async removeFavoriteTags(tagNames: string[]) {
        if(this.account?.data) {
            return this.removeFromAccountData({favorite_tag_names: tagNames});
        }
        this.settings.favoriteTags = this.settings.favoriteTags.filter(t => !tagNames.includes(t.name));
        this.saveSettings();
    },

    blacklist(): ComputedRef<Tag[]> {
        return computed(() => {
            if (this.account?.data) {
                return this.account.data.blacklist;
            }
            return this.settings.blacklist;
        });
    },

    async addToBlacklist(tags: Tag[]) {
        if(this.account?.data) {
            return this.addToAccountData({blacklist: tags});
        }
        this.settings.blacklist = tags.concat(this.settings.blacklist);
        this.saveSettings();
    },

    async removeFromBlacklist(tagNames: string[]) {
        if(this.account?.data) {
            return this.removeFromAccountData({blacklist_names: tagNames});
        }
        this.settings.blacklist = this.settings.blacklist.filter(t => !tagNames.includes(t.name));
        this.saveSettings();
    },

    searchHistory(): ComputedRef<SearchHistory[]> {
        return computed(() => {
            if (this.account?.data) {
                return this.account.data.search_history;
            }
            return this.settings.queryHistory;
        });
    },

    async addToSearchHistory(history: SearchHistory[]) {
        if(this.account?.data) {
            const serialized = history.map(h => ({
                date: h.date.toISOString(),
                query: h.query.toJSON()
            }));
            return this.addToAccountData({search_history: serialized});
        }
        this.settings.queryHistory = history.concat(this.settings.queryHistory);
        this.saveSettings();
    },

    async removeFromSearchHistory(queries: SimpleSerializedSearchQuery[]) {
        if(this.account?.data) {
            return this.removeFromAccountData({search_queries: queries});
        }
        this.settings.queryHistory = this.settings.queryHistory.filter(h =>
            !queries.some(q => h.query.equalsSimple(q))
        );
        this.saveSettings();
    },

    savedSearches(): ComputedRef<SavedSearch[]> {
        return computed(() => {
            if (this.account?.data) {
                return this.account.data.saved_searches;
            }
            return this.settings.savedSearches;
        });
    },

    async addToSavedSearches(searches: SavedSearch[]) {
        if(this.account?.data) {
            const serialized = searches.map(h => ({
                date: h.date.toISOString(),
                query: h.query.toJSON()
            }));
            return this.addToAccountData({saved_searches: serialized});
        }
        this.settings.savedSearches = searches.concat(this.settings.savedSearches);
        this.saveSettings();
    },

    async removeFromSavedSearches(queries: SimpleSerializedSearchQuery[]) {
        if(this.account?.data) {
            return this.removeFromAccountData({saved_queries: queries});
        }
        this.settings.savedSearches = this.settings.savedSearches.filter(h =>
            !queries.some(q => h.query.equalsSimple(q))
        );
        this.saveSettings();
    },
});

export default store;
