<script setup lang="ts">
import { computed, ref } from "vue";
import Collapsable from "@/components/Collapsable.vue";
import store from "@/store";
import Login from "./Login.vue";
import Register from "./Register.vue";

const showRegisterForm = ref(false);
const showDeleteConfirm = ref(false);
const usernameConfirm = ref("");

const canDelete = computed(
    () => store.account && store.account.username === usernameConfirm.value,
);

async function doDelete() {
    store.toast = {
        msg: "TODO",
        type: "error",
    };
    return;

    // try {
    // 	const resp = await fetch("/api/account", {
    // 		method: "DELETE",
    // 	});
    // } catch (e) {
    // 	console.log(e);
    // }
}

function logout() {
    store.account = null;
    store.saveAccount();
    store.toast = {
        msg: "You have been logged out",
        type: "info",
    };
}
</script>

<template>
    <div>
        <div class="faq">
            <Collapsable text="Account FAQ">
                <Collapsable
                    text="Do I have to register an account?"
                    :button="false"
                >
                    <div class="collapsable">
                        <p>No, accounts are entirely optional.</p>
                    </div>
                </Collapsable>

                <Collapsable text="What data is stored?" :button="false">
                    <div class="collapsable">
                        <p>
                            Accounts store your favorite tags and posts,
                            blacklist, and search history on Booruview.
                        </p>
                    </div>
                </Collapsable>
                <Collapsable text="How do accounts work?" :button="false">
                    <div class="collapsable">
                        <p>
                            It lets you access your data from multiple devices
                            and keeps those devices in sync.
                        </p>
                        <p>
                            Without an account, your data is only saved in your
                            browser. You're at risk of losing your data if your
                            device breaks or you accidentally clear your browser
                            data.
                        </p>
                    </div>
                </Collapsable>

                <Collapsable text="How is my data synced?" :button="false">
                    <div class="collapsable">
                        <p>
                            While logged in, your data is loaded when you open
                            Booruview or refresh the page. Any changes you make
                            are sent back to Booruview.
                        </p>
                    </div>
                </Collapsable>
            </Collapsable>
        </div>

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
        </template>
        <template v-else>
            <p v-if="store.account">
                <strong
                    >You are signed in as {{ store.account.username }}.</strong
                >
            </p>
            <p>
                <button
                    class="btn-primary btn-rounded btn-block"
                    @click="logout"
                >
                    Logout
                </button>
                <button
                    class="btn-danger btn-rounded btn-block"
                    @click="showDeleteConfirm = true"
                    :disabled="showDeleteConfirm"
                >
                    Delete Account
                </button>
            </p>
            <div v-if="showDeleteConfirm" class="confirm-delete">
                <p>Deleting your account is PERMANENT.</p>
                <p>
                    You'll keep your current favorites and settings, but you
                    won't be able to login and access them on other devices.
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
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";
@import "@/assets/form";

.faq {
    margin-bottom: 1em;
}

.collapsable {
    background-color: $color-primary-darker;
    color: $color-primary-light;
    padding: 0.1px 1em;
}

.confirm-delete {
    background-color: $color-darkgray;
    padding: 0.1px 1em;
}

.text-input {
    text-align: center;
}

.change-form {
    display: block;
    text-align: center;
}
</style>
