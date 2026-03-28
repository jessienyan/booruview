<script setup lang="ts">
import { ref } from "vue";
import store from "@/store";

const username = ref<string | null>(null);
const password = ref<string | null>(null);
const passwordConfirm = ref<string | null>(null);
const submitting = ref(false);

function toastError(msg: string) {
    store.toast = {
        msg,
        type: "error",
    };
}

async function onSubmit(e: Event) {
    if (submitting.value) {
        return;
    }

    const $form = e.target as HTMLFormElement;

    if (!$form.reportValidity()) {
        return;
    } else if (password.value !== passwordConfirm.value) {
        toastError("Password don't match");
        return;
    }

    try {
        submitting.value = true;

        const resp = await fetch("/api/register", {
            method: "POST",
            body: JSON.stringify({
                username: username.value,
                password: password.value,
            }),
            headers: {
                "Content-Type": "application/json",
            },
        });

        submitting.value = false;

        const data = await resp.json();
        if (data.error) {
            toastError(data.error);
            return;
        } else if (!resp.ok) {
            toastError(
                `Unexpected error (${resp.status}): ${await resp.text()}`,
            );
            return;
        }

        store.account = {
            username: username.value!,
            authToken: data.auth_token,
            data: {
                favorite_posts: [...store.settings.favorites],
                favorite_tags: [...store.settings.favoriteTags],
                blacklist: [...store.settings.blacklist],
                search_history: [...store.settings.queryHistory],
            },
        };
        store.saveAccountCredentials();
    } catch (e) {
        console.error(e);
        toastError("Something went wrong :(");
        submitting.value = false;
        return;
    } finally {
        submitting.value = false;
    }

    store.toast = {
        msg: "You are now logged in.",
        type: "info",
    };

    // Upload existing data
    store.addToAccountData({
        blacklist: store.settings.blacklist,
        favorite_posts: store.settings.favorites,
        favorite_tags: store.settings.favoriteTags,
        search_history: store.settings.queryHistory.map((q) => ({
            date: q.date.toISOString(),
            query: q.query.toJSON(),
        })),
    });
}
</script>

<template>
    <form @submit.prevent="onSubmit">
        <input
            v-model="username"
            class="text-input rounded"
            :class="{ touched: username != null }"
            type="text"
            placeholder="username"
            pattern="[a-zA-Z0-9_\-]+"
            maxlength="16"
            minlength="3"
            required
        />
        <span class="tip"
            >Between 3 and 16 characters. Your username is private and only for
            logging in.</span
        >

        <input
            v-model="password"
            class="text-input rounded"
            :class="{ touched: password != null }"
            type="password"
            placeholder="password"
            minlength="8"
            required
        />
        <span class="tip"
            >At least 8 characters. YOU CANNOT RESET YOUR PASSWORD. Use
            something that's easy to remember!</span
        >

        <input
            v-model="passwordConfirm"
            class="text-input rounded"
            :class="{ touched: passwordConfirm != null }"
            type="password"
            placeholder="confirm password"
            :minlength="8"
            required
        />

        <button
            class="submit btn-primary btn-rounded"
            type="submit"
            @click="onSubmit"
            :disabled="submitting"
        >
            register
        </button>
    </form>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/form";

form {
    padding: 0 1em;
}

.text-input {
    width: 100%;
}

.tip {
    display: block;
    font-size: 14px;
    font-style: italic;
    margin: 0.5em 0.5em 1.5em;
}

.submit {
    margin-top: 1.5em;
    width: 100%;
    display: block;
}
</style>
