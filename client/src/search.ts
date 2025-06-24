class SearchQuery {
    _include: Map<string, Tag>;
    _exclude: Map<string, Tag>;

    constructor() {
        this._include = new Map();
        this._exclude = new Map();
    }

    includedList(): Tag[] {
        return Array.from(this._include.values());
    }

    excludedList(): Tag[] {
        return Array.from(this._exclude.values());
    }

    clear() {
        this._include.clear();
        this._exclude.clear();
    }

    isExcluded(name: string) {
        return this._exclude.has(name);
    }

    isIncluded(name: string) {
        return this._include.has(name);
    }

    isEmpty() {
        return this._include.size + this._exclude.size === 0;
    }

    includeTag(t: Tag) {
        this._include.set(t.name, t);
        this._exclude.delete(t.name);
    }

    excludeTag(t: Tag) {
        this._exclude.set(t.name, t);
        this._include.delete(t.name);
    }

    removeTag(t: Tag) {
        this._include.delete(t.name);
        this._exclude.delete(t.name);
    }

    asList(): string[] {
        const include = Array.from(this._include.values(), (t) => t.name);
        const exclude = Array.from(this._exclude.values(), (t) => "-" + t.name);
        return include.concat(exclude);
    }

    equals(o: SearchQuery): boolean {
        if (
            this._include.size !== o._include.size ||
            this._exclude.size !== o._exclude.size
        ) {
            return false;
        }

        for (const k of this._include.keys()) {
            if (!o._include.has(k)) {
                return false;
            }
        }

        for (const k of this._exclude.keys()) {
            if (!o._exclude.has(k)) {
                return false;
            }
        }

        return true;
    }

    copy(): SearchQuery {
        const clone = new SearchQuery();
        clone._include = new Map(this._include);
        clone._exclude = new Map(this._exclude);
        return clone;
    }
}

export default SearchQuery;
