#!/bin/bash

{
  echo "# Test Results"
  echo ""
  echo "| Suite | ROM       | Status |"
  echo "|-------|-----------|--------|"

  jq -r -s '

    # 1. Define a regex to extract the custom test name from test output lines like "TESTCASE: suite/group/testname"
    ("(?<prefix>(SKIPPING )?TESTCASE): (?<name>[^\n]*)") as $regex |

    # 2. Iterate over all events (passed in as a stream of JSON objects), and reduce into a map keyed by test name
    reduce (.[] | select(.Test != null)) as $event ({};
      ($event.Test) as $key |

      # 3. If the event is a pass/fail/skip, store that as the status
      if $event.Action == "pass" or $event.Action == "fail" or $event.Action == "skip" then
        .[$key] += {status: $event.Action}

      # 4. If the event is an output line matching the regex, extract the custom name from the output
      elif $event.Action == "output" and ($event.Output | test($regex)) then
        .[$key] += {custom_name: ($event.Output | capture($regex) | .name)}

      # 5. Otherwise, just remember the original Go test name
      else
        .[$key] += {go_name: $key}
      end
    )

    # 6. Convert map to array of values, keep only tests that have a status
    | to_entries | map(.value) | map(select(.status))

    # 7. For each test, split its name (custom or fallback to Go test name) by `/` and assign suite, subgroup, name
    | map(
        . + (
          {
            path_parts: (.custom_name // .go_name | split("/"))
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
              suite: "Ungrouped",
              subgroup: null,
              name: .path_parts[0]
            }
          end
        )
      )

    # 8. Group tests by suite (top-level group)
    | group_by(.suite)

    | .[]
    | (
        # Emit suite group row
        "| **\(. [0].suite)** |  |  |"
      )
      + "\n"
      +
      (
        group_by(.subgroup)
        | map(
            # Subgroup header row (only if subgroup exists)
            (
              if .[0].subgroup != null then
                "| **\(. [0].suite) / \(. [0].subgroup)** |  |  |\n"
              else
                ""
              end
            )
            +
            # Tests inside subgroup
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
  ' test_results.json

} > TESTS.md

echo "--- Generated TESTS.md ---"
cat TESTS.md
echo "--------------------------"
