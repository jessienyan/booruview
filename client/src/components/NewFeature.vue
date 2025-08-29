<script setup lang="ts">
const strings = [
    "owo",
    "uwu",
    "OwO",
    "UwU",
    "◕⩊◕",
    "( ‘• ω • `)",
    "( • ̀ω•́ )✧",
    "( •̀ ω •́ )",
    "(*＞ω＜*)",
    "(,,>﹏<,,)",
    ">⩊<",
    "\(^ω^\ )",
    "ඞ",
    "( > 〰 < )",
    "(๑ > ᴗ < ๑)",
    "₍^.  ̫.^₎",
];
const randomString = strings[Math.round(Math.random() * (strings.length - 1))];
</script>

<template>
    <span class="text-pink new-feature">{{ randomString }}</span>
</template>

<style scoped>
.new-feature {
    font-weight: bold;
    text-wrap: nowrap;
}
</style>
