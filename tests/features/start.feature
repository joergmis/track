Feature: start a time entry

    Scenario: Start a new activity
        Given the date is "2024-07-01 08:48:12"
        And the customer "benni" with project "car" exists

        When starting a new activity for "benni" "car" "change_oil"

        Then a time entry is added for "2024-07-01 08:48:12" "benni" "car" "change_oil"
