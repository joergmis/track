Feature: start a time entry

    Scenario: Start a new activity
        Given the date is "2024-07-01 08:48:12"
        And the customer "benni" with project "car" and service "change_oil" exists

        When starting a new activity for "benni" "car" "service" "change_oil"

        Then a time entry is added for "2024-07-01 08:48:12" "benni" "car" "service" "change_oil"

    Scenario: Start a second activity
        Given the date is "2023-03-01 10:12:03"
        And the customer "beat" with project "garden" and service "cutting" exists
        And the customer "ralph" with project "house" and service "renovation" exists
        And there is an activity running for "beat" "garden" "cutting" "cut_roses" started on "2023-03-01 10:12:03"

        When starting a new activity for "ralph" "house" "renovation" "paint_walls"

        Then a time entry is added for "2023-03-01 10:12:03" "ralph" "house" "renovation" "paint_walls"
