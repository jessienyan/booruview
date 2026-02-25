<script lang="ts" setup>
import { computed, ref, watch } from "vue";
import store from "@/store";

const showDeleteConfirm = ref(false);
const usernameConfirm = ref("");

// Reset form when it's closed
watch([showDeleteConfirm], () => {
    if (!showDeleteConfirm.value) {
        usernameConfirm.value = "";
    }
});

const canDelete = computed(
    () => store.account && store.account.username === usernameConfirm.value,
);

async function doDelete() {
    if (!canDelete.value) {
        return;
    }

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
            <h3>Deleting your account is PERMANENT</h3>
            <p>
                All your data and account info will be erased. This cannot be
                undone.
                <strong>Backup your data now or it's gone forever.</strong>
            </p>
            <p>
                To continue, enter your username:
                <code>{{ store.account.username }}</code>
                <input
                    v-model="usernameConfirm"
                    class="text-input input-block rounded"
                    type="text"
                    placeholder="confirm username"
                />
            </p>
            <p>
                <button
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

    p {
        margin: 1.5em 0;

        &:last-child {
            margin-bottom: inherit;
        }
    }
}

h3 {
    text-align: center;
    margin: 1em 0;
    color: #c00;
}
</style>
