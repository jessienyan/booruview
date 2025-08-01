/// <reference types="vite/client" />

type TagType =
    | "tag"
    | "artist"
    | "copyright"
    | "character"
    | "metadata"
    | "deprecated"
    | "unknown";

type Tag = {
    count: number;
    name: string;
    type: TagType;
};

type ChipStyle = "blacklist" | "checkmark" | "strikethrough" | "default";

type TagChip = {
    style: ChipStyle;
    tag: Tag;
};

type Rating = "general" | "sensitive" | "questionable" | "explicit";

type Post = {
    id: number;
    created_at: number;
    score: number;
    rating: Rating;
    source_url: string;
    uploader: string;
    tags: string[];
    thumbnail_url: string;
    thumbnail_width: number;
    thumbnail_height: number;
    lowres_url: string;
    lowres_width: number;
    lowres_height: number;
    image_url: string;
    width: number;
    height: number;
};
