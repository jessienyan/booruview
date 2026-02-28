<script setup lang="ts">
import { ref, useTemplateRef, watch } from "vue";
import store, {
    type ColumnSizing,
    type FullscreenViewMenuAnchorPoint,
} from "@/store";
import Blacklist from "./Blacklist.vue";

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

function onChangeColCount(e: Event) {
    store.settings.columnCount = parseInt(
        (e.target as HTMLInputElement).value,
        10,
    );
    store.saveSettings();
}

function onChangeColSizing(e: Event) {
    store.settings.columnSizing = (e.target as HTMLInputElement)
        .value as ColumnSizing;
    store.saveSettings();
}

function onChangeColWidth(e: Event) {
    store.settings.columnWidth = parseInt(
        (e.target as HTMLInputElement).value,
        10,
    );
    store.saveSettings();
}

const MAX_POST_HEIGHT = 1500;

function onChangePostHeight(e: Event) {
    const $el = e.target as HTMLInputElement;

    if ($el.value === $el.max) {
        store.settings.maxPostHeight = null;
    } else {
        store.settings.maxPostHeight = parseInt($el.value, 10);
    }

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

function onChangeCheckForUpdates(e: Event) {
    store.settings.checkForUpdates = (e.target as HTMLInputElement).checked;
    store.saveSettings();
}

const exportCode = ref("");
const exportCodeRef = useTemplateRef("export-code");
const canGenerate = ref(true);
const generatingCode = ref(false);

// Re-enable code generation once a change is made
watch(
    () => store.settings,
    () => {
        canGenerate.value = true;
    },
);

const ignoredSettingsForImportExport: Array<
    Partial<keyof typeof store.settings>
> = ["consented", "queryHistory"];

function generateExportCode() {
    generatingCode.value = true;

    const data: Partial<typeof store.settings> = Object.assign(
        {},
        store.settings,
    );

    for (const k of ignoredSettingsForImportExport) {
        delete data[k];
    }

    fetch("/api/settings/export", {
        method: "POST",
        body: JSON.stringify(data),
    })
        .then((resp) => {
            resp.json().then(({ code }) => {
                exportCode.value = code;
            });
            canGenerate.value = false;
        })
        .catch((e) => console.error(e))
        .finally(() => {
            generatingCode.value = false;
        });
}

const importCode = ref("");
const importingData = ref(false);

function mergeSettings(data: Record<string, any>) {
    const blacklist = store.blacklist().value;
    const favPosts = store.favoritePosts().value;
    const favTags = store.favoriteTags().value;

    Object.entries(data).forEach(([_k, _v]) => {
        const k = _k as Exclude<
            keyof typeof store.settings,
            typeof ignoredSettingsForImportExport
        >;
        let v = _v;

        // Remove any duplicates
        if (k === "blacklist") {
            v = (v as Tag[])
                .filter(
                    (a) => blacklist.findIndex((b) => a.name === b.name) === -1,
                )
                .concat(blacklist);
        } else if (k === "favorites") {
            v = (v as Post[])
                .filter((a) => favPosts.findIndex((b) => a.id === b.id) === -1)
                .concat(favPosts);
        } else if (k === "favoriteTags") {
            v = (v as Tag[])
                .filter(
                    (a) => favTags.findIndex((b) => a.name === b.name) === -1,
                )
                .concat(favTags);
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
        .then((resp) => {
            if (resp.status >= 400) {
                resp.json()
                    .then((val) => {
                        let msg = "Something went wrong";

                        if ("error" in val) {
                            msg = val.error;
                        }

                        store.toast = {
                            msg,
                            type: "error",
                        };
                    })
                    .catch(() => {
                        store.toast = {
                            msg: "Something went wrong",
                            type: "error",
                        };
                    });
                return;
            }

            resp.json().then((data) => mergeSettings(data));
            store.toast = {
                msg: "data imported successfully",
                type: "info",
            };
            importCode.value = "";
        })
        .catch((e) => console.error(e))
        .finally(() => {
            importingData.value = false;
        });
}
</script>

<template>
    <div class="settings-container">
        <h2>app</h2>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.checkForUpdates"
                    @change="onChangeCheckForUpdates"
                />
                periodically check for updates</label
            >
        </div>

        <h2>blacklist</h2>

        <Blacklist />

        <h2>content</h2>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.highResImages"
                    @change="onChangeHighResImages"
                    disabled
                />
                high resolution images (uncheck if slow connection)</label
            >
            <p class="text-pink">
                hi-res images are temporarily disabled while we improve load
                times
            </p>
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

        <h2>layout</h2>

        <div class="input-group">
            <label># of columns</label>
            <div class="input-container">
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
            <div class="input-container">
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
            <div class="input-container">
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
            <label>max post height </label>
            <div class="input-container">
                <input
                    type="range"
                    min="100"
                    :max="MAX_POST_HEIGHT"
                    step="50"
                    :value="store.settings.maxPostHeight ?? MAX_POST_HEIGHT"
                    @input="onChangePostHeight"
                />
                <span class="value">{{
                    store.settings.maxPostHeight != null
                        ? `${store.settings.maxPostHeight}px`
                        : "none"
                }}</span>
            </div>
        </div>

        <div class="input-group">
            <label>fullscreen view menu position</label>
            <div class="input-container">
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
                vertical controls/menu</label
            >
        </div>

        <h2>search</h2>

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

        <h2>import/export</h2>

        <p>
            Copy your data to another device (settings/blacklist/favorites).
            Your blacklist and favorites will be merged together.
        </p>

        <p>
            First, click export to generate a code from your device. Then on
            another device, enter the code and click import.
        </p>

        <div class="input-group">
            <label
                >1. Export your data and generate a code. The code will expire
                after 15 minutes.</label
            >
            <div class="input-container text-btn-combo">
                <input
                    class="text-input rounded-start"
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
            <label
                >2. Import your data onto another device using the code
                generated above.</label
            >
            <div class="input-container text-btn-combo">
                <input
                    class="text-input rounded-start"
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
@import "@/assets/form";

.input-container {
    display: flex;

    input,
    select {
        flex: 1;
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

label + .input-container {
    margin-top: 0.4rem;
}

label {
    display: block;
}

.text-btn-combo {
    .text-input {
        border-right: 0;
        text-align: center;
    }

    button {
        border-top-right-radius: 5px;
        border-bottom-right-radius: 5px;
        padding-left: 1rem;
        padding-right: 1rem;
    }
}
</style>
