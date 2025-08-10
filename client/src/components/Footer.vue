<script setup lang="ts">
import store from "@/store";
import { computed, ref } from "vue";

const clickedWhich = ref<"prev" | "next" | null>(null);
const fmt = new Intl.NumberFormat();
const currentPageText = computed(() => fmt.format(store.currentPage));
const maxPageText = computed(() => fmt.format(store.maxPage()));
const totalPostCountText = computed(() => fmt.format(store.totalPostCount));

function nextPage() {
    clickedWhich.value = "next";
    store.nextPage()?.finally(() => (clickedWhich.value = null));
}

function prevPage() {
    clickedWhich.value = "prev";
    store.prevPage()?.finally(() => (clickedWhich.value = null));
}
</script>

<template>
    <footer class="page-nav">
        <div class="nav-btns">
            <button
                class="btn-primary btn-rounded"
                :disabled="store.fetchingPosts"
                @click="prevPage"
                v-if="store.currentPage > 1"
            >
                <div
                    class="spinner"
                    v-if="store.fetchingPosts && clickedWhich === 'prev'"
                >
                    <span class="spinner-inner"></span>
                </div>
                <template v-else>
                    <i class="bi bi-arrow-left"></i> prev
                </template>
            </button>
            <button
                class="btn-primary btn-rounded"
                :disabled="store.fetchingPosts"
                @click="nextPage"
                v-if="store.currentPage < store.maxPage()"
            >
                <div
                    class="spinner"
                    v-if="store.fetchingPosts && clickedWhich === 'next'"
                >
                    <span class="spinner-inner"></span>
                </div>
                <template v-else>
                    next <i class="bi bi-arrow-right"></i>
                </template>
            </button>
        </div>
        <p>
            page {{ currentPageText }} of {{ maxPageText }} ({{
                totalPostCountText
            }}
            results)
        </p>
    </footer>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.page-nav {
    margin-top: 40px;
    text-align: center;
}

.nav-btns {
    display: flex;
    justify-content: center;
    gap: 20px;

    button {
        position: relative;
        width: 80px;
        height: 40px;
    }
}

.spinner {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);

    width: 20px;
    height: 20px;

    .spinner-inner {
        width: 100%;
        height: 100%;

        border-radius: 50%;
        border: 2px solid #fff;
        border-bottom-color: transparent;
        animation: linear 1s spin-anim infinite;
        display: block;

        @keyframes spin-anim {
            from {
                transform: rotate(0);
            }
            to {
                transform: rotate(360deg);
            }
        }
    }
}
</style>
