<script setup lang="ts">
import { ref } from "vue";
import store from "@/store";

const username = ref("");
const password = ref("");

async function onSubmit() {
    try {
        const resp = await fetch("/api/login", {
            method: "POST",
            body: JSON.stringify({
                username: username.value,
                password: password.value,
            }),
            headers: { "Content-Type": "application/json" },
        });

        const data = await resp.json();
        if (data.error) {
            store.toast = {
                msg: data.error,
                type: "error",
            };
            return;
        }

        store.account = {
            authToken: data.auth_token,
            username: data.username,
        };
        store.saveAccount();
        store.toast = {
            msg: "Logged in successfully",
            type: "info",
        };
    } catch (e) {
        console.error(e);
        store.toast = {
            msg: "Something went wrong :(",
            type: "error",
        };
    }
}
</script>

<template>
    <form @submit.prevent="onSubmit">
        <input
            v-model="username"
            class="text-input rounded"
            type="text"
            placeholder="username"
            maxlength="16"
            minLength="3"
            required
        />
        <input
            v-model="password"
            class="text-input rounded"
            type="password"
            placeholder="password"
            required
        />

        <button class="submit btn-primary btn-rounded" type="submit">
            login
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
    margin-bottom: 1em;
}

.submit {
    width: 100%;
    display: block;
}
</style>
