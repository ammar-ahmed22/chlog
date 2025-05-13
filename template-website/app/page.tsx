"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { compareAsc, format, parseISO, compareDesc } from "date-fns";
import {
  ChevronDown,
  ChevronRight,
  Loader,
  Search,
  Tag,
} from "lucide-react";
import TagFilter from "./tag-filter";
import DateRangeFilter from "./date-range-filter";
import { DateRange } from "react-day-picker";

type ChangelogChange = {
  commit_hash: string;
  tags: string[];
  details: string;
};

type ChangelogEntry = {
  version: string;
  date: string;
  from_ref: string;
  to_ref: string;
  changes: ChangelogChange[];
};

export default function Home() {
  const [changelog, setChangelog] = useState<ChangelogEntry[]>([]);
  const [fetchingError, setFetchingError] = useState(false);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    setLoading(true);
    fetch("/changelog.json")
      .then((response) => response.json())
      .then((data) => {
        setChangelog(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Error fetching changelog:", error);
        setFetchingError(true);
        setLoading(false);
      });
  }, []);
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [selectedVersion, setSelectedVersion] = useState<string>("");
  const [dateRange, setDateRange] = useState<DateRange | undefined>(
    undefined,
  );
  const [expandedEntries, setExpandedEntries] = useState<
    Record<string, boolean>
  >({});

  // Get all unique tags from the changelog
  const allTags = useMemo(() => {
    const tags = new Set<string>();
    changelog.forEach((entry) => {
      entry.changes.forEach((change) => {
        change.tags.forEach((tag) => tags.add(tag));
      });
    });
    return Array.from(tags).sort();
  }, [changelog]);

  // Get all unique versions from the changelog
  const allVersions = useMemo(() => {
    return changelog.map((entry) => entry.version);
  }, [changelog]);

  // Toggle expanded state for an entry
  const toggleExpanded = (version: string) => {
    setExpandedEntries((prev) => ({
      ...prev,
      [version]: !prev[version],
    }));
  };

  // Add a tag to the filter
  const addTagFilter = (tag: string) => {
    if (!selectedTags.includes(tag)) {
      setSelectedTags([...selectedTags, tag]);
    }
  };

  // Remove a tag from the filter
  const removeTagFilter = (tag: string) => {
    setSelectedTags(selectedTags.filter((t) => t !== tag));
  };

  // Clear all filters
  const clearFilters = () => {
    setSearchQuery("");
    setSelectedTags([]);
    setSelectedVersion("");
    setDateRange(undefined);
  };

  const hasFilters = useMemo(() => {
    return (
      searchQuery.length > 0 ||
      selectedTags.length > 0 ||
      (selectedVersion !== "" && selectedVersion !== "all") ||
      dateRange !== undefined
    );
  }, [searchQuery, selectedTags, selectedVersion, dateRange]);

  const matchChange = useCallback(
    (change: ChangelogChange) => {
      // Filter by search query
      const matchesSearch =
        !searchQuery ||
        change.details
          .toLowerCase()
          .includes(searchQuery.toLowerCase()) ||
        change.commit_hash
          .toLowerCase()
          .includes(searchQuery.toLowerCase());

      // Filter by tags
      const matchesTags =
        selectedTags.length === 0 ||
        selectedTags.every((tag) => change.tags.includes(tag));

      return matchesSearch && matchesTags;
    },
    [searchQuery, selectedTags],
  );

  // Filter the changelog based on all filters
  const filteredChangelog = useMemo(() => {
    return changelog.filter((entry) => {
      // Filter by version
      if (
        selectedVersion &&
        selectedVersion !== "all" &&
        entry.version !== selectedVersion
      ) {
        return false;
      }

      // Filter by date range
      if (dateRange?.from && parseISO(entry.date) < dateRange?.from) {
        return false;
      }
      if (dateRange?.to && parseISO(entry.date) > dateRange?.to) {
        return false;
      }

      // Filter by search query and tags at the change level
      const matchingChanges = entry.changes.filter(matchChange);

      return matchingChanges.length > 0;
    });
  }, [changelog, selectedVersion, dateRange, matchChange]);

  const dateRangeBounds: DateRange | undefined = useMemo(() => {
    if (changelog.length !== 0) {
      const parsedDates = changelog.map((entry) =>
        parseISO(entry.date),
      );
      const oldest = parsedDates.reduce((min, curr) =>
        compareAsc(min, curr) <= 0 ? min : curr,
      );
      const newest = parsedDates.reduce((max, curr) =>
        compareDesc(max, curr) <= 0 ? max : curr,
      );
      return {
        from: oldest,
        to: newest,
      };
    }

    return undefined;
  }, [changelog]);

  useEffect(() => {
    // If has filters, expand all
    if (hasFilters && filteredChangelog.length > 0) {
      setExpandedEntries((prev) => {
        const copy = { ...prev };
        filteredChangelog.forEach((entry) => {
          if (!(entry.version in prev) || !prev[entry.version]) {
            copy[entry.version] = true;
          }
        });
        return copy;
      });
    } else {
      setExpandedEntries({});
    }
  }, [hasFilters, filteredChangelog, setExpandedEntries]);
  // Format date for display
  const formatDate = (dateString: string) => {
    return format(parseISO(dateString), "MMMM d, yyyy");
  };

  return (
    <div className="container mx-auto py-8 px-4 max-w-5xl">
      <h1 className="text-3xl font-bold mb-8">Changelog</h1>

      {/* Filters */}
      <div className="mb-8 space-y-4">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="relative flex-grow">
            <Search
              className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
              size={18}
            />
            <Input
              placeholder="Search changes..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>

          <Select
            value={selectedVersion}
            onValueChange={setSelectedVersion}>
            <SelectTrigger className="w-full md:w-[180px] cursor-pointer">
              <Tag className="size-4" />
              <SelectValue placeholder="Filter by version" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All versions</SelectItem>
              {allVersions.map((version) => (
                <SelectItem key={version} value={version}>
                  {version}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <DateRangeFilter
            dateRange={dateRange}
            onDateRangeChange={setDateRange}
            disabled={(date) => {
              if (dateRangeBounds) {
                if (
                  dateRangeBounds.from &&
                  date < dateRangeBounds.from
                ) {
                  return true;
                }
                if (dateRangeBounds.to && date > dateRangeBounds.to) {
                  return true;
                }
              }
              return false;
            }}
          />
        </div>

        <TagFilter
          tags={allTags}
          selectedTags={selectedTags}
          onTagSelect={addTagFilter}
          onTagDeselect={removeTagFilter}
        />
      </div>

      {loading && (
        <div className="flex items-center justify-center py-12">
          <Loader className="animate-spin size-8" />
        </div>
      )}

      {/* Fetching error */}
      {!loading && fetchingError && (
        <div className="text-center py-12">
          <p className="text-xl font-bold">Oops!</p>
          <p className="text-gray-500">
            There was an error fetching the changelog.
          </p>
        </div>
      )}
      {/* Changelog entries */}
      {!loading &&
      !fetchingError &&
      filteredChangelog.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500">
            No changelog entries match your filters.
          </p>
          <Button variant="link" onClick={clearFilters}>
            Clear all filters
          </Button>
        </div>
      ) : (
        <div className="space-y-6">
          {filteredChangelog.map((entry) => {
            // Filter changes based on search query and selected tags
            const filteredChanges = entry.changes.filter(matchChange);

            if (filteredChanges.length === 0) return null;

            return (
              <Card
                key={entry.version}
                className="overflow-hidden py-0 gap-0">
                <CardHeader className="bg-gray-50 dark:bg-gray-900 py-4">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-xl flex items-center">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="p-0 mr-2 h-6 w-6 cursor-pointer"
                        onClick={() => toggleExpanded(entry.version)}>
                        {expandedEntries[entry.version] ? (
                          <ChevronDown size={18} />
                        ) : (
                          <ChevronRight size={18} />
                        )}
                        <span className="sr-only">
                          {expandedEntries[entry.version]
                            ? "Collapse"
                            : "Expand"}
                        </span>
                      </Button>
                      Version {entry.version}
                    </CardTitle>
                    <div className="text-sm text-gray-500">
                      {formatDate(entry.date)}
                    </div>
                  </div>
                </CardHeader>

                <CardContent className="p-0">
                  <Collapsible
                    open={expandedEntries[entry.version]}
                    onOpenChange={() =>
                      toggleExpanded(entry.version)
                    }>
                    <CollapsibleTrigger className="w-full text-left cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors p-4 border-t border-gray-100 dark:border-gray-800">
                      <div className="flex justify-between items-center">
                        <div className="text-sm text-gray-500">
                          {filteredChanges.length} change
                          {filteredChanges.length !== 1 ? "s" : ""}
                        </div>
                        <div className="text-sm text-gray-500">
                          {entry.from_ref} → {entry.to_ref}
                        </div>
                      </div>
                    </CollapsibleTrigger>

                    <CollapsibleContent>
                      <ul className="divide-y divide-gray-100 dark:divide-gray-800">
                        {filteredChanges.map((change, index) => (
                          <li
                            key={change.commit_hash + index}
                            className="p-4 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors">
                            <div className="flex flex-col gap-2">
                              <div className="flex items-center gap-2">
                                <code className="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded font-mono">
                                  {change.commit_hash.substring(0, 7)}
                                </code>
                                <div className="flex flex-wrap gap-1">
                                  {change.tags.map((tag) => (
                                    <Badge
                                      key={tag}
                                      variant="outline"
                                      className="text-xs">
                                      {tag}
                                    </Badge>
                                  ))}
                                </div>
                              </div>
                              <p className="text-sm">
                                {change.details}
                              </p>
                            </div>
                          </li>
                        ))}
                      </ul>
                    </CollapsibleContent>
                  </Collapsible>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
