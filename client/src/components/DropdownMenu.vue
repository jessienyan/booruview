<script setup lang="ts">
import { computed, type MaybeRefOrGetter, toValue, useTemplateRef } from "vue";
import { useDismiss, useViewportSize } from "@/composable";

const container = useTemplateRef("container");
const props = defineProps<{ el: MaybeRefOrGetter<HTMLElement | null> }>();
const show = defineModel("show", { required: true });

const viewport = useViewportSize();
useDismiss(
	() => [container, props.el],
	() => {
		show.value = false;
	},
);

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

	const underneath = menuRect.height + elRect.bottom < viewport.value.height;
	const leftAlign = menuRect.width + elRect.left < viewport.value.width;

	if (underneath) {
		pos.top = `${elRect.bottom}px`;
	} else {
		pos.bottom = `${viewport.value.height - elRect.top}px`;
	}

	if (leftAlign) {
		pos.left = `${elRect.left}px`;
	} else {
		pos.right = `${viewport.value.width - elRect.right}px`;
	}

	return pos;
});
</script>

<template>
    <Teleport to="body">
        <div v-if="show" ref="container" class="dropdown" :style="menuPosition">
            <slot>
                <!-- The .dropdown-option class is available for slotted content -->
            </slot>
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
    margin: 4px 0;
}

.dropdown :deep(.dropdown-option) {
    padding: 12px;
    text-align: left;

    .bi {
        margin-right: 0.4rem;
    }

    &:first-of-type {
        border-top-left-radius: 4px;
        border-top-right-radius: 4px;
    }

    &:last-of-type {
        border-bottom-left-radius: 4px;
        border-bottom-right-radius: 4px;
    }

    &:not(:last-of-type) {
        border-bottom: 0;
    }
}
</style>
