import { useNewFeatureIndicator } from "@/composable";

export default {
    maxPostHeight: useNewFeatureIndicator(
        "max-post-height",
        new Date("2025-10-17"),
    ),
};
