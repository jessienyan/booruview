<script setup lang="ts">
import store, {
    type ColumnSizing,
    type FullscreenViewMenuAnchorPoint,
    type PageLoadAutoSearch,
} from "@/store";
import { ref, useTemplateRef, watch } from "vue";

const exportCodeRef = useTemplateRef("export-code");

const columnSizingOptions: Record<ColumnSizing, string> = {
    dynamic: "dynamic",
    fixed: "fixed",
};

const fullscreenViewMenuAnchorOptions: Record<
    FullscreenViewMenuAnchorPoint,
    string
> = {
    topleft: "top left",
    topcenter: "top center",
    topright: "top right",
    right: "right",
    bottomright: "bottom right",
    bottomcenter: "bottom center",
    bottomleft: "bottom left",
    left: "left",
};

const searchOnPageLoadOptions: Record<PageLoadAutoSearch, string> = {
    always: "always",
    "if-query": "if search isn't empty",
    never: "never",
};

function onChangeColCount(e: Event) {
    store.settings.columnCount = parseInt((e.target as HTMLInputElement).value);
    store.saveSettings();
}

function onChangeColSizing(e: Event) {
    store.settings.columnSizing = (e.target as HTMLInputElement)
        .value as ColumnSizing;
    store.saveSettings();
}

function onChangeColWidth(e: Event) {
    store.settings.columnWidth = parseInt((e.target as HTMLInputElement).value);
    store.saveSettings();
}

function onChangeFullscreenViewMenuAnchor(e: Event) {
    store.settings.fullscreenViewMenuAnchor = (e.target as HTMLInputElement)
        .value as FullscreenViewMenuAnchorPoint;
    store.saveSettings();
}

function onChangeFullscreenViewMenuRotate(e: Event) {
    store.settings.fullscreenViewMenuRotate = (
        e.target as HTMLInputElement
    ).checked;
    store.saveSettings();
}

function onChangeCloseSidebarOnSearch(e: Event) {
    store.settings.closeSidebarOnSearch = (
        e.target as HTMLInputElement
    ).checked;
    store.saveSettings();
}

function onChangeSearchOnPageLoad(e: Event) {
    store.settings.searchOnPageLoad = (e.target as HTMLInputElement)
        .value as any;
    store.saveSettings();
}

function onChangeHighResImages(e: Event) {
    store.settings.highResImages = (e.target as HTMLInputElement).checked;
    store.saveSettings();
}

function onChangeAutoplayVideos(e: Event) {
    store.settings.autoplayVideo = (e.target as HTMLInputElement).checked;
    store.saveSettings();
}

function onChangeMuteVideos(e: Event) {
    store.settings.muteVideo = (e.target as HTMLInputElement).checked;
    store.saveSettings();
}

const exportCode = ref("");
const canGenerate = ref(true);
const generatingCode = ref(false);

// Re-enable code generation once a change is made
watch(
    () => store.settings,
    () => (canGenerate.value = true),
);

function generateExportCode() {
    generatingCode.value = true;

    const data: Partial<typeof store.settings> = Object.assign(
        {},
        store.settings,
    );
    delete data["queryHistory"];

    fetch("/api/settings/export", {
        method: "POST",
        body: JSON.stringify(data),
    })
        .then((resp) => {
            resp.json().then(({ code }) => (exportCode.value = code));
            canGenerate.value = false;
        })
        .catch((e) => console.error(e))
        .finally(() => (generatingCode.value = false));
}

const importCode = ref("");
const importingData = ref(false);

function mergeSettings(data: Record<string, any>) {
    Object.entries(data).forEach(([_k, v]) => {
        const k = _k as Exclude<keyof typeof store.settings, "queryHistory">;

        // Remove any duplicates
        if (k === "blacklist") {
            // prettier-ignore
            v = (v as Tag[]).filter((a) => store.settings.blacklist.findIndex((b) => a.name === b.name) === -1);
            v = v.concat(store.settings.blacklist);
        } else if (k === "favorites") {
            // prettier-ignore
            v = (v as Post[]).filter((a) => store.settings.favorites.findIndex((b) => a.id === b.id) === -1);
            v = v.concat(store.settings.favorites);
        }

        (store.settings as any)[k] = v;
    });

    store.saveSettings();
}

function importData() {
    importingData.value = true;

    fetch("/api/settings/import", {
        method: "POST",
        body: JSON.stringify({ code: importCode.value }),
    })
        .then((resp) => resp.json().then((data) => mergeSettings(data)))
        .catch((e) => console.error(e))
        .finally(() => (importingData.value = false));
}
</script>

