Feature: start a time entry

    Background:
        Given the date is "2024-07-01 08:48:12"
        And no time entries are present
        And the customer "benni" with project "car" exists

    Scenario: Start a new activity
        When starting a new activity for "benni" "car" "service" "change_oil"
        Then a time entry is added for "2024-07-01 08:48:12" "benni" "car" "service" "change_oil"
        And there are 1 time entries in the database

    Scenario: Start a new activity
        Given the date is "2024-07-01 08:48:12"
        When starting a new activity for "benni" "car" "service" "change_oil"
        Then a time entry is added for "2024-07-01 08:48:12" "benni" "car" "service" "change_oil"

        Given the date is "2024-07-01 08:50:12"
        When starting a new activity for "benni" "car" "service" "change_tires"
        Then a time entry is added for "2024-07-01 08:50:12" "benni" "car" "service" "change_tires"

        And there are 2 time entries in the database

