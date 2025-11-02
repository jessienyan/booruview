const tagTypeSortOrder: Record<TagType, number> = {
    artist: 0,
    character: 1,
    copyright: 2,
    tag: 3,
    deprecated: 3,
    metadata: 4,
    unknown: 5,
};

export function sortTags(tags: Tag[]): Tag[] {
    return [...tags].sort((a, b) => {
        const typeCompare = tagTypeSortOrder[a.type] - tagTypeSortOrder[b.type];

        if (typeCompare !== 0) {
            return typeCompare;
        }

        return a.name.localeCompare(b.name);
    });
}