<template>
    <div class="settings-container">
        <h3>layout</h3>

        <div class="input-group">
            <label># of columns</label>
            <div class="input">
                <select
                    :value="store.settings.columnSizing"
                    @change="onChangeColSizing"
                >
                    <option
                        v-for="(label, val) in columnSizingOptions"
                        :value="val"
                    >
                        {{ label }}
                    </option>
                </select>
            </div>
        </div>

        <div class="input-group" v-if="store.settings.columnSizing === 'fixed'">
            <label>column count</label>
            <div class="input">
                <input
                    type="range"
                    min="1"
                    max="20"
                    step="1"
                    :value="store.settings.columnCount"
                    @input="onChangeColCount"
                />
                <span class="value">{{ store.settings.columnCount }}</span>
            </div>
        </div>

        <div
            class="input-group"
            v-if="store.settings.columnSizing === 'dynamic'"
        >
            <label>max column width</label>
            <div class="input">
                <input
                    type="range"
                    min="100"
                    max="1000"
                    step="10"
                    :value="store.settings.columnWidth"
                    @input="onChangeColWidth"
                />
                <span class="value">{{ store.settings.columnWidth }}px</span>
            </div>
        </div>

        <div class="input-group">
            <label>fullscreen view menu position</label>
            <div class="input">
                <select
                    :value="store.settings.fullscreenViewMenuAnchor"
                    @change="onChangeFullscreenViewMenuAnchor"
                >
                    <option
                        v-for="(label, val) in fullscreenViewMenuAnchorOptions"
                        :value="val"
                    >
                        {{ label }}
                    </option>
                </select>
            </div>
        </div>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.fullscreenViewMenuRotate"
                    @change="onChangeFullscreenViewMenuRotate"
                />
                rotate menu 90 degrees</label
            >
        </div>

        <h3>search</h3>

        <div class="input-group">
            <label>auto-search when page loads</label>
            <div class="input">
                <select
                    :value="store.settings.searchOnPageLoad"
                    @change="onChangeSearchOnPageLoad"
                >
                    <option
                        v-for="(label, val) in searchOnPageLoadOptions"
                        :value="val"
                    >
                        {{ label }}
                    </option>
                </select>
            </div>
        </div>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.closeSidebarOnSearch"
                    @change="onChangeCloseSidebarOnSearch"
                />
                searching closes sidebar</label
            >
        </div>

        <h3>content</h3>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.highResImages"
                    @change="onChangeHighResImages"
                />
                high resolution images (uncheck if slow connection)</label
            >
        </div>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.autoplayVideo"
                    @change="onChangeAutoplayVideos"
                />
                autoplay videos</label
            >
        </div>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.muteVideo"
                    @change="onChangeMuteVideos"
                />
                mute videos</label
            >
        </div>

        <h3>import/export</h3>

        <p>
            <i class="bi bi-info-circle"></i> Copy your data to another device
            (settings, blacklist, favorites). Importing will combine/merge your
            blacklist and favorites.
        </p>

        <div class="input-group">
            <label
                >1. Export data and generate a code. The code will expire after
                15 minutes.</label
            >
            <div class="input text-btn-combo">
                <input
                    ref="export-code"
                    type="text"
                    placeholder="export to generate code"
                    readonly
                    :value="exportCode"
                    @focus="exportCodeRef?.setSelectionRange(0, 999)"
                />
                <button
                    class="btn-primary"
                    @click="generateExportCode"
                    :disabled="generatingCode || !canGenerate"
                >
                    export
                </button>
            </div>
        </div>

        <div class="input-group">
            <label>2. Import using the generated code.</label>
            <div class="input text-btn-combo">
                <input
                    type="text"
                    placeholder="xxxx-xxxx-xxxx"
                    v-model="importCode"
                />
                <button
                    class="btn-primary"
                    @click="importData"
                    :disabled="
                        importingData || !/^\d{4}-\d{4}-\d{4}$/.test(importCode)
                    "
                >
                    import
                </button>
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.input {
    display: flex;

    input,
    select {
        flex: 1;
    }

    input[type="text"] {
        background-color: #252525;
        border: 1px solid #555;
    }

    .value {
        min-width: 60px;
        text-align: center;
    }
}

p,
.input-group {
    margin: 1rem 0 1rem 1rem;
}

label + .input {
    margin-top: 0.4rem;
}

label {
    display: block;
}

.text-btn-combo {
    input[type="text"] {
        border-top-left-radius: 5px;
        border-bottom-left-radius: 5px;
        border-right: 0;
        text-align: center;
        color: #aaa;
    }

    button {
        border-top-right-radius: 5px;
        border-bottom-right-radius: 5px;
        padding-left: 1rem;
        padding-right: 1rem;
    }
}
</style>
