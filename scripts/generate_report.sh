#!/bin/bash

# A helper function to run the jq pipeline on a given input file
generate_table() {
  local input="$1"
  jq -r -s '
    ("(?<prefix>(SKIPPING )?TESTCASE): (?<name>[^\n]*)") as $regex |
    # Build a map: test -> { status, custom_name, go_name }
    reduce (.[] | select(.Test != null)) as $event ({};
      ($event.Test) as $key |
      if $event.Action == "pass" or $event.Action == "fail" or $event.Action == "skip" then
        .[$key] += {status: $event.Action}
      elif $event.Action == "output" and ($event.Output | test($regex)) then
        .[$key] += {custom_name: ($event.Output | capture($regex) | .name)}
      else
        .[$key] += {go_name: $key}
      end
    )
    | to_entries | map(.value) | map(select(.status))
    # Extract suite/subgroup/name
    | map(
        . + (
          {
            # Clean the path first: remove leading ../ or ./
            clean_path: ((.custom_name // .go_name) | gsub("^(\\.\\./|\\./)"; "")),
            path_parts: ((.custom_name // .go_name) | gsub("^(\\.\\./|\\./)"; "") | split("/"))
          } |
          if (.path_parts | length) >= 3 then
            {
              suite: .path_parts[0],
              subgroup: .path_parts[1],
              name: .path_parts[-1]
            }
          elif (.path_parts | length) == 2 then
            {
              suite: .path_parts[0],
              subgroup: null,
              name: .path_parts[1]
            }
          else
            {
              suite: "other",
              subgroup: null,
              name: .path_parts[0]
            }
          end
        )
      )
    # Group by suite
    | group_by(.suite)
    # For each suite:
    | .[]
    | (
        # Suite header row
        "| **\(. [0].suite)** |  |  |"
      )
      + "\n"
      +
      (
        group_by(.subgroup)
        | map(
            # Subgroup header (if exists)
            (
              if .[0].subgroup != null then
                "| **\(. [0].suite) / \(. [0].subgroup)** |  |  |\n"
              else
                ""
              end
            )
            +
            (
              map(
                "|  | `\(.name)` | " +
                (if .status == "pass" then "✅"
                 elif .status == "fail" then "❌"
                 elif .status == "skip" then "⏩"
                 else "❓" end) + " |"
              )
              | join("\n")
            )
          )
        | join("\n")
      )
  ' "$input"
}

{
  #
  # Main test results
  #
  echo "# Test Results"
  echo ""
  echo "| Suite | ROM       | Status |"
  echo "|-------|-----------|--------|"
  generate_table "test_results.json"
  echo ""
  #
  # Screenshot test results
  #
  echo "# Screenshot Test Results"
  echo ""
  echo "| Suite | ROM       | Status |"
  echo "|-------|-----------|--------|"
  generate_table "screenshots_test_results.json"
  echo ""
} > TESTS.md

echo "--- Generated TESTS.md ---"
cat TESTS.md
echo "--------------------------"