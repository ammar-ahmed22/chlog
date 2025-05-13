import { Badge } from "@/components/ui/badge";

export interface TagFilterProps {
  tags: string[];
  selectedTags: string[];
  onTagSelect: (tag: string) => void;
  onTagDeselect: (tag: string) => void;
}

export default function TagFilter({
  tags,
  selectedTags,
  onTagSelect,
  onTagDeselect,
}: TagFilterProps) {
  if (tags.length === 0) {
    return null;
  }
  return (
    <div className="flex flex-wrap gap-2 items-center">
      {tags.map((tag) => {
        const isSelected = selectedTags.includes(tag);
        return (
          <Badge
            key={`tag-${tag}`}
            variant={isSelected ? "default" : "outline"}
            className="text-xs cursor-pointer"
            onClick={() =>
              isSelected ? onTagDeselect(tag) : onTagSelect(tag)
            }>
            {tag}
          </Badge>
        );
      })}
    </div>
  );
}
