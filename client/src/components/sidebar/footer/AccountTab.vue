<script setup lang="ts">
import { ref } from "vue";
import store from "@/store";
import ChangePassword from "./account/ChangePassword.vue";
import DeleteAccount from "./account/DeleteAccount.vue";
import FAQ from "./account/FAQ.vue";
import Login from "./account/Login.vue";
import Register from "./account/Register.vue";

const showRegisterForm = ref(false);
const showChangePasswordForm = ref(false);

function logout() {
    store.account = null;
    store.saveAccountCredentials();
}

function downloadData() {
    // Shouldn't happen
    if (!store.account) {
        return;
    }

    const localFavPosts = new Set(store.settings.favorites.map((p) => p.id));
    const addToFavPosts = store.account.data.favorite_posts.filter(
        (p) => !localFavPosts.has(p.id),
    );
    store.settings.favorites = addToFavPosts.concat(store.settings.favorites);

    const localFavTags = new Set(
        store.settings.favoriteTags.map((t) => t.name),
    );
    const addToFavTags = store.account.data.favorite_tags.filter(
        (t) => !localFavTags.has(t.name),
    );
    store.settings.favoriteTags = addToFavTags.concat(
        store.settings.favoriteTags,
    );

    const localBlacklist = new Set(store.settings.blacklist.map((t) => t.name));
    const addToBlacklist = store.account.data.blacklist.filter(
        (t) => !localBlacklist.has(t.name),
    );
    store.settings.blacklist = addToBlacklist.concat(store.settings.blacklist);

    store.saveSettings();
    store.toast = {
        msg: "Download OK!",
        type: "info",
    };
}

const isUploading = ref(false);
async function uploadData() {
    // Shouldn't happen
    if (!store.account) {
        return;
    }

    isUploading.value = true;

    const accountFavPosts = new Set(
        store.account.data.favorite_posts.map((p) => p.id),
    );
    const addToFavPosts = store.settings.favorites.filter(
        (p) => !accountFavPosts.has(p.id),
    );
    store.account.data.favorite_posts = addToFavPosts.concat(
        store.account.data.favorite_posts,
    );

    const accountFavTags = new Set(
        store.account.data.favorite_tags.map((t) => t.name),
    );
    const addToFavTags = store.settings.favoriteTags.filter(
        (t) => !accountFavTags.has(t.name),
    );
    store.account.data.favorite_tags = addToFavTags.concat(
        store.account.data.favorite_tags,
    );

    const accountBlacklist = new Set(
        store.account.data.blacklist.map((t) => t.name),
    );
    const addToBlacklist = store.settings.blacklist.filter(
        (t) => !accountBlacklist.has(t.name),
    );
    store.account.data.blacklist = addToBlacklist.concat(
        store.account.data.blacklist,
    );

    try {
        await store.saveAccountData({
            favorite_posts: true,
            favorite_tags: true,
            blacklist: true,
        });
        store.toast = {
            msg: "Upload OK!",
            type: "info",
        };
    } finally {
        isUploading.value = false;
    }
}
</script>

<template>
    <template v-if="!store.account">
        <template v-if="showRegisterForm">
            <h3>register</h3>
            <Register />
            <p>
                <a
                    class="change-form"
                    href="#"
                    @click.prevent="showRegisterForm = false"
                    >I already have an account</a
                >
            </p>
        </template>
        <template v-else>
            <h3>login</h3>
            <Login />
            <p>
                <a
                    class="change-form"
                    href="#"
                    @click.prevent="showRegisterForm = true"
                    >Create an account</a
                >
            </p>
        </template>
        <FAQ />
    </template>
    <template v-else>
        <p>
            <strong>You are signed in as {{ store.account.username }}.</strong>
        </p>

        <button
            class="btn-primary btn-rounded btn-block"
            @click="uploadData"
            :disabled="isUploading"
        >
            Upload Data <i class="bi bi-cloud-upload"></i>
        </button>
        <p class="hint">
            Upload data from your device and add it to your account.
        </p>

        <button class="btn-primary btn-rounded btn-block" @click="downloadData">
            Download Data <i class="bi bi-cloud-download"></i>
        </button>
        <p class="hint">
            Download data from your account and add it to your device.
        </p>

        <button class="btn-primary btn-rounded btn-block" @click="logout">
            Logout
        </button>
        <ChangePassword />
        <DeleteAccount />
    </template>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";
@import "@/assets/form";

.text-input {
    text-align: center;
}

.change-form {
    display: block;
    text-align: center;
}

.hint {
    color: #888;
}
</style>
