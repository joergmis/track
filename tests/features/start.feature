Feature: start a time entry

    Scenario: Start a new activity
        Given the date is <date>
        And the customer <customer> with project <project> and service <service> exists
        When starting a new activity for <customer> <project> <service> <description>
        Then a time entry is added for <date> <customer> <project> <service> <description>

        Examples:
            | date                | customer | project | service | description |
            | 2023-03-01 10:12:03 | a        | b       | c       | d           |
