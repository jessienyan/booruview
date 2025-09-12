<script setup lang="ts">
import { useDismiss, useViewportSize } from "@/composable";
import { computed, toValue, useTemplateRef, type MaybeRef } from "vue";

const container = useTemplateRef("container");
const props = defineProps<{ el: MaybeRef<HTMLElement | null> }>();
const show = defineModel("show", { default: false });

const viewport = useViewportSize();
useDismiss([container, props.el], () => (show.value = false));

const menuPosition = computed(() => {
    const el = toValue(props.el);

    if (!el || !container.value) {
        return;
    }

    const menuRect = container.value.getBoundingClientRect();
    const elRect = el.getBoundingClientRect();
    const pos = {
        top: "",
        left: "",
        right: "",
        bottom: "",
    };

    const underneath = menuRect.height + elRect.bottom >= viewport.value.height;
    const leftAlign = menuRect.width + elRect.left >= viewport.value.width;

    if (underneath) {
        pos.top = elRect.bottom + "px";
    } else {
        pos.bottom = viewport.value.height - elRect.top + "px";
    }

    if (leftAlign) {
        pos.left = elRect.left + "px";
    } else {
        pos.right = viewport.value.width - elRect.right + "px";
    }

    return pos;
});
</script>

<template>
    <Teleport to="body">
        <div v-if="show" ref="container" class="dropdown" :style="menuPosition">
            <slot></slot>
        </div>
    </Teleport>
</template>

<style lang="scss" scoped>
.dropdown {
    position: absolute;
    display: flex;
    flex-direction: column;
    z-index: 999;
    width: max-content;
    box-shadow: 0 0 0.8rem black;
}
</style>
