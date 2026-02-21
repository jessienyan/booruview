<script lang="ts" setup>
import { computed, ref } from "vue";
import store from "@/store";

const showDeleteConfirm = ref(false);
const usernameConfirm = ref("");

const canDelete = computed(
    () => store.account && store.account.username === usernameConfirm.value,
);

async function doDelete() {
    try {
        const resp = await fetch("/api/account", {
            method: "DELETE",
            body: JSON.stringify({ permanently_delete_account: true }),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (resp.ok) {
            store.account = null;
            store.saveAccountCredentials();
            store.toast = {
                msg: "Account deleted successfully",
                type: "info",
            };
            return;
        }

        const data = await resp.json();
        if (data.error) {
            store.toast = {
                msg: data.error,
                type: "error",
            };
        }
    } catch (e) {
        console.log(e);
        store.toast = {
            msg: "Something went wrong :(",
            type: "error",
        };
    }
}
</script>

<template>
    <template v-if="store.account">
        <button
            class="btn-danger btn-rounded btn-block"
            @click="showDeleteConfirm = true"
            :disabled="showDeleteConfirm"
        >
            Delete Account
        </button>
        <div v-if="showDeleteConfirm" class="confirm-delete">
            <p>Deleting your account is PERMANENT.</p>
            <p>
                You'll keep your current favorites and settings, but you won't
                be able to login and access them on other devices.
            </p>
            <p>
                To continue, enter your username:
                {{ store.account.username }}
            </p>
            <p>
                <input
                    v-model="usernameConfirm"
                    class="text-input input-block rounded"
                    type="text"
                    placeholder="confirm username"
                /><button
                    class="btn-danger btn-rounded btn-block"
                    @click="doDelete"
                    :disabled="!canDelete"
                >
                    I understand, delete it</button
                ><button
                    class="btn-gray btn-rounded btn-block"
                    @click="showDeleteConfirm = false"
                >
                    Nevermind
                </button>
            </p>
        </div>
    </template>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";
@import "@/assets/form";

.confirm-delete {
    background-color: $color-darkgray;
    padding: 0.1px 1em;
}
</style>
