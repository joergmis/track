Feature: fetch customers

    Scenario: Fetch customers on startup
        Given the customer repository is set up
        When adding a new activity
        Then it should fetch the customers from the repository
