<script setup lang="ts">
import CodebergLogo from "@/assets/codeberg.svg";
import { changeLog } from "@/changelog";
import Collapsable from "@/components/Collapsable.vue";
import {
    COMMIT_SHA,
    DISCORD_LINK,
    LAST_COMMIT_DATE,
    SURVEY_LINK,
} from "@/config";
import ChangeLog from "./ChangeLog.vue";
</script>

<template>
    <p>
        Booruview is an open source
        <a href="https://gelbooru.com" target="_blank">Gelbooru</a> viewer. It's
        simple, fast, and works on desktop and mobile.
    </p>
    <p>Check the help tab for tips and useful info.</p>
    <ul class="links">
        <li>
            <a href="https://codeberg.org/jessienyan/booruview" target="_blank"
                ><img class="codeberg" :src="CodebergLogo" />Source code</a
            >
        </li>
        <li>
            <a
                href="https://codeberg.org/jessienyan/booruview/issues"
                target="_blank"
                ><i class="bi bi-bug-fill"></i>Report a bug</a
            >
        </li>
        <li>
            <a :href="DISCORD_LINK" target="_blank"
                ><i class="bi bi-discord"></i>Discord (feedback+art)</a
            >
        </li>
        <li>
            <a :href="SURVEY_LINK" target="_blank"
                ><i class="bi bi-send-fill"></i>Give anonymous feedback</a
            >
        </li>
    </ul>

    <p>
        <span class="text-pink"
            >Booruview costs $25/mo to host. Help keep it running with a
            donation. Any amount helps!</span
        >
    </p>

    <p>
        <Collapsable text="donate crypto">
            <p>
                Bitcoin (BTC)
                <code class="wallet">35jCje24EGQLDzW2gkmhdT1D8pxep4nN67</code>
            </p>

            <p>
                Bitcoin Cash (BCH)
                <code class="wallet">19VgrR32NZGhS4pJxnBDWZD293SaCrphbs</code>
            </p>

            <p>
                Ethereum Classic (ETC)
                <code class="wallet"
                    >0xB4cF6eCD1152D68204875B00d8A2D250105f51c1</code
                >
            </p>

            <p>
                Ethereum (ETH)
                <code class="wallet"
                    >0x0f9dF73967fD57f04decb9AF9c4aa04342c153F7</code
                >
            </p>

            <p>
                Litecoin (LTC)
                <code class="wallet">MLGX8DiuWKuSKoLXdzRtw1QHxZ6Hu6Veqr</code>
            </p>

            <p>
                Tether (USDT)
                <code class="wallet"
                    >0xc5614040d18B1fC8e8E77e6599553F5EEC946C69</code
                >
            </p>

            <p>
                XRP
                <code class="wallet">rw2ciyaNshpHe7bCHo4bRWq6pqqynnWKQg</code>
                <code class="wallet">memo: 783456807</code>
            </p>
        </Collapsable>
    </p>

    <p>
        <Collapsable text="changelog">
            <ChangeLog
                v-for="change in changeLog"
                :date="change.date"
                :changes="change.changes"
            />

            <p class="version">
                latest commit
                <a
                    :href="`https://codeberg.org/jessienyan/booruview/commit/${COMMIT_SHA}`"
                    target="_blank"
                    ><code>{{ COMMIT_SHA }}</code></a
                >
                @ {{ LAST_COMMIT_DATE }}
            </p>
        </Collapsable>
    </p>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";

.version {
    font-family: "Courier New", Courier, monospace;
    font-size: 0.75rem;
    text-align: right;
}

.links .bi,
.codeberg {
    margin-right: 0.4rem;
    vertical-align: middle;
}

.wallet {
    display: block;
    width: fit-content;
    font-family: "Courier New", Courier, monospace;
    background-color: rgba(121, 121, 121, 0.2);
    padding: 0 3px;
    border: 1px solid #4c4c4c;
    border-radius: 2px;
    margin-top: 4px;
}
</style>
